#!/usr/bin/env bash
# Update the coordinated Kubernetes staging module versions for a release.
#
# Usage: K8S_RELEASE=1.37.0 make update-mod

set -euo pipefail

if [[ -z "${K8S_RELEASE:-}" ]]; then
	echo "K8S_RELEASE not set. Example: export K8S_RELEASE=1.37.0" >&2
	exit 1
fi

# Kubernetes staging modules use v0.<minor>.0 tags while the main module uses
# v1.<minor>.0. This target updates stable Kubernetes releases only.
if [[ ! "${K8S_RELEASE}" =~ ^1\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)$ ]]; then
	echo "Invalid K8S_RELEASE '${K8S_RELEASE}'. Expected a stable release such as 1.37.0." >&2
	exit 1
fi

minor="${BASH_REMATCH[1]}"
staging_version="v0.${minor}.0"
kubernetes_version="v1.${minor}.0"

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
TEMP_MOD=""
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
	if [[ -n "${TEMP_MOD}" ]]; then
		if ! rm -f "${TEMP_MOD}"; then
			echo "Failed to remove temporary file ${TEMP_MOD}." >&2
			if (( status == 0 )); then
				status=1
			fi
		fi
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

update_versions() {
	local go_mod="$1"
	local temp_mod
	temp_mod="$(mktemp "$(dirname "${go_mod}")/.go.mod.XXXXXX")"
	TEMP_MOD="${temp_mod}"

	if ! awk -v staging_version="${staging_version}" -v kubernetes_version="${kubernetes_version}" '
		function replace_versions(line, version_pattern, replacement, position, version_length, following, output, comment_position, comment) {
			comment_position = index(line, "//")
			comment = ""
			if (comment_position > 0) {
				comment = substr(line, comment_position)
				line = substr(line, 1, comment_position - 1)
			}
			output = ""
			while ((position = match(line, version_pattern)) != 0) {
				version_length = RLENGTH
				following = substr(line, position + version_length, 1)
				if (following == "" || following == " " || following == "\t") {
					output = output substr(line, 1, position - 1) replacement
					line = substr(line, position + version_length)
				} else {
					output = output substr(line, 1, position + version_length - 1)
					line = substr(line, position + version_length)
				}
			}
			return output line comment
		}
		{
			line = $0
			if (line ~ /^[[:space:]]*(require[[:space:]]+)?k8s\.io\/kubernetes[[:space:]]+v1\.[0-9][0-9]*\.[0-9][0-9]*([[:space:]]|$)/) {
				candidate_count++
				line = replace_versions(line, "v1\\.[0-9][0-9]*\\.[0-9][0-9]*", kubernetes_version)
			} else if (line ~ /^[[:space:]]*(require[[:space:]]+)?k8s\.io\/[^[:space:]]+[[:space:]]+v0\.[0-9][0-9]*\.[0-9][0-9]*([[:space:]]|$)/) {
				candidate_count++
				line = replace_versions(line, "v0\\.[0-9][0-9]*\\.[0-9][0-9]*", staging_version)
			} else {
				arrow = index(line, "=>")
				if (arrow > 0) {
					left = substr(line, 1, arrow - 1)
					right = substr(line, arrow + 2)
					if (left ~ /^[[:space:]]*(replace[[:space:]]+)?k8s\.io\/[^[:space:]]+([[:space:]]+v0\.[0-9][0-9]*\.[0-9][0-9]*)?[[:space:]]*$/ &&
					    right ~ /^[[:space:]]*k8s\.io\/[^[:space:]]+[[:space:]]+v0\.[0-9][0-9]*\.[0-9][0-9]*([[:space:]]|$)/) {
						candidate_count++
						line = replace_versions(line, "v0\\.[0-9][0-9]*\\.[0-9][0-9]*", staging_version)
					}
				}
			}
			print line
		}
		END {
			if (candidate_count == 0) {
				exit 2
			}
		}
	' "${go_mod}" > "${temp_mod}"; then
		echo "Could not update Kubernetes versions in ${go_mod}." >&2
		return 1
	fi

	if ! cmp -s "${go_mod}" "${temp_mod}"; then
		cat "${temp_mod}" > "${go_mod}"
	fi
	rm -f "${temp_mod}"
	TEMP_MOD=""
}

echo "Updating Kubernetes module pins for Kubernetes ${K8S_RELEASE} (${staging_version}, ${kubernetes_version})."
for module in "${MODULES[@]}"; do
	go_mod="${REPO_ROOT}/${module}/go.mod"
	update_versions "${go_mod}"

	echo "Running go mod tidy in ${module}."
	(
		cd "${REPO_ROOT}/${module}"
		go mod tidy
	)
done
