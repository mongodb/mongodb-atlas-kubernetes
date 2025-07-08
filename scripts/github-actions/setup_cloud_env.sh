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

export MCLI_PUBLIC_API_KEY="${MCLI_PUBLIC_API_KEY:-${{ secrets.ATLAS_PUBLIC_KEY }}}"
export MCLI_PRIVATE_API_KEY="${MCLI_PRIVATE_API_KEY:-${{ secrets.ATLAS_PRIVATE_KEY }}}"
export MCLI_ORG_ID="${MCLI_ORG_ID:-${{ secrets.ATLAS_ORG_ID }}}"
export MCLI_OPS_MANAGER_URL="https://cloud-qa.mongodb.com/"
export IMAGE_PULL_SECRET_REGISTRY="ghcr.io"
export IMAGE_PULL_SECRET_USERNAME="$"
export IMAGE_PULL_SECRET_PASSWORD="${{ secrets.GITHUB_TOKEN }}"

export AWS_ACCESS_KEY_ID="${AWS_ACCESS_KEY_ID:-${{ secrets.AWS_ACCESS_KEY_ID }}}"
export AWS_ACCOUNT_ARN_LIST="${AWS_ACCOUNT_ARN_LIST:-${{ secrets.AWS_ACCOUNT_ARN_LIST }}}"
export AWS_SECRET_ACCESS_KEY="${AWS_SECRET_ACCESS_KEY:-${{ secrets.AWS_SECRET_ACCESS_KEY }}}"

export AZURE_CLIENT_ID="${AZURE_CLIENT_ID:-${{ secrets.AZURE_CLIENT_ID }}}"
export AZURE_TENANT_ID="${AZURE_TENANT_ID:-${{ secrets.AZURE_TENANT_ID }}}"
export AZURE_CLIENT_SECRET="${AZURE_CLIENT_SECRET:-${{ secrets.AZURE_CLIENT_SECRET }}}"
export AZURE_SUBSCRIPTION_ID="${AZURE_SUBSCRIPTION_ID:-${{ secrets.AZURE_SUBSCRIPTION_ID }}}"

export GCP_SA_CRED="${GCP_SA_CRED:-${{ secrets.GCP_SA_CRED }}}"
export DATADOG_KEY="${DATADOG_KEY:-${{ secrets.DATADOG_KEY }}}"
export PAGER_DUTY_SERVICE_KEY="${PAGER_DUTY_SERVICE_KEY:-${{ secrets.PAGER_DUTY_SERVICE_KEY }}}"
