#!/bin/bash
# Verify that a swagger.json was generated with OpenAPIEnums=true by counting
# how many non-empty "enum" arrays it contains. Fails if the count is below a
# conservative minimum, so a contributor cannot proceed with enum-free swagger.
#
# Usage: verify-enum-swagger.sh <path-to-swagger.json> [min-enum-arrays]
#
# Reusable across the temporary-checkout, maintainer-checkout, and future
# artifact-fetch workflows.

set -euo pipefail

SWAGGER="${1:-}"
MIN_ENUMS="${2:-50}"

if [ -z "${SWAGGER}" ]; then
	echo "Usage: $0 <path-to-swagger.json> [min-enum-arrays]" >&2
	exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
	echo "jq is required but was not found in PATH." >&2
	exit 1
fi

if [ ! -f "${SWAGGER}" ]; then
	echo "swagger file not found: ${SWAGGER}" >&2
	exit 1
fi

# Recursively count every non-empty "enum" array anywhere in the document.
ENUM_COUNT="$(jq '[.. | objects | select(has("enum")) | .enum | select(length > 0)] | length' "${SWAGGER}")"

if [ "${ENUM_COUNT}" -lt "${MIN_ENUMS}" ]; then
	echo "Enum verification FAILED: ${SWAGGER} has ${ENUM_COUNT} non-empty enum arrays (expected at least ${MIN_ENUMS})." >&2
	echo "The swagger was likely generated without OpenAPIEnums=true." >&2
	exit 1
fi

echo "Enum verification passed: ${ENUM_COUNT} non-empty enum arrays in ${SWAGGER}."
