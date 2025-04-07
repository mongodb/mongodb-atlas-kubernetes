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

Note that in the case of the `safe-to-test` label, such label is automatically removed by the CI (see workflow `remove-label.yml`) to ensure re-inspection before running **Cloud Tests**.

Additionally, the configuration variable in GitHub `SKIP_OPENSHIFT` can be set to `true` to skip the OpenShift upgrade test, should there be issues or ongoing maintenance in the cluster.

### CI Testing Flow

The workflow [test.yml](../../.github/workflows/test.yml) is the main entry point for the whole test flow.

Most of the times, it will trigger due to a GitHub `pull_request` event, which for PRs from forked repositories, will use read-only credentials and should not have access to cloud resources, so **Cloud Tests** will not be run. For PRs from official maintainers of the repository, this restriction does not apply and the `pull_request` event can run all tests needed.

This workflow also runs on a nightly schedule at midnight on each day of the working week to ensure all tests are run against both the oldest and newest Kubernetes versions supported.

Apart from that, tests can also run on `push` (merges) or on demand by `workflow_dispatch`. Both options are only accessible to official maintainers.

The [test.yml](../../.github/workflows/test.yml) workflow calls the local tests workflows directly:
- [lint](../../.github/workflows/lint.yaml)
- [test-unit](../../.github/workflows/test-unit.yml)
- [validate-manifests](../../.github/workflows/validate-manifests.yml)
- [check-licenses](../../.github/workflows/check-licenses.yml)

And also calls an special workflow called [cloud-tests-filter.yml](../../.github/workflows/cloud-tests-filter.yml) which is in charge of:

- Checking whether or not the **production code was changed** in this PR.
- Deciding whether or not [cloud tests](../../.github/workflows/cloud-tests.yml) should run** according to the logic we decide.

That workflow also shows relevant CI context values that allow us to debug why the CI took one decision or another, depending of whether the PR is a **draft**, it is **forked or not**, etc.

The [test.yml](../../.github/workflows/test.yml) workflow takes the output from [cloud-tests-filter.yml](../../.github/workflows/cloud-tests-filter.yml) and will ONLY invoke the [cloud tests](../../.github/workflows/cloud-tests.yml) workflow IF `cloud-tests-filter.yml` had decided **Cloud Tests** should run.

### Linting

The lint workflow runs three seperate linters; `golangci` (via `make lint`), `shellcheck`, and `govulncheck`.

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
  compute:
    environment: test
    name: "Compute test matrix"
    runs-on: ubuntu-latest
    outputs:
      test_matrix: ${{ steps.test.outputs.matrix }}
    steps:
      - id: test
        name: Compute Test Matrix
        run: |
          # Note the use of external single quotes to allow for double quotes at inline YAML array
          matrix='["v1.27.1-kind"]'
          if [ "${{ github.ref }}" == "refs/heads/main" ];then
            matrix='["v1.27.1-kind", "v1.29.2-kind"]'
          fi
          echo "matrix=${matrix}" >> "${GITHUB_OUTPUT}"
          cat "${GITHUB_OUTPUT}"
    ...
  e2e:
    name: E2E tests
    ...
    strategy:
      fail-fast: false
      matrix:
        k8s: ${{fromJson(needs.compute.outputs.test_matrix)}}
        ...
```

Adjust the `matrix` variable in the above workflow to match the desired Kubernetes versions. It migh also be necessary to bump the `kind` version and the `kind-action` version in various workflows, see https://github.com/mongodb/mongodb-atlas-kubernetes/pull/2082 as an example.

Additionally, adjust the `ENVTEST_K8S_VERSION` variable in the `Makefile` as well.

Adjust the minimum Kubernetes version ("1.27.1" in the above example) in the [Atlas Kubernetes CLI repository](https://github.com/mongodb/atlas-cli-plugin-kubernetes] plugin) as well. Here, a Kubernetes cluster is being created for e2e tests programmatically. Bump and adjust the Kubernetes version in its `go.mod` file: https://github.com/mongodb/atlas-cli-plugin-kubernetes/blob/d34c4b18930b0cd77dc6013d52669161edb224d5/go.mod#L32 for the kind version and https://github.com/mongodb/atlas-cli-plugin-kubernetes/blob/d5b2610dd50e312e315b63d1bfd0d7dde244b262/test/e2e/operator_helper_test.go#L91-L98 for the actual Kubernetes version.

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

A **release PR** is prepared by the [release-branch.yml](../../.github/workflows/release-branch.yml) workflow.

Once the PR is approved and **closed**, the [tag.yml](../../.github/workflows/release-branch.yml) workflow takes over:
- The branch `release/X.Y.Z` would be tagged as release `vX.Y.Z`.
  - A pre-release branch `prerelease/X.Y.Z-...` would become `vX.Y.Z-...`.
- Then does a workflow call to [release-post-merge.yml](../../.github/workflows/release-opost-merge.yml).

Note `release-post-merge` can also be triggered by tagging a release or pre-release manually.

After a release is published, daily rebuilds are run by [rebuild-released-images](../../.github/workflows/rebuild-released-images.yaml):
  - The list of [supported releases is computed dynamically](../../scripts/supported-releases.sh).

## Other Workflows

### Update dependabot PR Licenses

Workflow [update-licenses.yml](../../.github/workflows/update-licenses.yml) runs to patch `dependabot`'s go module update PRs.

When `dependabot` updates go modules, dependencies change and license dependencies might also change. But `dependabot` does not know how to update your code when dependencies change. This workflow is triggered on `dependabot` PRs, runs `make recompute-licenses` and patches the PR as needed.
