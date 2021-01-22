#!/bin/bash


export ATLAS_ORG_ID="${INPUT_ATLAS_ORG_ID}"
export ATLAS_PUBLIC_KEY="${INPUT_ATLAS_PUBLIC_KEY}"
export ATLAS_PRIVATE_KEY="${INPUT_ATLAS_PRIVATE_KEY}"
# otherwise we may get strange "Detected Programmatic Focus - setting exit status to 197"
export GINKGO_EDITOR_INTEGRATION="true"

ginkgo -v ./test/int