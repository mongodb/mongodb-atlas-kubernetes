# CI

## CI Tests

Atlas Kubernetes Operator testing can be divided into tow different types of tests:

- **Local Tests**: which includes `unit tests`, `linting` and things like that.
- **Cloud Tests**: which test the operator against Atlas QA cloud performing real resource provisioning actions. These include not both `integration` and `end to end` (`e2e`) tests.

Both tests differ mainly in cost: **Local Tests** are cheap because they run in a few minutes and are reliable, while **Cloud Tests** take about an hour to run, can be flaky and require access to certain secrets to manage cloud infrastructure. For this reason not all PRs should always run **Cloud Tests**.

Here are the reasons or situations to skip *Cloud Tests*:

- **Draft PRs should not run Cloud Tests** by default.
- **Changes not affecting production code should not need to run Cloud tests** most of the time.
- **External PRs from forked repositories should also not run Cloud Tests** by default, as they will get access to cloud resources without any vet the changes to check whether or not such access is reasonable and safe.

On other occasions, project maintainers will want to enforce that **Cloud Tests** will run, for example:

- A maintainer sets the `cloud-tests` label on the PR so that the CI tests the code even if production code was not changed, maybe because the CI code did and for this particular change it makes sense to exercise the whole battery test.
- A maintainer sets the `safe-to-test` label on a PR from an external contributor which has been inspected, seems safe and is a candidate for further review and a eventual merge.

Note that in the case of the `safe-to-test` label, such label is automatically removed by the CI (see workflow `remove-label.yml`) to ensure re-inspection before running **Cloud Tests**.

### CI Testing Flow

The workflow `test.yml` is the main entry point for the whole test flow.

Most of the times, it will trigger a GitHub `pull_request` event, which for PRS from forked repositories, will only have read only credentials and should not be given access to cloud resources. For PRs from maintainers, from the official repository, this restriction does not apply and `pull_request` event can run all tests needed.

For `safe-to-test` forked PRs, a GitHub `pull_request_target` event is triggered ONLY on the `labeled` event type. This event allows for external PRs to run with write credentials and will be given cloud access. This the need to protect access behind a label maintainers need to explicitly set.

Apart from that, tests can also run on `push` (merges) or on demand by `workflow_dispatch`. Both actions are actions only accessible to official maintainers.

The `test.yml` workflow calls the local tests workflows directly, and then also calls an special workflow called `cloud-tests-filter.yml` which is in charge of:

- Checking whether or not the **production code was changed** in this PR.
- Deciding whether or not the **cloud tests should run** according to the logic we have decided and outlined above.

That workflow also shows relevant CI context values that allow us to debug why the CI took one decision or another, depending of whether the PR is a **draft**, it is **forked or not**, etc.

The `test.yml` takes the output from `cloud-tests-filter.yml` and will ONLY invoke the `cloud-tests.yml` workflow IF `cloud-tests-filter.yml` had decided **Cloud Tests* should run.
