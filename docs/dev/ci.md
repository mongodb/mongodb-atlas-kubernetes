# CI

## CI Tests

Atlas Kubernetes Operator testing can be divided into tow different types of tests:

- **Local Tests**: which includes `unit tests`, `linting` and things like that.
- **Cloud Tests**: which test the operator against Atlas QA cloud performing real resource provisioning actions. These include not both `integration` and `end to end` (`e2e`) tests.

Both tests differ mainly in cost: **Local Tests** are cheap because they don't consume remote resources, while **Cloud Tests** a longer time to run, can be flaky, and consume remote cloud infrastructure. For this reason not all PRs should always run **Cloud Tests**.

Here are the reasons or situations to skip *Cloud Tests*:

- **Draft PRs should not run Cloud Tests** by default.
- **Changes not affecting production code should not need to run Cloud tests** most of the time.
- **External PRs from forked repositories should also not run Cloud Tests** by default.

On other occasions, project maintainers will want to enforce that **Cloud Tests** will run, for example:

- A maintainer sets the `cloud-tests` label on the PR so that the CI tests the code even if production code was not changed, maybe because the CI code did and for this particular change it makes sense to exercise the whole battery test.
- A maintainer sets the `safe-to-test` label on a PR from an external contributor which has been inspected, seems safe and is a candidate for further review and a eventual merge.

Note that in the case of the `safe-to-test` label, such label is automatically removed by the CI (see workflow `remove-label.yml`) to ensure re-inspection before running **Cloud Tests**.

Additionally, the configuration variable in GitHub `SKIP_OPENSHIFT` can be set to `true` to skip the OpenShift upgrade test, should there be issues or maintenance with the cluster.

### CI Testing Flow

The workflow `test.yml` is the main entry point for the whole test flow.

Most of the times, it will trigger a GitHub `pull_request` event, which for PRs from forked repositories, will only read only credentials and should not be given access to cloud resources. For PRs from maintainers, from the official repository, this restriction does not apply and the `pull_request` event can run all tests needed.

For `safe-to-test` forked PRs, a GitHub `pull_request_target` event is triggered ONLY on the `labeled` event type. This event allows for external PRs to run with write credentials and will be given cloud access. Thus the need to protect access behind a label which maintainers need to set explicitly.

Apart from that, tests can also run on `push` (merges) or on demand by `workflow_dispatch`. Both actions are actions only accessible to official maintainers.

The `test.yml` workflow calls the local tests workflows directly, and then also calls an special workflow called `cloud-tests-filter.yml` which is in charge of:

- Checking whether or not the **production code was changed** in this PR.
- Deciding whether or not the **cloud tests should run** according to the logic we have decided and outlined above.

That workflow also shows relevant CI context values that allow us to debug why the CI took one decision or another, depending of whether the PR is a **draft**, it is **forked or not**, etc.

The `test.yml` takes the output from `cloud-tests-filter.yml` and will ONLY invoke the `cloud-tests.yml` workflow IF `cloud-tests-filter.yml` had decided **Cloud Tests* should run.

## Kubernetes version in CI tests

The kubernetes version in CI tests is set purposefully to the oldest kubernetes version we need to support.

Such version is set by parameterising the kind image tag within the **strategy** **matrix** at the `test-e2e.yml` workflow. Eg:

```yaml
    e2e:
    name: E2E tests
      ...
    strategy:
      fail-fast: false
      matrix:
        k8s: [ "v1.21.1-kind" ] # <K8sGitVersion>-<Platform>
        ...
```
