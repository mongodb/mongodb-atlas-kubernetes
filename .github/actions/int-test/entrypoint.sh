#!/bin/bash


export ATLAS_ORG_ID="${INPUT_ATLAS_ORG_ID}"
export ATLAS_PUBLIC_KEY="${INPUT_ATLAS_PUBLIC_KEY}"
export ATLAS_PRIVATE_KEY="${INPUT_ATLAS_PRIVATE_KEY}"

ginkgo -v ./test/int/*