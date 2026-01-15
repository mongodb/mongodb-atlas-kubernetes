# Updating Kubernetes Versions

When the Kubernetes versions check alert triggers, it indicates that the supported Kubernetes version range in `kubernetes-versions.json` needs to be updated. This document outlines the steps required to update the test infrastructure accordingly.

For the original support policy, see [Kubernetes Operators Versioning](https://wiki.corp.mongodb.com/spaces/MMS/pages/193108905/Kubernetes+Operators+Versioning).

## Overview

The Kubernetes versions check runs weekly and validates that our supported version range aligns with the support policy. When an update is required, you need to:

1. Update the test version matrices in GitHub Actions workflows (if Kubernetes versions changed)
2. Coordinate with the wider MCK team to upgrade the OpenShift cluster (if OpenShift version changed)
3. Update `kubernetes-versions.json` to reflect the actual supported versions after infrastructure is ready

## Step 1: Update Test Version Matrices

The test workflows use Kubernetes version matrices that need to be updated to match the new supported range. Update the following files:

### `.github/workflows/test-e2e.yml`

Update the `Compute K8s matrix/versions for testing` step:

```yaml
- name: Compute K8s matrix/versions for testing
  id: compute
  run: |
    matrix='["v<MIN_VERSION>-kind"]'  # Update MIN_VERSION
    if [ "${{ github.ref }}" == "refs/heads/main" ]; then
      matrix='["v<MIN_VERSION>-kind", "v<MAX_VERSION>-kind"]'  # Update both
    fi
    echo "matrix=${matrix}" >> "${GITHUB_OUTPUT}"
```

### `.github/workflows/tests-selectable.yaml`

Update the `Compute test matrix for k8s versions` step in the `compute` job:

```yaml
- id: test
  name: Compute test matrix for k8s versions
  run: |
    matrix='["v<MIN_VERSION>-kind"]'  # Update MIN_VERSION
    if [ "${{ github.ref }}" == "refs/heads/main" ];then
      matrix='["v<MIN_VERSION>-kind", "v<MAX_VERSION>-kind"]'  # Update both
    fi
    echo "matrix=${matrix}" >> "${GITHUB_OUTPUT}"
```

### `.github/workflows/tests-e2e2.yaml`

Update the `Compute test matrix for k8s versions` step in the `compute` job:

```yaml
- id: test
  name: Compute test matrix for k8s versions
  run: |
    matrix='["v<MIN_VERSION>-kind"]'  # Update MIN_VERSION
    if [ "${{ github.ref }}" == "refs/heads/main" ];then
      matrix='["v<MIN_VERSION>-kind", "v<MAX_VERSION>-kind"]'  # Update both
      echo "Nightly runs oldest and newest Kubernetes supported versions"
    fi
    echo "matrix=${matrix}" >> "${GITHUB_OUTPUT}"
```

**Note:** Replace `<MIN_VERSION>` and `<MAX_VERSION>` with the actual version numbers from `kubernetes-versions.json` (e.g., `1.32`, `1.34`). The format should be `v<version>-kind` (e.g., `v1.32.8-kind`, `v1.34.0-kind`).

**Tip:** If only Kubernetes versions changed (not OpenShift), steps 1 and 3 can be done in the same PR - update the workflow matrices and `kubernetes-versions.json` together.

## Step 2: Coordinate OpenShift Cluster Upgrade

The OpenShift version specified in `kubernetes-versions.json` needs to be supported by our test infrastructure. **This requires coordination with the wider MCK (MongoDB Controllers for Kubernetes) team** before making changes.

**Important:** The OpenShift cluster upgrade happens out of the repo. Ensure to update `kubernetes-versions.json` as soon as the upgrade of the cluster is completed.

For detailed instructions on updating the ROSA cluster, see [Updating the ROSA cluster](release.md#updating-the-rosa-cluster) in the release documentation.

## Step 3: Update kubernetes-versions.json

After completing steps 1 and/or 2, update `kubernetes-versions.json` to reflect the actual supported versions. This file should be kept in sync with the reality of what versions are actually tested and supported in the infrastructure.

**Important:** 
- `kubernetes-versions.json` should reflect what is actually supported and tested, not what we plan to support
- For OpenShift upgrades, this file is what registers the new version (the upgrade itself happens offline)
- If only Kubernetes versions changed (not OpenShift), this step can be combined with step 1 in the same PR

### Process

1. **Verify infrastructure readiness**:
   - Test matrices have been updated (if Kubernetes versions changed)
   - OpenShift cluster upgrade is complete (if OpenShift version changed)
2. **Update the file** with the new supported versions:
   ```json
   {
     "kubernetes": {
       "min": "<MIN_VERSION>",
       "max": "<MAX_VERSION>"
     },
     "openshift": "<OPENSHIFT_VERSION>"
   }
   ```
3. **Verify the changes**:
   - Run the check script locally: `./scripts/check-kube-versions.sh`
   - Verify the version matrices match the range in `kubernetes-versions.json`
   - Ensure all three workflow files are updated consistently (if Kubernetes versions changed)
4. **Test and commit**:
   - Test the workflows on a branch before merging to main
   - Commit the changes (workflow updates and `kubernetes-versions.json` can be in the same PR if only Kubernetes versions changed)

## Related Files

- `kubernetes-versions.json` - Source of truth for supported versions
- `scripts/check-kube-versions.sh` - Version check script
- `.github/workflows/check-kubernetes-versions.yaml` - Weekly check workflow

