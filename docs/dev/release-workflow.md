# Release Workflow

This document describes the release workflow steps, highlights operations with external side effects, and explains how to recover from failures.

## Workflow Steps

### Preparation (No Side Effects)
- Resolve commit SHA from image
- Validate version
- Generate artifacts (SBOMs, SDLC reports, licenses)

### External Artifacts Created ⚠️
1. **Release PR Creation** ⚠️
   - Creates branch `new-release/v{VERSION}` and pushes to GitHub
   - Creates draft PR
   - Updates `version.json` and `helm-charts/` on branch

2. **Git Tag Creation** ⚠️
   - Creates and pushes tag `v{VERSION}` to GitHub
   - Tag is publicly visible

3. **GitHub Release Creation** ⚠️
   - Creates draft GitHub release
   - Attaches artifacts (tar.gz, SBOMs, SDLC report)
   - **Note**: Release remains draft until manually published

### Safety Check
4. **Pre-validation Certification (Dry-run)**
   - Validates image will pass OpenShift certification
   - Runs against prerelease image before publishing
   - Fails workflow if issues detected (prevents wasted work)

### Image Publishing 🚨 (Point of No Return)
5. **Image Publishing** 🚨
   - Pushes images to Docker Hub and Quay.io release registries
   - Signs all images with cosign
   - **Images become publicly available** (cannot be easily removed)
   - Processes targets atomically: Docker → Quay → Quay certified

6. **OpenShift Certification**
   - Submits certification to Red Hat (official releases only)
   - Creates certified image tag on Quay.io

## Recovery from Failures

**TL;DR: Just retry the workflow. It's idempotent.**

The workflow is designed to be safely re-runnable. All steps check for existing artifacts and skip or update them as needed.

### If Workflow Fails

1. **Check the failure point** in the workflow logs
2. **Re-run the workflow** with the same inputs
   - The workflow will:
     - Skip steps that already completed
     - Update existing artifacts (PR, tag, release) if needed
     - Continue from where it left off
3. **Verify completion** - Check that all steps succeeded

### Special Cases

- **Partial image publishing**: If images were partially published, retry will complete the missing targets (each target is atomic)
- **Draft release**: The release will remain as draft during retries (safeguard ensures this)
- **Certification failure**: Can be re-run independently after images are published

### Manual Cleanup (Rarely Needed)

Only needed if you want to start completely fresh:

1. Delete the release PR and branch
2. Delete the git tag: `git push origin --delete v{VERSION}`
3. Delete the GitHub release (if it exists)
4. Re-run the workflow

**Note**: Manual cleanup is usually unnecessary - just retry the workflow.

## Idempotency

All critical steps are idempotent:
- ✅ PR creation checks if PR exists
- ✅ Tag creation uses force update
- ✅ Release creation updates existing releases
- ✅ Image publishing checks for existing images
- ✅ Image signing checks for existing signatures
- ✅ Certification can be re-run

## Important Notes

- **Draft releases**: Releases are created as drafts and remain drafts during retries
- **Image availability**: Images are published after artifacts are created (ordering limitation)
- **Certification**: Pre-validation runs before publishing to catch issues early
- **Retry-safe**: Workflow can be safely re-run from any point

## Local Testing

### Testing Prepare Released Branch

The `prepare-released-branch` step can be tested locally without a sandbox since it only creates local directories:

```bash
# Test with a commit SHA (e.g., from a previous release)
# rm -rf -rf released-branch .ci-tooling ## shoudl you need to retest
make prepare-released-branch COMMIT_SHA=<commit-sha>
```

This will create `released-branch/` directory with the checked-out commit and `.ci-tooling/` directory with CI tooling.

### Testing Prepare Release PR

The `prepare-release-pr` step combines `prepare-released-branch`, `build-release-pr`, and `create-release-pr`. It requires setting up some variables, as it has external side effects (pushes branches and creates PRs on GitHub). Use `release-sandbox.sh` which provides safe defaults for PR testing variables:

```bash
# Prepare released branch, build artifacts, and create PR in one step
# release-sandbox.sh provides safe defaults for AUTHORS, IMAGE_URL, OPERATOR_REGISTRY, etc.
GITHUB_TOKEN=$(gh auth token) \
COMMIT_SHA="${COMMIT_SHA}" \
RELEASE_TAG="v${VERSION}" VERSION="${VERSION}" \
./scripts/release-sandbox.sh make prepare-release-pr
```

Or run the steps separately:

```bash
# First prepare the released branch
export COMMIT_SHA=<commit-sha>
./scripts/release-sandbox.sh make prepare-released-branch COMMIT_SHA=${COMMIT_SHA}

# Then build release artifacts (release-sandbox.sh provides defaults for all vars)
export VERSION=2.14.0-test
./scripts/release-sandbox.sh make build-release-pr VERSION=${VERSION}

# Finally create the PR (will create branch and PR in test repo)
export RELEASE_TAG=v2.14.0-test
GITHUB_TOKEN=$(gh auth token) \
./scripts/release-sandbox.sh make create-release-pr RELEASE_TAG=${RELEASE_TAG} VERSION=${VERSION}
```

**Note:** This will create real branches and PRs in the specified repository. Use a test repository or fork for testing. To retest, delete the PR and branch manually, or the script will update an existing PR if the branch already exists. The `release-sandbox.sh` script provides safe defaults for PR testing variables (AUTHORS, IMAGE_URL, OPERATOR_REGISTRY, RELEASED_OPERATOR_IMAGE, IMAGE_TAG, SKIP_SIGNATURE_VERIFY).

### Testing Release Image Operations

To test release image operations locally, use the `release-sandbox.sh` script.

You might need to do `gh auth login --scopes write:packages` beforehand.

```bash
# Test image publishing, specify a custom registry
TMPDIR=./tmp \
GH_TOKEN=$(gh auth token) \
PKCS11_URI="${PKCS11_URI}" \
GRS_USERNAME="${GRS_USERNAME}" \
GRS_PASSWORD="${GRS_PASSWORD}" \
SANDBOX_REGISTRY="${MY_REGISTRY}" \
./scripts/release-sandbox.sh make push-release-images
```

The script sets up sandbox environment variables that redirect image pushes to a test registry, allowing you to validate release logic without affecting production registries.

