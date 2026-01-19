# CI

## CI Tests

Atlas Kubernetes Operator testing can be divided into 2 different types of tests:

- **Local Tests**: which includes `unit tests`, `linting` and things like that.
- **Cloud Tests**: which test the operator against Atlas QA cloud performing real resource provisioning actions. These include not both `integration` and `end to end` (`e2e`) tests.

Both tests differ mainly in cost: **Local Tests** are fast, cheap and more reliable because they don't interact with remote resources, while **Cloud Tests** take longer time to run, can be flaky, and consume remote cloud infrastructure. For this reason not all PRs should always run **Cloud Tests**.

Here are the reasons or situations to skip *Cloud Tests*:

- **Draft PRs should not run Cloud Tests** by default.
- **Changes not affecting production code should not need to run Cloud tests** most of the time.
- **External PRs from forked repositories should also not run Cloud Tests** by default, as they should not get access to any credentials without prior inspection.

On other occasions, project maintainers will want to enforce that **Cloud Tests** will run, for example:

- A maintainer sets the `cloud-tests` label on the PR so that the CI tests the code even if production code was not changed, maybe because the CI code did and for this particular change it makes sense to exercise the whole battery test.
- A maintainer sets the `safe-to-test` label on a PR from an external contributor which has been inspected, seems safe and is a candidate for further review and a eventual merge.

Note that in the case of the `safe-to-test` label, cloud tests will only run when the label is first applied (when the action is "labeled"). This ensures that maintainers must re-inspect and re-apply the label if they want to run cloud tests again after changes to the PR.

Additionally, the configuration variable in GitHub `SKIP_OPENSHIFT` can be set to `true` to skip the OpenShift upgrade test, should there be issues or ongoing maintenance in the cluster.

### CI Testing Flow

The workflow [test.yml](../../.github/workflows/test.yml) is the main entry point for the whole test flow.

Most of the times, it will trigger due to a GitHub `pull_request` event, which for PRs from forked repositories, will use read-only credentials and should not have access to cloud resources, so **Cloud Tests** will not be run. For PRs from official maintainers of the repository, this restriction does not apply and the `pull_request` event can run all tests needed.

This workflow also runs on a nightly schedule at midnight on each day of the working week to ensure all tests are run against both the oldest and newest Kubernetes versions supported.

Apart from that, tests can also run on `push` (merges) or on demand by `workflow_dispatch`. Both options are only accessible to official maintainers.

The [test.yml](../../.github/workflows/test.yml) workflow calls the [ci.yml](../../.github/workflows/ci.yml) workflow which runs local tests including:
- Unit tests (via `make unit-test`)
- Linting (via `make all-lints` which includes `golangci`, `shellcheck`, `govulncheck`, `validate-manifests`, `validate-api-docs`, `check-licenses`, `addlicense-check`)

And also calls an special workflow called [cloud-tests-filter.yml](../../.github/workflows/cloud-tests-filter.yml) which is in charge of:

- Checking whether or not the **production code was changed** in this PR.
- Deciding whether or not [cloud tests](../../.github/workflows/cloud-tests.yml) should run** according to the logic we decide.

That workflow also shows relevant CI context values that allow us to debug why the CI took one decision or another, depending of whether the PR is a **draft**, it is **forked or not**, etc.

The [test.yml](../../.github/workflows/test.yml) workflow takes the output from [cloud-tests-filter.yml](../../.github/workflows/cloud-tests-filter.yml) and will ONLY invoke the [cloud tests](../../.github/workflows/cloud-tests.yml) workflow IF `cloud-tests-filter.yml` had decided **Cloud Tests** should run.

### Linting

The lint workflow runs multiple linters via `make all-lints`, including `golangci` (via `make lint`), `shellcheck`, `govulncheck`, `validate-manifests`, `validate-api-docs`, `check-licenses`, and `addlicense-check`.

`golangci` is a tool that makes use of a defined collection of other linters, such as `gosec` and `govet`. The enabled linters (and other configuration) for `golangci` can be seen in [this repo's config file](../../.golangci.yml).

`shellcheck` lints shell scripts in the repo. This is performed with default settings, using [`shellcheck-action`](https://github.com/bewuethr/shellcheck-action). This tool makes use of a regex to find all files within the codebase that have shell scripts that should be assessed.

`govulncheck` checks the Go packages used in the codebase, and flags any that have known vulnerabilities. [`vuln-ignore`](../../vuln-ignore) contains a list of vulnerabilities that we are explicitly ignoring; for use when there is not an available fix, and `govulncheck` is blocking.

#### Cloud tests

The [cloud tests](../../.github/workflows/cloud-tests.yml) workflow is also worth an explanation. It is in charge of running all expensive and slow tests such as:
- [test-int](../../.github/workflows/test-int.yaml)
- [test-e2e](../../.github/workflows/test-e2e.yml)
- [openshift-upgrade-test](../../.github/workflows/openshift-upgrade-test.yml)
- [test-e2e-gov](../../.github/workflows/test-e2e-gov.yml)

Note **Gov e2e tests** are never run on PRs.

The [test-e2e.yml](../../.github/workflows/test-e2e.yml) workflow builds a test image and a bundle before running the tests. It also has to *compute* the version(s) of Kubernetes to test against. The Kubernetes version in PRs is set purposefully to the oldest kubernetes version. On scheduled nightly runs we test on both the latest and oldest supported versions.

##### Kubernetes Version Matrix

The version list selection is done by parameterising the kind image tag within the **strategy** **matrix** at the [test-e2e](../../.github/workflows/test-e2e.yml) workflow. Eg:

```yaml
  prepare-e2e:
    runs-on: ubuntu-latest
    outputs:
      test_matrix: ${{ steps.compute.outputs.matrix }}
    steps:
      - name: Compute K8s matrix/versions for testing
        id: compute
        run: |
          matrix='["v<oldest-k8s-version>-kind"]'
          if [ "${{ github.ref }}" == "refs/heads/main" ]; then
            matrix='["v<oldest-k8s-version>-kind", "v<newest-k8s-version>-kind"]'
          fi
          echo "matrix=${matrix}" >> "${GITHUB_OUTPUT}"
    ...
  e2e:
    needs: [prepare-e2e]
    ...
    strategy:
      fail-fast: false
      matrix:
        k8s: ${{ fromJson(needs.prepare-e2e.outputs.test_matrix) }}
        ...
```

To update Kubernetes versions, check the current supported versions in [test-e2e.yml](../../.github/workflows/test-e2e.yml) and adjust the `matrix` variable in the above workflow accordingly. It may also be necessary to bump the `kind` version and the `kind-action` version in various workflows, see https://github.com/mongodb/mongodb-atlas-kubernetes/pull/2082 as an example.

Additionally, adjust the `ENVTEST_K8S_VERSION` variable in the `Makefile` to match the minimum supported Kubernetes version.

Update the minimum Kubernetes version in the [Atlas Kubernetes CLI repository](https://github.com/mongodb/atlas-cli-plugin-kubernetes] plugin) as well. Here, a Kubernetes cluster is being created for e2e tests programmatically. Bump and adjust the Kubernetes version in its `go.mod` file for the kind version and in the test helper files for the actual Kubernetes version.

Finally, adjust the `com.redhat.openshift.versions` setting in all relevant files to reflect the currently supported OpenShift versions, most notably:
- `scripts/release-redhat-certified.sh`
- `.github/actions/gen-install-scripts/entrypoint.sh`
- `bundle.Dockerfile`

### Test Variants

- **PRs**:
  - Skip cloud tests on non production changes
  - Run e2e tests only on oldest Kubernetes version

- **Merges**:
  - Skip cloud tests on non production changes
  - Run e2e tests in both oldest and newest Kubernetes version

- **Releases & Nightlies**
  - Run ALL test always
  - Run e2e tests in both oldest and newest Kubernetes version

## Release CI

A **release** is initiated by the [`release-image.yml`](../../.github/workflows/release-image.yml) workflow, which takes four inputs: the version to release, the image SHA to be published for the promoted image, the authors for compliance reporting, and the release type. The process is fully automated; the only manual step is approving and merging the release PR. This PR does **not** re-run any tests.

The workflow inputs are:
- `version`: The version to be released in X.Y.Z format (without the `v` prefix)
- `authors`: A comma-separated list of MongoDB email addresses responsible for the release
- `image_sha`: The 7-character Git commit SHA used for the promoted image, or `latest` for the most recent (defaults to `latest`)
- `release_type`: Either `pre-release` (default) or `official-release`. Official releases post to official registries (`mongodb/mongodb-atlas-kubernetes-operator`), while pre-releases post to prerelease registries (`mongodb/mongodb-atlas-kubernetes-operator-prerelease`)

The `image_sha` refers to a previously tested and promoted operator image stored in official prerelease registries (`docker.io`, `quay.io`), traceable to a specific Git commit. Using `latest` here will use the most recent successful image tested. The release workflow uses this image to generate the `release/<version>` directory containing `deploy/`, `helm-charts/`, and `bundle/` folders with all necessary metadata.

The release workflow performs the following automated steps:
- Moves the promoted image from prerelease registries to official release registries (for official releases)
- Creates OpenShift certified images on Quay.io (tagged with `-certified` suffix)
- Signs all released images using PKCS11 signing
- Generates SBOMs for both `linux-amd64` and `linux-arm64` platforms
- Creates SDLC compliance reports
- Generates deployment configurations (bundle, helm charts)
- Bumps the `version.json` file automatically (sets `current` to the released version and `next` to the next minor version)
- Triggers the Helm charts sync workflow to update the [helm-charts repository](https://github.com/mongodb/helm-charts)
- Creates a release PR with all artifacts
- Creates a Git tag of the form `v<version>`
- Publishes a GitHub release with zipped `all-in-one.yml` and SDLC-compliant artifacts (SBOMs and compliance reports)

For more information, see [`release.md`](./release.md).

### Promotion Logic

Operator images are promoted to official prerelease registries after passing all tests. Promotion occurs via:

- Scheduled CI runs on the `main` branch  
- Merges to `main` that modify production code  
- Manual dispatch of `test.yml` with the `promote` input set to `true`

The [`promote-image.yml`](../../.github/workflows/promote-image.yml) workflow runs after all tests, including cloud-based Helm tests, have passed. If successful, it:

- Copies the tested image from `ghcr.io` to `docker.io` and `quay.io`
- Tags the image as `promoted-<commit-sha>` for traceability
- Updates the `promoted-latest` tag to point to this image (only when running on the `main` branch)

For more information, see [`release.md`](./release.md).

### Daily Rebuilds

Daily rebuilds of released images are triggered by [`rebuild-released-images`](../../.github/workflows/rebuild-released-images.yaml), using a dynamically computed list of [supported releases](../../scripts/supported-releases.sh). This workflow:

- Runs on a schedule (weekdays at 1 AM UTC) and can be manually dispatched
- Rebuilds images for all supported releases
- Can target either prerelease or official release registries (via `image_repo` input)
- Tags images with both the release version and a daily tag (`<version>-<date>`)
- Signs and verifies images for supported versions
- Generates SBOMs during the build process

### Helm Charts Synchronization

The [`sync-helm-charts.yaml`](../../.github/workflows/sync-helm-charts.yaml) workflow automatically synchronizes Helm chart changes to the [helm-charts repository](https://github.com/mongodb/helm-charts). It:

- Runs automatically when PRs are merged (excluding dependabot PRs)
- Can be manually dispatched via `workflow_dispatch`
- Verifies if CRDs or RBAC configurations have changed
- Creates a PR in the helm-charts repository with any detected changes
- Is also triggered automatically by the release workflow after a successful release

### Kubernetes Version Monitoring

The [`check-kubernetes-versions.yaml`](../../.github/workflows/check-kubernetes-versions.yaml) workflow monitors Kubernetes and OpenShift version support policy compliance:

- Runs on a weekly schedule (every Monday at midnight UTC)
- Can be manually dispatched
- Checks if version updates are required based on support policies
- Sends Slack alerts when version updates are needed or checks fail

### SBOM Upload to Kondukto

The [`send-sboms.yaml`](../../.github/workflows/send-sboms.yaml) workflow handles SBOM upload and augmentation for compliance:

- Manually triggered via `workflow_dispatch` with a version input
- Downloads SBOM files from the GitHub release for both `linux-amd64` and `linux-arm64` platforms
- Augments both SBOMs with Kondukto scan results using Silkbomb 2.0
- Completes the SSDLC compliance process

For more information, see [`release.md`](./release.md).


## Other Workflows

### Additional Test Workflows

The [test.yml](../../.github/workflows/test.yml) workflow also calls:
- [tests-selectable.yaml](../../.github/workflows/tests-selectable.yaml) - Selectable test suites based on PR labels
- [tests-e2e2.yaml](../../.github/workflows/tests-e2e2.yaml) - Additional e2e test variants

The `test.yml` workflow can be manually dispatched with a `promote` input option (set to `true` or `false`) to control whether image promotion should occur after successful tests. This is useful for manually triggering promotion of images that have passed all tests.
