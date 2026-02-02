# OpenShift Image Certification

## Overview

OpenShift certification validates container images for Red Hat's OpenShift platform. Certified images can be listed in Red Hat's catalog, used in production with Red Hat support, and distributed through official channels.

The `certify-openshift-images.sh` script uses Red Hat's **preflight** tool to validate and optionally certify images. **Note**: Registry login must be performed before running the script.

## Two Modes

- **Check Only (`SUBMIT=false`)** - NON-DESTRUCTIVE: Runs validation locally, reports results, does NOT write to Red Hat. Safe for testing.
- **Submit (`SUBMIT=true`)** - DESTRUCTIVE ⚠️: Writes certification results to Red Hat's Pyxis API, makes certification official. Point of no return.

## What Preflight Checks

- Image security vulnerabilities (CVE scanning)
- Image structure and metadata compliance
- Required labels and annotations
- Container best practices
- OpenShift compatibility

## Testing

**Required environment variables:** `REGISTRY`, `REPOSITORY`, `VERSION`, `RHCC_TOKEN`, `RHCC_PROJECT`

**Use Make targets to simplify testing:**

**Pre-certification testing:**
```bash
docker login docker.io -u <username> -p <password>
PROMOTED_TAG=promoted-latest \
RHCC_TOKEN=$RH_CERTIFICATION_PYXIS_API_TOKEN \
RHCC_PROJECT=$RH_CERTIFICATION_OSPID \
make pre-cert-sandbox certify-openshift-images
```

**Certification testing (check only):**
```bash
docker login quay.io -u <username> -p <password>
VERSION=2.13.0-certified \
RHCC_TOKEN=$RH_CERTIFICATION_PYXIS_API_TOKEN \
RHCC_PROJECT=$RH_CERTIFICATION_OSPID \
make cert-sandbox certify-openshift-images
```

**Production submit (⚠️ destructive):**
```bash
docker login quay.io -u <username> -p <password>
SUBMIT=true \
REGISTRY=quay.io \
REPOSITORY=mongodb/mongodb-atlas-kubernetes-operator \
VERSION=2.13.0-certified \
RHCC_TOKEN=$RH_CERTIFICATION_PYXIS_API_TOKEN \
RHCC_PROJECT=$RH_CERTIFICATION_OSPID \
make certify-openshift-images
```

## In Release Workflow

In `.github/workflows/release-image.yml`:
- **Pre-releases**: `SUBMIT=false` (safe, no writes to Red Hat)
- **Official releases**: `SUBMIT=true` (destructive, writes to Red Hat)

## Best Practices

1. Always test with `SUBMIT=false` first
2. Only use `SUBMIT=true` for official releases
3. Verify the image is correct before submitting
4. This is the final step in the release process for a reason

## Related Documentation

- [Red Hat Preflight](https://github.com/redhat-openshift-ecosystem/openshift-preflight)
- [Red Hat Container Certification](https://connect.redhat.com/en/products/red-hat-openshift-container-platform-certification)
- [Pyxis API](https://pyxis.engineering.redhat.com/)

