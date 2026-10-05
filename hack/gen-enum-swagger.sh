#!/bin/bash
# Generate enum-enabled OpenAPI swagger.json from a temporary Kubernetes source
# checkout, so release contributors do not need a maintainer-managed, manually
# patched k/k clone.
#
# Steps: shallow-clone the release tag, run k/k's hack/update-openapi-spec.sh
# with enums kept, copy only api/openapi-spec/swagger.json into gen-apidocs,
# verify enum metadata, and delete the temporary checkout (KEEP_TMP=1 preserves
# it for debugging).
#
# Required env: K8S_RELEASE (e.g. 1.36.0)
# Pass-through env (read directly by k/k): TMP_DIR, ETCD_PORT, API_PORT, API_LOGFILE
# Debug: KEEP_TMP=1 keeps the temporary checkout and generation log.

set -euo pipefail

if [ -z "${K8S_RELEASE:-}" ]; then
	echo "K8S_RELEASE not set. Example: export K8S_RELEASE=1.36.0" >&2
	exit 1
fi

# Preflight: obvious local tools only. k/k's build reports deeper problems.
for tool in git go jq curl openssl; do
	if ! command -v "${tool}" >/dev/null 2>&1; then
		echo "${tool} is required but was not found in PATH." >&2
		exit 1
	fi
done

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(dirname "${SCRIPT_DIR}")"

TAG="v${K8S_RELEASE}"
# 1.36.0 -> v1_36, matching set_version_dirs.sh and the Makefile.
VERSION_DIR="v$(echo "${K8S_RELEASE}" | cut -c 1-4 | sed "s/\./_/g")"
OUT_DIR="${REPO_ROOT}/gen-apidocs/config/${VERSION_DIR}"
OUT_SWAGGER="${OUT_DIR}/swagger.json"

TMPROOT="$(mktemp -d)"
KK="${TMPROOT}/kubernetes"
GEN_LOG="${TMPROOT}/gen-openapi.log"

cleanup() {
	if [ "${KEEP_TMP:-}" = "1" ]; then
		echo "KEEP_TMP=1 set; preserving temporary checkout:"
		echo "  checkout: ${KK}"
		echo "  log:      ${GEN_LOG}"
	else
		chmod -R u+w "${TMPROOT}" 2>/dev/null || true
		rm -rf "${TMPROOT}"
	fi
}
trap cleanup EXIT

echo "Cloning kubernetes/kubernetes at ${TAG} (shallow) into ${KK}"
git clone --depth 1 --branch "${TAG}" \
	https://github.com/kubernetes/kubernetes.git "${KK}"

# k/k omits enums from its checked-in spec. Since v1.38 it keeps them when
# KUBE_OPENAPI_SPEC_KEEP_ENUMS=true. Older releases ignore that variable and
# hardcode OpenAPIEnums=false on the kube-apiserver --feature-gates line, so
# flip it to true in this temporary checkout only.
UPDATE_SCRIPT="${KK}/hack/update-openapi-spec.sh"
if grep -q 'KUBE_OPENAPI_SPEC_KEEP_ENUMS' "${UPDATE_SCRIPT}"; then
	echo "Keeping enums with KUBE_OPENAPI_SPEC_KEEP_ENUMS=true"
elif grep -q 'OpenAPIEnums=false' "${UPDATE_SCRIPT}"; then
	echo "Enabling OpenAPIEnums=true in the temporary checkout"
	sed -i.bak 's/OpenAPIEnums=false/OpenAPIEnums=true/' "${UPDATE_SCRIPT}"
	rm -f "${UPDATE_SCRIPT}.bak"
else
	echo "Cannot keep enums: ${UPDATE_SCRIPT} has neither KUBE_OPENAPI_SPEC_KEEP_ENUMS nor OpenAPIEnums=false." >&2
	echo "The k/k script format may have changed for ${TAG}; patch it manually." >&2
	exit 1
fi

echo "Running k/k hack/update-openapi-spec.sh (logging to ${GEN_LOG})"
( cd "${KK}" && KUBE_OPENAPI_SPEC_KEEP_ENUMS=true hack/update-openapi-spec.sh ) 2>&1 | tee "${GEN_LOG}"

mkdir -p "${OUT_DIR}"
echo "Copying swagger.json into ${OUT_SWAGGER}"
cp "${KK}/api/openapi-spec/swagger.json" "${OUT_SWAGGER}"

echo "Verifying enum metadata"
"${SCRIPT_DIR}/verify-enum-swagger.sh" "${OUT_SWAGGER}"

echo "Enum-enabled swagger.json ready at ${OUT_SWAGGER}"
