#!/usr/bin/env bash
# Executor for Go toolchain bump. All policy and filtering live in
# scripts/check-go-bump-policy.sh.
#
# Usage: bump-go.sh <version>
#   <version> is the exact go directive (e.g. 1.26.2), no "go" prefix.

set -euo pipefail

if [[ $# -lt 1 || -z "${1}" ]]; then
  echo "usage: bump-go.sh <version>" >&2
  echo "  example: bump-go.sh 1.26.2" >&2
  exit 1
fi

version="${1#go}"

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
ROOT_DIR="$(cd "${SCRIPT_DIR}/.." && pwd)"

# All go.mod files that should track the same Go version.
GO_MOD_FILES=(
  "go.mod"
  "test/app/go.mod"
  "tools/clean/go.mod"
  "tools/githubjobs/go.mod"
  "tools/makejwt/go.mod"
  "tools/openapi2crd/go.mod"
  "tools/openapi2crd/hack/tools/go.mod"
  "tools/scaffolder/go.mod"
  "tools/scandeprecation/go.mod"
  "tools/toolbox/go.mod"
)

printf 'bump-go: bumping Go version to %s\n' "${version}"

# TEST_BUMP_DRY_RUN=1 lets the test suite confirm the script is reached
# without touching real files (check-go-bump-policy.sh uses an absolute path
# to invoke this script so PATH-based stubbing cannot intercept it).
if [[ "${TEST_BUMP_DRY_RUN:-}" == "1" ]]; then
  printf 'bump-go: dry-run, skipping file updates\n'
  exit 0
fi

# Update the go directive in each go.mod file.
# Use a temp-file swap for cross-platform compatibility (GNU vs BSD sed).
for rel in "${GO_MOD_FILES[@]}"; do
  go_mod="${ROOT_DIR}/${rel}"
  if [[ ! -f "${go_mod}" ]]; then
    printf 'bump-go: warning: %s not found, skipping\n' "${go_mod}" >&2
    continue
  fi
  tmpfile=$(mktemp)
  sed -e "s|^go [0-9][0-9.]*$|go ${version}|" \
      -e "/^toolchain go[0-9]/d" \
      "${go_mod}" > "${tmpfile}"
  mv "${tmpfile}" "${go_mod}"
  printf 'bump-go: updated %s\n' "${go_mod}"
done
