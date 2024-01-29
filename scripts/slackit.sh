#!/bin/bash
set -euo pipefail
MESSAGE=$(cat)
WEBHOOK=${1}
curl -X POST -d "{\"text\":\"${MESSAGE}\"}" "${WEBHOOK}"
