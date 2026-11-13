#!/usr/bin/env bash
# Update the coordinated Kubernetes staging module versions for a release.
#
# Usage: K8S_RELEASE=1.37.0 make update-mod

set -euo pipefail

if [[ -z "${K8S_RELEASE:-}" ]]; then
	echo "K8S_RELEASE not set. Example: export K8S_RELEASE=1.37.0" >&2
	exit 1
fi

# Accept stable releases and release candidates. Keep the full patch and RC
# suffix when translating the version to the staging module's v0.x.y tags.
if [[ ! "${K8S_RELEASE}" =~ ^1\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)(-rc\.(0|[1-9][0-9]*))?$ ]]; then
	echo "Invalid K8S_RELEASE '${K8S_RELEASE}'. Expected 1.MINOR.PATCH or 1.MINOR.PATCH-rc.N (for example, 1.37.0 or 1.37.0-rc.1); patch versions are supported." >&2
	exit 1
fi

if ! command -v go >/dev/null 2>&1; then
	echo "go is required to update module files." >&2
	exit 1
fi
if ! command -v jq >/dev/null 2>&1; then
	echo "jq is required to discover Kubernetes staging modules from go.mod." >&2
	exit 1
fi

staging_version="v0.${K8S_RELEASE#1.}"
kubernetes_version="v${K8S_RELEASE}"

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
MODULES=(gen-compdocs genref)

for module in "${MODULES[@]}"; do
	if [[ ! -f "${REPO_ROOT}/${module}/go.mod" || ! -f "${REPO_ROOT}/${module}/go.sum" ]]; then
		echo "Expected ${module}/go.mod and ${module}/go.sum under ${REPO_ROOT}." >&2
		exit 1
	fi
done

BACKUP_DIR="$(mktemp -d "${TMPDIR:-/tmp}/update-go-mod.XXXXXX")"
restore_on_failure=0

cleanup() {
	status=$?
	restore_failed=0
	trap - EXIT
	if (( status != 0 && restore_on_failure )); then
		echo "Update failed; restoring the original go.mod and go.sum files." >&2
		for module in "${MODULES[@]}"; do
			for file in go.mod go.sum; do
				backup="${BACKUP_DIR}/${module}.${file}"
				target="${REPO_ROOT}/${module}/${file}"
				if ! cp "${backup}" "${target}"; then
					echo "Failed to restore ${target} from ${backup}." >&2
					restore_failed=1
				fi
			done
		done
	fi
	if (( restore_failed )); then
		echo "Recovery backups retained at ${BACKUP_DIR}." >&2
	else
		if ! rm -rf "${BACKUP_DIR}"; then
			echo "Failed to remove backup directory ${BACKUP_DIR}." >&2
			if (( status == 0 )); then
				status=1
			fi
		fi
	fi
	exit "${status}"
}
trap cleanup EXIT

for module in "${MODULES[@]}"; do
	cp "${REPO_ROOT}/${module}/go.mod" "${BACKUP_DIR}/${module}.go.mod"
	cp "${REPO_ROOT}/${module}/go.sum" "${BACKUP_DIR}/${module}.go.sum"
done
restore_on_failure=1

module_dir() {
	local module="$1"
	(cd "${REPO_ROOT}/${module}" && pwd)
}

module_json() {
	local directory="$1"
	(cd "${directory}" && go mod edit -json)
}

module_edit() {
	local directory="$1"
	shift
	(cd "${directory}" && go mod edit "$@")
}

module_version() {
	local directory="$1"
	local path="$2"
	module_json "${directory}" | jq -er --arg path "${path}" '[.Require[]? | select(.Path == $path) | .Version] | .[0]'
}

contains_staging_module() {
	local candidate="$1"
	local staging_module
	for staging_module in "${STAGING_MODULES[@]}"; do
		if [[ "${staging_module}" == "${candidate}" ]]; then
			return 0
		fi
	done
	return 1
}

is_independently_versioned_module() {
	case "$1" in
		k8s.io/gengo|k8s.io/gengo/v2|k8s.io/kube-openapi|k8s.io/utils|k8s.io/klog/v2|k8s.io/system-validators)
			return 0
			;;
	esac
	return 1
}

echo "Updating k8s.io/kubernetes to ${kubernetes_version} to discover the release staging set."
COMP_DOCS_DIR="$(module_dir gen-compdocs)"
module_edit "${COMP_DOCS_DIR}" "-require=k8s.io/kubernetes@${kubernetes_version}"
KUBERNETES_GOMOD="$(go mod download -json "k8s.io/kubernetes@${kubernetes_version}" | jq -r '.GoMod // empty')" || true
if [[ -z "${KUBERNETES_GOMOD}" || ! -f "${KUBERNETES_GOMOD}" ]]; then
	echo "Could not download go.mod for k8s.io/kubernetes ${kubernetes_version}. Does the tag exist?" >&2
	exit 1
fi

STAGING_MODULES_FILE="${BACKUP_DIR}/staging-modules.txt"
go mod edit -json "${KUBERNETES_DIR}/go.mod" |
	jq -er '.Replace[]? | select(.New.Path | startswith("./staging/src/k8s.io/")) | select(.Old.Path | startswith("k8s.io/")) | .Old.Path' \
	> "${STAGING_MODULES_FILE}"
STAGING_MODULES=()
while IFS= read -r staging_module; do
	[[ -z "${staging_module}" ]] || STAGING_MODULES+=("${staging_module}")
done < "${STAGING_MODULES_FILE}"
if (( ${#STAGING_MODULES[@]} == 0 )); then
	echo "No Kubernetes staging modules found in ${KUBERNETES_DIR}/go.mod." >&2
	exit 1
fi

echo "Found ${#STAGING_MODULES[@]} staging modules in ${KUBERNETES_DIR}/go.mod."

update_module() {
	local module="$1"
	local directory
	local previous_api_version
	local stale_modules_file
	local stale_module
	local old_replacements_file
	local old_replacement

	directory="$(module_dir "${module}")"
	previous_api_version="$(module_version "${directory}" k8s.io/api)"

	# Drop release-pinned modules that were in the previous staging set but no
	# longer appear in this Kubernetes release. The main module's replace set is
	# authoritative for gen-compdocs; genref follows the documented OLD-version
	# selection rule. Explicitly preserve known independent repositories even if
	# one happens to use the same version string as the staging release.
	stale_modules_file="${BACKUP_DIR}/${module}.stale-requires.txt"
	module_json "${directory}" |
		jq -r --arg version "${previous_api_version}" \
			'.Require[]? | select(.Path | startswith("k8s.io/")) | select(.Version == $version) | .Path' \
		> "${stale_modules_file}"
	while IFS= read -r stale_module; do
		[[ -z "${stale_module}" ]] && continue
		if ! contains_staging_module "${stale_module}" && ! is_independently_versioned_module "${stale_module}"; then
			module_edit "${directory}" "-droprequire=${stale_module}"
		fi
	done < "${stale_modules_file}"

	if [[ "${module}" == "gen-compdocs" ]]; then
		old_replacements_file="${BACKUP_DIR}/${module}.stale-replaces.txt"
		module_json "${directory}" |
			jq -r '.Replace[]? | select(.Old.Path | startswith("k8s.io/")) | select(.Old.Path == .New.Path) | select(.New.Version | test("^v0\\.[0-9]+\\.[0-9]+(-rc\\.[0-9]+)?$")) | .Old.Path' \
			> "${old_replacements_file}"
		while IFS= read -r old_replacement; do
			[[ -z "${old_replacement}" ]] && continue
			if ! contains_staging_module "${old_replacement}" && ! is_independently_versioned_module "${old_replacement}"; then
				module_edit "${directory}" "-dropreplace=${old_replacement}"
				module_edit "${directory}" "-droprequire=${old_replacement}"
			fi
		done < "${old_replacements_file}"
	fi

	local staging_module
	for staging_module in "${STAGING_MODULES[@]}"; do
		module_edit "${directory}" "-require=${staging_module}@${staging_version}"
		if [[ "${module}" == "gen-compdocs" ]]; then
			module_edit "${directory}" "-replace=${staging_module}=${staging_module}@${staging_version}"
		fi
	done

	if [[ "${module}" == "gen-compdocs" ]]; then
		module_edit "${directory}" "-require=k8s.io/kubernetes@${kubernetes_version}"
	else
		if module_json "${directory}" | jq -e '.Require[]? | select(.Path == "k8s.io/kubernetes")' >/dev/null; then
			module_edit "${directory}" "-require=k8s.io/kubernetes@${kubernetes_version}"
		fi
	fi

	echo "Running go mod tidy in ${module}."
	(cd "${directory}" && go mod tidy)
}

echo "Updating Kubernetes module pins for Kubernetes ${K8S_RELEASE} (${staging_version}, ${kubernetes_version})."
for module in "${MODULES[@]}"; do
	update_module "${module}"
done
