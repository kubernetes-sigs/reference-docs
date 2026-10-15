#!/usr/bin/env bash
# Fixture tests for update-go-mod.sh. No network access or repo module changes.

set -euo pipefail

SCRIPT_DIR="$(cd -- "$(dirname -- "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd -- "${SCRIPT_DIR}/.." && pwd)"
TEST_ROOT="$(mktemp -d "${TMPDIR:-/tmp}/update-go-mod-test.XXXXXX")"
trap 'rm -rf "${TEST_ROOT}"' EXIT

FIXTURE="${TEST_ROOT}/fixture"
mkdir -p "${FIXTURE}/hack" "${FIXTURE}/gen-compdocs" "${FIXTURE}/genref" "${TEST_ROOT}/bin"
cp "${SCRIPT_DIR}/update-go-mod.sh" "${FIXTURE}/hack/update-go-mod.sh"
cp "${SCRIPT_DIR}/set_version_dirs.sh" "${FIXTURE}/hack/set_version_dirs.sh"

cat > "${TEST_ROOT}/bin/go" <<'EOF'
#!/usr/bin/env bash
set -euo pipefail
if [[ "${1:-}" != "mod" || "${2:-}" != "tidy" ]]; then
	echo "unexpected go invocation: $*" >&2
	exit 97
fi
printf '%s\n' "$PWD" >> "${GO_TIDY_LOG}"
if [[ "${GO_TIDY_FAIL:-}" == "1" && "$PWD" == */genref ]]; then
	printf 'tidy side effect that must roll back\n' >> go.sum
	exit 1
fi
EOF
chmod +x "${TEST_ROOT}/bin/go"

cat > "${TEST_ROOT}/bin/cp" <<'EOF'
#!/usr/bin/env bash
if [[ "${CP_FAIL_RESTORE:-}" == "1" && "$#" == "2" &&
	"$1" == */genref.go.sum && "$2" == */genref/go.sum ]]; then
	echo "simulated restore copy failure" >&2
	exit 73
fi
exec /bin/cp "$@"
EOF
chmod +x "${TEST_ROOT}/bin/cp"

cat > "${FIXTURE}/gen-compdocs/go.mod" <<'EOF'
module example.test/gen-compdocs

go 1.26.0

require (
	k8s.io/component-base v0.34.5
	k8s.io/apimachinery v0.34.2 // pin note v0.34.2 remains documentation
	k8s.io/klog/v2 v2.140.0
	k8s.io/gengo v0.0.0-20230829151522-9cce18d56c01
	k8s.io/utils v0.0.0-20260626114624-be93311217bd // indirect
	k8s.io/system-validators v1.12.1 // indirect
)

require k8s.io/kubernetes v1.34.0

replace (
	k8s.io/api => k8s.io/api v0.34.2
	k8s.io/apimachinery => k8s.io/apimachinery v0.34.1
)

replace k8s.io/component-base => k8s.io/component-base v0.34.0
replace k8s.io/local-fork => ../kubernetes
EOF

cat > "${FIXTURE}/genref/go.mod" <<'EOF'
module example.test/genref

go 1.26.0

require k8s.io/api v0.34.5

require (
	k8s.io/apimachinery v0.34.2
	k8s.io/gengo v0.0.0-20230829151522-9cce18d56c01
	k8s.io/klog/v2 v2.140.0
	k8s.io/kube-openapi v0.0.0-20260721132016-d427ff9ee9ad // indirect
)
EOF

for module in gen-compdocs genref; do
	printf 'sum fixture for %s\n' "${module}" > "${FIXTURE}/${module}/go.sum"
done

count_staging_pins() {
	local file="$1"
	local minor="$2"
	local version_pattern="v0\\.${minor}\\.0"
	local require_count
	local replace_count
	require_count="$(grep -E -c "^[[:space:]]*k8s\\.io/[^[:space:]]+[[:space:]]+${version_pattern}([[:space:]]|$)" "${file}" || true)"
	replace_count="$(grep -E -c "^[[:space:]]*(replace[[:space:]]+)?k8s\\.io/[^[:space:]]+([[:space:]]+${version_pattern})?[[:space:]]*=>[[:space:]]*k8s\\.io/[^[:space:]]+[[:space:]]+${version_pattern}([[:space:]]|$)" "${file}" || true)"
	printf '%s\n' "$((require_count + replace_count))"
}

assert_contains() {
	local file="$1"
	local text="$2"
	if ! grep -Fq -- "${text}" "${file}"; then
		echo "expected '${text}' in ${file}" >&2
		exit 1
	fi
}

export PATH="${TEST_ROOT}/bin:${PATH}"
export GO_TIDY_LOG="${TEST_ROOT}/go-tidy.log"
: > "${GO_TIDY_LOG}"

(cd "${TEST_ROOT}" && K8S_RELEASE=1.37.0 "${FIXTURE}/hack/update-go-mod.sh")

compdocs_mod="${FIXTURE}/gen-compdocs/go.mod"
genref_mod="${FIXTURE}/genref/go.mod"
assert_contains "${compdocs_mod}" 'k8s.io/component-base v0.37.0'
assert_contains "${compdocs_mod}" 'k8s.io/apimachinery v0.37.0 // pin note v0.34.2 remains documentation'
assert_contains "${compdocs_mod}" 'require k8s.io/kubernetes v1.37.0'
assert_contains "${compdocs_mod}" 'k8s.io/api => k8s.io/api v0.37.0'
assert_contains "${compdocs_mod}" 'k8s.io/apimachinery => k8s.io/apimachinery v0.37.0'
assert_contains "${compdocs_mod}" 'replace k8s.io/component-base => k8s.io/component-base v0.37.0'
assert_contains "${compdocs_mod}" 'k8s.io/klog/v2 v2.140.0'
assert_contains "${compdocs_mod}" 'k8s.io/gengo v0.0.0-20230829151522-9cce18d56c01'
assert_contains "${compdocs_mod}" 'k8s.io/utils v0.0.0-20260626114624-be93311217bd'
assert_contains "${compdocs_mod}" 'k8s.io/system-validators v1.12.1'
assert_contains "${compdocs_mod}" 'replace k8s.io/local-fork => ../kubernetes'
assert_contains "${genref_mod}" 'require k8s.io/api v0.37.0'
assert_contains "${genref_mod}" 'k8s.io/api v0.37.0'
assert_contains "${genref_mod}" 'k8s.io/apimachinery v0.37.0'
assert_contains "${genref_mod}" 'k8s.io/klog/v2 v2.140.0'
assert_contains "${genref_mod}" 'k8s.io/gengo v0.0.0-20230829151522-9cce18d56c01'
assert_contains "${genref_mod}" 'k8s.io/kube-openapi v0.0.0-20260721132016-d427ff9ee9ad'

if [[ "$(wc -l < "${GO_TIDY_LOG}" | tr -d ' ')" != "2" ]]; then
	echo "expected go mod tidy to run once per module" >&2
	exit 1
fi

cp "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.after-first.go.mod"
cp "${genref_mod}" "${TEST_ROOT}/genref.after-first.go.mod"
(cd "${TEST_ROOT}" && K8S_RELEASE=1.37.0 "${FIXTURE}/hack/update-go-mod.sh")
cmp -s "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.after-first.go.mod"
cmp -s "${genref_mod}" "${TEST_ROOT}/genref.after-first.go.mod"

# Exercise selection across every concrete staging pin in the checked-in
# modules, while faking tidy to keep this test offline and read-only.
FULL_FIXTURE="${TEST_ROOT}/full-fixture"
mkdir -p "${FULL_FIXTURE}/hack" "${FULL_FIXTURE}/gen-compdocs" "${FULL_FIXTURE}/genref"
cp "${SCRIPT_DIR}/update-go-mod.sh" "${FULL_FIXTURE}/hack/update-go-mod.sh"
current_minor="$(awk '$1 == "k8s.io/kubernetes" { split($2, v, /[.]/); if (v[1] == "v1") { print v[2]; exit } }' "${REPO_ROOT}/gen-compdocs/go.mod")"
if [[ -z "${current_minor}" ]]; then
	echo "could not find the versioned k8s.io/kubernetes requirement" >&2
	exit 1
fi
target_minor="$((current_minor + 1))"
for module in gen-compdocs genref; do
	cp "${REPO_ROOT}/${module}/go.mod" "${FULL_FIXTURE}/${module}/go.mod"
	cp "${REPO_ROOT}/${module}/go.sum" "${FULL_FIXTURE}/${module}/go.sum"
	cp "${REPO_ROOT}/${module}/go.mod" "${TEST_ROOT}/${module}.full.before.go.mod"
done
for module in gen-compdocs genref; do
	before_count="$(count_staging_pins "${REPO_ROOT}/${module}/go.mod" "${current_minor}")"
	if [[ "${before_count}" == "0" ]]; then
		echo "expected concrete v0.${current_minor}.0 pins in ${module}/go.mod" >&2
		exit 1
	fi
done
(cd "${TEST_ROOT}" && K8S_RELEASE="1.${target_minor}.0" "${FULL_FIXTURE}/hack/update-go-mod.sh")
for module in gen-compdocs genref; do
	new_count="$(count_staging_pins "${FULL_FIXTURE}/${module}/go.mod" "${target_minor}")"
	old_count="$(count_staging_pins "${FULL_FIXTURE}/${module}/go.mod" "${current_minor}")"
	before_count="$(count_staging_pins "${REPO_ROOT}/${module}/go.mod" "${current_minor}")"
	if [[ "${new_count}" != "${before_count}" || "${old_count}" != "0" ]]; then
		echo "not all concrete staging pins were updated in ${module}/go.mod" >&2
		exit 1
	fi
done
assert_contains "${FULL_FIXTURE}/gen-compdocs/go.mod" "k8s.io/kubernetes v1.${target_minor}.0"
for module in gen-compdocs genref; do
	for path in k8s.io/klog/v2 k8s.io/system-validators k8s.io/gengo k8s.io/kube-openapi k8s.io/utils; do
		if ! diff -u \
			<(awk -v path="${path}" '$1 == path { print }' "${TEST_ROOT}/${module}.full.before.go.mod") \
			<(awk -v path="${path}" '$1 == path { print }' "${FULL_FIXTURE}/${module}/go.mod"); then
			echo "${path} changed in ${module}/go.mod" >&2
			exit 1
		fi
	done
done

cp "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-invalid.go.mod"
cp "${genref_mod}" "${TEST_ROOT}/genref.before-invalid.go.mod"
invalid_releases=(
	""
	"2.35.0"
	"1.37"
	"1..0"
	"1.37."
	"1.037.0"
	"1.37.00"
	"1.37.0-rc.1"
	"1.37.0+build.1"
	"v1.37.0"
	" 1.37.0"
)
for release in "${invalid_releases[@]}"; do
	if (cd "${TEST_ROOT}" && K8S_RELEASE="${release}" "${FIXTURE}/hack/update-go-mod.sh") > "${TEST_ROOT}/invalid-release.log" 2>&1; then
		echo "expected invalid K8S_RELEASE '${release}' to fail" >&2
		exit 1
	fi
	cmp -s "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-invalid.go.mod"
	cmp -s "${genref_mod}" "${TEST_ROOT}/genref.before-invalid.go.mod"
done

cp "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-tidy-failure.go.mod"
cp "${compdocs_mod%/go.mod}/go.sum" "${TEST_ROOT}/gen-compdocs.before-tidy-failure.go.sum"
cp "${genref_mod}" "${TEST_ROOT}/genref.before-tidy-failure.go.mod"
cp "${genref_mod%/go.mod}/go.sum" "${TEST_ROOT}/genref.before-tidy-failure.go.sum"
if (cd "${TEST_ROOT}" && GO_TIDY_FAIL=1 K8S_RELEASE=1.36.0 "${FIXTURE}/hack/update-go-mod.sh"); then
	echo "expected go mod tidy failure to fail the update" >&2
	exit 1
fi
cmp -s "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-tidy-failure.go.mod"
cmp -s "${compdocs_mod%/go.mod}/go.sum" "${TEST_ROOT}/gen-compdocs.before-tidy-failure.go.sum"
cmp -s "${genref_mod}" "${TEST_ROOT}/genref.before-tidy-failure.go.mod"
cmp -s "${genref_mod%/go.mod}/go.sum" "${TEST_ROOT}/genref.before-tidy-failure.go.sum"

# If restoring a file fails, report it and retain the recovery backups.
cp "${TEST_ROOT}/genref.after-first.go.mod" "${genref_mod}"
cp "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-restore-failure.go.mod"
cp "${compdocs_mod%/go.mod}/go.sum" "${TEST_ROOT}/gen-compdocs.before-restore-failure.go.sum"
cp "${genref_mod}" "${TEST_ROOT}/genref.before-restore-failure.go.mod"
cp "${genref_mod%/go.mod}/go.sum" "${TEST_ROOT}/genref.before-restore-failure.go.sum"
if (cd "${TEST_ROOT}" && TMPDIR="${TEST_ROOT}" CP_FAIL_RESTORE=1 GO_TIDY_FAIL=1 K8S_RELEASE=1.36.0 "${FIXTURE}/hack/update-go-mod.sh") > "${TEST_ROOT}/restore-failure.log" 2>&1; then
	echo "expected tidy and restore failure to fail the update" >&2
	exit 1
fi
if ! grep -Fq 'Failed to restore ' "${TEST_ROOT}/restore-failure.log" ||
	! grep -Fq '/genref/go.sum from ' "${TEST_ROOT}/restore-failure.log"; then
	echo "expected rollback failure to identify genref/go.sum" >&2
	cat "${TEST_ROOT}/restore-failure.log" >&2
	exit 1
fi
backup_dir="$(sed -n 's/^Recovery backups retained at \(.*\)\.$/\1/p' "${TEST_ROOT}/restore-failure.log")"
if [[ -z "${backup_dir}" || ! -d "${backup_dir}" ]]; then
	echo "expected recovery backups to be retained after a restore failure" >&2
	exit 1
fi
cmp -s "${backup_dir}/gen-compdocs.go.mod" "${TEST_ROOT}/gen-compdocs.before-restore-failure.go.mod"
cmp -s "${backup_dir}/gen-compdocs.go.sum" "${TEST_ROOT}/gen-compdocs.before-restore-failure.go.sum"
cmp -s "${backup_dir}/genref.go.mod" "${TEST_ROOT}/genref.before-restore-failure.go.mod"
cmp -s "${backup_dir}/genref.go.sum" "${TEST_ROOT}/genref.before-restore-failure.go.sum"
if cmp -s "${genref_mod%/go.mod}/go.sum" "${TEST_ROOT}/genref.before-restore-failure.go.sum"; then
	echo "expected the injected restore failure to leave genref/go.sum needing recovery" >&2
	exit 1
fi

# A module with no concrete staging pins indicates an unexpected go.mod layout;
# fail and roll back the preceding module's update.
cp "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-no-pins.go.mod"
cp "${genref_mod}" "${TEST_ROOT}/genref.before-no-pins.go.mod"
cat > "${genref_mod}" <<'EOF'
module example.test/genref

go 1.26.0

require k8s.io/gengo v0.0.0-20230829151522-9cce18d56c01
EOF
cp "${genref_mod}" "${TEST_ROOT}/genref.no-pins.go.mod"
if (cd "${TEST_ROOT}" && K8S_RELEASE=1.36.0 "${FIXTURE}/hack/update-go-mod.sh"); then
	echo "expected a module without concrete staging pins to fail" >&2
	exit 1
fi
cmp -s "${compdocs_mod}" "${TEST_ROOT}/gen-compdocs.before-no-pins.go.mod"
cmp -s "${genref_mod}" "${TEST_ROOT}/genref.no-pins.go.mod"

# The moved version-directory script should resolve the repository root from
# its own path, even when invoked from another working directory.
mkdir -p "${FIXTURE}/gen-kubectldocs/generators/v1_36" "${FIXTURE}/gen-apidocs/config/v1_36"
printf 'toc fixture\n' > "${FIXTURE}/gen-kubectldocs/generators/v1_36/toc.yaml"
printf 'config fixture\n' > "${FIXTURE}/gen-apidocs/config/v1_36/config.yaml"
(cd "${TEST_ROOT}" && K8S_RELEASE=1.37.0 "${FIXTURE}/hack/set_version_dirs.sh")
assert_contains "${FIXTURE}/gen-kubectldocs/generators/v1_37/toc.yaml" 'toc fixture'
assert_contains "${FIXTURE}/gen-apidocs/config/v1_37/config.yaml" 'config fixture'
assert_contains "${FIXTURE}/gen-apidocs/config/config.yaml" 'config fixture'

echo "release script tests passed"
