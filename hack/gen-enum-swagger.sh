#!/bin/bash
# Generate enum-enabled OpenAPI swagger.json from a temporary Kubernetes source
# checkout, so release contributors do not need a maintainer-managed, manually
# patched k/k clone.
#
# Steps: shallow-clone the release tag, patch only that temporary checkout to
# enable OpenAPIEnums=true, run k/k's existing hack/update-openapi-spec.sh, copy
# only api/openapi-spec/swagger.json into gen-apidocs, verify enum metadata, and
# delete the temporary checkout (KEEP_TMP=1 preserves it for debugging).
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

# Patch only this temporary checkout. k/k hardcodes OpenAPIEnums=false on the
# kube-apiserver --feature-gates line; flip it to true for enum-enabled output.
echo "Enabling OpenAPIEnums=true in the temporary checkout"
sed -i.bak 's/OpenAPIEnums=false/OpenAPIEnums=true/' "${KK}/hack/update-openapi-spec.sh"
rm -f "${KK}/hack/update-openapi-spec.sh.bak"
if ! grep -q 'OpenAPIEnums=true' "${KK}/hack/update-openapi-spec.sh"; then
	echo "Failed to enable OpenAPIEnums in ${KK}/hack/update-openapi-spec.sh." >&2
	echo "The k/k script format may have changed for ${TAG}; patch it manually." >&2
	exit 1
fi

echo "Running k/k hack/update-openapi-spec.sh (logging to ${GEN_LOG})"
( cd "${KK}" && hack/update-openapi-spec.sh ) 2>&1 | tee "${GEN_LOG}"

mkdir -p "${OUT_DIR}"
echo "Copying swagger.json into ${OUT_SWAGGER}"
cp "${KK}/api/openapi-spec/swagger.json" "${OUT_SWAGGER}"

echo "Verifying enum metadata"
"${SCRIPT_DIR}/verify-enum-swagger.sh" "${OUT_SWAGGER}"

echo "Enum-enabled swagger.json ready at ${OUT_SWAGGER}"
