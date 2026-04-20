# Atlas Operator Release Instructions

For the various PRs involved, seek at least someone else to approve. In case of doubts, engage the team member(s) who might be able to clarify and seek their review as well.

## Prerequisites

To get PRs to be auto-committed for RedHat community & Openshift you need to make sure you are listed in the [team members ci.yaml list for community-operators](https://github.com/k8s-operatorhub/community-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml) and [team members ci.yaml list for community-operators-prod](https://github.com/redhat-openshift-ecosystem/community-operators-prod/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml).

This is not required for [Certified Operators](https://github.com/redhat-openshift-ecosystem/certified-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml).

Finally, make sure you have a "RedHat Connect" account and are a [team member with org administrator role in the team list](https://connect.redhat.com/account/team-members).

> [!CAUTION]
> Ensure that the commit you are releasing (the most recent change included in the image) did not contain changes to 
> GitHub workflows. This causes the release workflows to error.

### Tools

Most tools are automatically installed for you. Most of them are Go binaries and use `go install`. There are a few that might cause issues and you might want to pre-install manually:

- [devbox](https://www.jetify.com/devbox) to be able to enter a sandbox development environment that includes necessary tools for the release process.
- [Docker](https://www.docker.com/) to be able to deal with containers.

## Kubernetes version updates

The Kubernetes Version testing matrix is automatically checked via CI, Every week an automated workflow:
- Monitors Kubernetes and OpenShift version support policy compliance
- Alerts the team via Slack when version updates are required or the check fails

**No manual action is required** unless you receive an alert or the automated check fails. If an alert is received or the check fails, please refer to the [Updating Kubernetes Versions documentation](update-kubernetes-version.md)

Before starting a release, verify that the [Check Kubernetes Versions workflow](.github/workflows/check-kubernetes-versions.yaml) ran successfully recently to ensure versions are up to date.

## A Note on Versioning with `version.json`

The `version.json` file is the definitive **source of truth** for managing software versions in this project.

* **Source of Truth**: The file contains two primary fields: **`current`**, which reflects the latest stable version that has been released, and **`next`**, which designates the version targeted for the upcoming release.

* **Automatic Updates for next Releases**: After a successful release, the CI pipeline will **automatically** update `version.json`. The `current` field is set to the version that was just released, and the `next` field is incremented to the next minor version (e.g., `2.11.0` would become `2.12.0`, or `2.11.4` would become `2.12.0`).

* **Manual Updates for Patch Releases**: Creating a **patch release** (e.g., `v2.10.1`) is a deliberate exception to the automated process. To prepare for a patch, you must **manually** update the `next` field to the exact patch version via a Pull Request (PR). This manual step ensures the release workflow targets the specific patch instead of accidentally creating the next minor release. In case of modifying the `version.json`, you should trigger the `Test` workflow with the `promote image` option set to `true` so that it builds a new image with the updated version from `version.json` file BEFORE triggering a `Release image` workflow.

## Release Notes

- Create a draft of the release notes in a Google Document and share with Product and the Docs team.
  - In confluence, look for the `AKO Release Internal details` page for more details.
- Ensure as well that supporting documents for new features are in review.
- Wait for approval of the release notes and availability of the associated documents.

**DO NOT** start the release process until the release notes are approved and associated documentation is, at least, in review state. Always seek explicit approval by Product and/or Management.

The reason for this preparatory step is to avoid customers getting new or breaking changes before their supporting documentation.

## Create the Release

Once release notes and documentation are approved, trigger the [`release-image.yml`](../../.github/workflows/release-image.yml) workflow.

You will be prompted to enter:

| Input          | Description                                                                                         | Required | Default       | Example                               |
|----------------|-----------------------------------------------------------------------------------------------------|----------|---------------|---------------------------------------|
| `version`      | The version to be released in the similar X.Y.Z format, without the `v` prefix                    | Yes      | None          | 2.11.0                                |
| `authors`      | A comma-separated list of MongoDB email addresses responsible for the release                       | Yes      | None          | `alice@mongodb.com,bob@mongodb.com`   |
| `image_sha`    | The 7-character Git commit SHA used for the promoted image, or `'latest'` for the most recent       | No       | `latest`      | `3e79a3f`                             |
| `release_type` | Either `pre-release` or `official-release`. Official releases post to official registries           | No       | `pre-release` | `official-release`                    |

The input `version` is a **safety check** to ensure the intent is to release the same version that was already tagged in the image already tested and ready for release. The workflow will fail if the input does not match the expected version from `version.json`.

The input `authors` must be filled out every time you trigger the release workflow. The `image_sha` is optional and defaults to `latest` if left empty.

The `image_sha` corresponds exactly to the 7-character Git commit SHA used to build the operator image; for example, `image_sha: 3e79a3f` means the image was built from Git commit `3e79a3f`. Using `latest` as the `image_sha` means the workflow will release the most recently promoted and tested operator image—not necessarily the latest Git commit—and when `latest` is used, the workflow will echo the corresponding Git commit during the internal steps so the user knows exactly which source is being released.

The `release_type` input determines which Docker registries the images are published to:
- `pre-release` (default): Images are published to `mongodb/mongodb-atlas-kubernetes-operator-prerelease` registries
- `official-release`: Images are published to `mongodb/mongodb-atlas-kubernetes-operator` registries (official release registries)

### Example Release Input

```yaml
version: 2.10.2
authors: alice@mongodb.com,bob@mongodb.com
image_sha: 3e79a3f
```

or

```yaml
version: 2.10.2
authors: alice@mongodb.com,bob@mongodb.com
image_sha: latest
```

### What Happens Next

The release workflow performs automated steps including artifact generation, PR creation, image publishing, and certification. 

**See [Release Workflow](release-workflow.md) for detailed step-by-step information, including which steps have external side effects and how to recover from failures.**

The only manual step is to **review and merge** the release PR. This PR does **not** re-run any of the expensive tests on cloud-qa.

> [!NOTE]
> **Retry Capability**: The workflow is fully idempotent and can be safely re-run if it fails. If a release fails, simply re-run the workflow with the same inputs - it will continue from where it left off.

**Note:** this directory-based approach avoids merge conflicts entirely. Because each release introduces a clean, isolated `release/<version>` folder, it can be merged directly into `main` without conflicting with prior or future releases. This enables a linear and conflict-free release history while maintaining clear traceability for each version.

---

## Image Promotion

The `image_sha` used in a release must already be tested and promoted via CI. Promotion can occur in one of two ways:

- A scheduled CI run on the `main` branch
- A manual dispatch of the `tests.yml` workflow with the `promote` flag enabled

### How Promotion Works

During promotion, the operator image used in Helm-based E2E tests is first built and published as a dummy image in `ghcr.io`. Once **all** tests—including the cloud-based Helm scenarios—complete successfully, the [`promote-image.yml`](../../.github/workflows/promote-image.yml) workflow is triggered.

This workflow:

- Moves the image from `ghcr.io` to official prerelease registries in `docker.io` and `quay.io`
- Tags the image in the official prerelease registires as:
  - `promoted-<git_sha>` — uniquely maps the image to the source Git commit
  - `promoted-latest` — always points to the most recent image that passed all tests

The `promoted-<git_sha>` builds the one-to-one correspondence between the 7-character Git commit and the `image_sha`. For the correspondence between the 7-character Git commit and `image_sha: latest`, we internally store a label within the image `promoted-latest` that has the exact git commit used for that image. Moreover, the `promoted-latest` tag is only updated by events that run on the main branch—whether triggered by a schedule or a workflow dispatch. Manual promotions on any other branch will never overwrite this tag.

One can find promoted images by browsing the prerelease Docker registries at:

- Docker Hub: `mongodb/mongodb-atlas-kubernetes-prerelease`
- Quay.io: `mongodb/mongodb-atlas-kubernetes-prerelease`

**Note:** When releasing, you omit the `promoted-` prefix and specify only the image SHA or `latest`. The `promoted-` prefix is used internally to organize images in the registries.

### Best Practice

Releases should generally use `latest` as the `image_sha`. This ensures that you are releasing the most recently tested and CI-verified image.

## Edit the Release Notes and publish the release

Follow the format described in the [release-notes-template.md](../release-notes/release-notes-template.md) file.
Paste the release notes content approved before the release was started.
Once the image is out, publish the release notes draft as soon as possible.

## Upload SBOMs to Kondukto

> [!IMPORTANT]
> The GitHub release notes must have been **published** before the SBOMs can be sent to Kondukto.

The SBOM upload and augmentation process is automated via the [Send SBOMs to Kondukto](../../.github/workflows/send-sboms.yaml) GitHub workflow.

1. Navigate to the [Actions tab](https://github.com/mongodb/mongodb-atlas-kubernetes/actions) in the GitHub repository.
2. Select the "Send SBOMs to Kondukto" workflow from the workflow list.
3. Click "Run workflow" and provide:
   - **version**: The version to send SBOMs for (e.g., `2.11.0`), without the `v` prefix.
4. The workflow will:
   - Download the SBOM files from the GitHub release for both `linux-amd64` and `linux-arm64` platforms
   - Augment both SBOMs with Kondukto scan results using Silkbomb 2.0
   - Complete the SSDLC compliance process

**Note**: SBOMs are automatically generated during the release process, but the Kondukto upload must be triggered manually after the release is published.

## Synchronize configuration changes with the Helm Charts

The [helm-charts syncer job](https://github.com/mongodb/helm-charts/actions/workflows/atlas-operator-chart-sync.yaml) automatically creates a Pull Request against the [helm-charts repository](https://github.com/mongodb/helm-charts) for every PR made against this AKO repository. The syncer does **not** create a PR until the AKO release PR is merged.

After merging the AKO release PR:

1. Go to the [helm-charts repo Pull Requests](https://github.com/mongodb/helm-charts/pulls) and locate the PR automatically created by the helm-charts syncer job. It is named "Release Atlas Operator x.y.z."

2. The PR will update two Helm charts:
   * [atlas-operator-crds](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator-crds)
   * [atlas-operator](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator)

3. Merge the PR - the chart will get released automatically.

## Create Pull Requests to publish OLM bundles

All bundles/package manifests for Operators for operatorhub.io reside in the following repositories:
* https://github.com/k8s-operatorhub/community-operators - Kubernetes Operators that appear on [OperatorHub.io](https://operatorhub.io/)
* https://github.com/redhat-openshift-ecosystem/community-operators-prod - Kubernetes Operators that appear on [OpenShift](https://openshift.com/) and [OKD](https://www.okd.io/)
* https://github.com/redhat-openshift-ecosystem/certified-operators - Red Hat certified Kubernetes Operators

All 3 PRs for those repos can be pushed from workflow [Push release PRs to RedHat](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/release-rh.yml). This can also be run from the CLI using make, please look [at the code](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/release-rh.yml) for more details.

The workflow needs 4 parameters:
- The `version` to release, which is **required**.
- Whether or not this is a `dryrun`. Set to `false` to do the actual release.
- The `author` name, which is **required**. This should be the author's username [as registered in the RedHat repos like this](https://raw.githubusercontent.com/k8s-operatorhub/community-operators/refs/heads/main/operators/mongodb-atlas-kubernetes/ci.yaml). See [Prerequisites](#prerequisites).

Note when the dryrun is `true` the workflow does everything, except the `git push` is dry run, which should test access credentials, but not make the push really happen.

The `author` and `email` parameters are used to configure the git user identity for commits made during the release process.

Once the workflow ends sucessfully, please go to the projects PR tabs and complete the PRs to review at:
* https://github.com/k8s-operatorhub/community-operators/pulls
* https://github.com/redhat-openshift-ecosystem/community-operators-prod/pulls
* https://github.com/redhat-openshift-ecosystem/certified-operators/pulls

The job log should end with direct links for you to create the PRs, they will look like this:
```
https://github.com/mongodb-forks/community-operators/pull/new/mongodb-atlas-operator-community-${version}
https://github.com/mongodb-forks/community-operators-prod/pull/new/mongodb-atlas-operator-community-${version}
https://github.com/mongodb-forks/certified-operators/pull/new/mongodb-atlas-kubernetes-operator-${version}
```

# Post install hook release

If changes have been made to the post install hook (mongodb-atlas-kubernetes/cmd/post-install/main.go).
You must also release this image. Run the "Release Post Install Hook" workflow manually specifying the desired 
release version.

# Post Release actions

If the release is a new minor version, then the CLI must be updated with the new version (and any new CRDs) [here](https://github.com/mongodb/atlas-cli-plugin-kubernetes/blob/main/internal/kubernetes/operator/features/crds.go).

If necessary, a CLI plugin release can be created as detailed [here](https://github.com/mongodb/atlas-cli-plugin-kubernetes/blob/main/RELEASING.md).

# Updating the ROSA cluster

For the Openshift upgrade tests we rely on a service account to be present in the OpenShift cluster and its login token to be present in CI.

## Setup Kubectl against the new cluster

1. Go to https://console.redhat.com/openshift
1. Use your RedHat account credentials to log in, see [Prerequisites](#prerequisites) on the RedHat Connect account you need to setup before this.
1. Form the list of Clusters, click of the name of the one to be used now.
1. CLick the `Open Console` in the top right of the page.
1. Use the cluster `htpasswd` credentials you should have been given beforehand to login to the cluster itself.
1. On the landing page, click the account drop down on the top right corner if the page and click on `Copy login command` there.
1. Login again with the `htpasswd`credentials.
1. On the white page click `Display token`.
1. Copy the `oc` command there and run it. You need to have [oc installed](https://docs.openshift.com/container-platform/4.8/cli_reference/openshift_cli/getting-started-cli.html) for this step to work.

After that if you do `kubectl config current-context` it should display you are connected to your new cluster.

## Create the cluster managing service account

Using the kubectl context against the new cluster, create the service account and its token:

```shell
$ kubectl create ns atlas-upgrade-test-tokens
$ kubectl -n atlas-upgrade-test-tokens create serviceaccount atlas-operator-upgrade-test
$ oc create token --duration=87600h -n atlas-upgrade-test-tokens atlas-operator-upgrade-test >token.txt
```

Give this service account enough permissions, currently this is cluster-admin:

```shell
$ oc adm policy add-cluster-role-to-user cluster-admin system:serviceaccount:atlas-upgrade-test-tokens:atlas-operator-upgrade-test
```

Copy & Paste token.txt into the `OPENSHIFT_UPGRADE_TOKEN` secret in Github Actions.

Run `kubectl cluster-info` Eg:

```shell
% kubectl cluster-info
Kubernetes control plane is running at https://***somehostname***.com:6443
...
```

And use the URL there to set `OPENSHIFT_UPGRADE_SERVER_API` so that openshift upgrade tests to run successfully.

## Troubleshooting

### Major version issues when executing the "Create Release Branch" workflow

The release creation will fail if the major version of the release you're creating is incompatible with the `current` version defined in the `version.json` file. This file acts as the single source of truth for the codebase's version, which helps prevent mistakes.

This check allows us to:

1.  **Prevent Mismatches**: It stops the workflow if the version being released (e.g., `v3.0.0`) has a different major version than what the codebase expects (e.g., `current` is `v2.11.0`).
2.  **Avoid Incorrect Patching**: It stops you from accidentally trying to release a patch for an older major version (e.g., `v1.2.3`) when the codebase has already moved on to a new major version.
3.  **Skip Tests**: Certain tests that are expected to fail, like Helm upgrades, they should be skipped.

If the "Create Release Branch" job fails with an error like `Bad major version for X... expected Y...`, you should review the `current` field in `version.json`. Ensure it correctly reflects the codebase's state and is compatible with the version you intend to release.
