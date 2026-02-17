#!/bin/bash
# Copyright 2025 MongoDB Inc
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# Split generated CRDs into production (individual files for kustomize) and experimental (aggregated file)
#
# This script reads config/generated/crd/bases/crds.yaml and:
# - Production CRDs: extracts to individual files in config/crd/bases/{name}.yaml and adds to kustomization
# - Experimental CRDs: aggregates into config/generated/crd/bases/crds.experimental.yaml
#
# Production CRDs are determined by the kinds listed in config/generated/crd/production-kinds.txt

set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
REPO_ROOT="$(cd "$SCRIPT_DIR/.." && pwd)"
CRD_BASES_DIR="$REPO_ROOT/config/generated/crd/bases"
PRODUCTION_KINDS_FILE="$REPO_ROOT/config/generated/crd/production-kinds.json"
ALL_CRDS_FILE="$CRD_BASES_DIR/crds.yaml"
EXP_CRDS_FILE="$CRD_BASES_DIR/crds.experimental.yaml"
LEGACY_CRD_BASES_DIR="$REPO_ROOT/config/crd/bases"

cd "$REPO_ROOT"

if [ ! -f "$ALL_CRDS_FILE" ]; then
    echo "Error: $ALL_CRDS_FILE not found. Run 'make gen-crds' first." >&2
    exit 1
fi

if ! command -v yq >/dev/null 2>&1; then
    echo "Error: yq is required but not found. Please install yq." >&2
    exit 1
fi

if ! command -v jq >/dev/null 2>&1; then
    echo "Error: jq is required but not found. Please install jq." >&2
    exit 1
fi

# Read production kinds from JSON file
PRODUCTION_KINDS_JQ="[]"
if [ -f "$PRODUCTION_KINDS_FILE" ]; then
    PRODUCTION_KINDS_JQ=$(jq '.' "$PRODUCTION_KINDS_FILE" 2>/dev/null || echo "[]")
fi

if [ "$PRODUCTION_KINDS_JQ" = "[]" ] || [ -z "$PRODUCTION_KINDS_JQ" ]; then
    echo "Warning: No production kinds found. All CRDs will be experimental." >&2
    echo "# Experimental CRDs" > "$EXP_CRDS_FILE"
    cp "$ALL_CRDS_FILE" "$EXP_CRDS_FILE"
    echo "✅ CRDs split successfully (all experimental)"
    exit 0
fi

echo "Production kinds: $(echo "$PRODUCTION_KINDS_JQ" | jq -r '.[]' | tr '\n' ' ')"
echo "Initializing experimental CRDs file..."
echo "# Experimental CRDs" > "$EXP_CRDS_FILE"

prod_count=0
exp_count=0
first_experimental=true

# Convert CRDs to JSON array and process with jq
temp_json=$(mktemp)
yq eval-all 'select(.kind == "CustomResourceDefinition")' "$ALL_CRDS_FILE" -o json | jq -s '.' > "$temp_json"

# Process each CRD using jq - output name, kind, and is_production flag
# Use process substitution (< <(...)) instead of pipe to avoid subshell issues with counters
while IFS='|' read -r crd_name crd_kind is_production; do
    if [ -z "$crd_name" ] || [ -z "$crd_kind" ]; then
        continue
    fi
    
    # Extract the CRD JSON for this specific CRD
    crd_json=$(jq -r --arg name "$crd_name" '.[] | select(.metadata.name == $name and .kind == "CustomResourceDefinition")' "$temp_json")
    
    if [ "$is_production" = "true" ]; then
        # Production CRD: write to individual file as YAML
        output_file="$LEGACY_CRD_BASES_DIR/$crd_name.yaml"
        echo "$crd_json" | jq . | yq eval -P - > "$output_file"
        echo "  Production: $crd_name"
        prod_count=$((prod_count + 1))
        
        # Add to kustomization.yaml if it exists
        if [ -f "$LEGACY_CRD_BASES_DIR/kustomization.yaml" ]; then
            cd "$LEGACY_CRD_BASES_DIR"
            # Add resource if not already present
            if ! yq eval ".resources[] | select(. == \"$crd_name.yaml\")" kustomization.yaml >/dev/null 2>&1; then
                yq eval -i "(.resources // []) += [\"$crd_name.yaml\"] | .resources |= unique" kustomization.yaml
            fi
            cd "$REPO_ROOT"
        fi
    else
        # Experimental CRD: append to experimental file as YAML
        # Add YAML document separator before each CRD (except the first)
        if [ "$first_experimental" = false ]; then
            echo "---" >> "$EXP_CRDS_FILE"
        fi
        echo "$crd_json" | jq . | yq eval -P - >> "$EXP_CRDS_FILE"
        echo "  Experimental: $crd_name"
        exp_count=$((exp_count + 1))
        first_experimental=false
    fi
done < <(jq -r --argjson prod_kinds "$PRODUCTION_KINDS_JQ" '
    .[] |
    select(.kind == "CustomResourceDefinition") |
    .spec.names.kind as $kind |
    .metadata.name as $name |
    ($prod_kinds | index($kind) != null) as $is_prod |
    "\($name)|\($kind)|\($is_prod)"
' "$temp_json")

rm -f "$temp_json"

echo "✅ Split $prod_count production CRDs (individual files) and $exp_count experimental CRDs (aggregated)"
