# Atlas Operator Release Instructions

For the various PRs involved, seek at least someone else to approve. In case of doubts, engage the team member(s) who might be able to clarify and seek their review as well.

## Prerequisites

To get PRs to be auto-committed for RedHat community & Openshift you need to make sure you are listed in the [team members ci.yaml list for community-operators](https://github.com/k8s-operatorhub/community-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml) and [team members ci.yaml list for community-operators-prod](https://github.com/redhat-openshift-ecosystem/community-operators-prod/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml).

This is not required for [Certified Operators](https://github.com/redhat-openshift-ecosystem/certified-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml).

Finally, make sure you have a "RedHat Connect" account and are a [team member with org administrator role in the team list](https://connect.redhat.com/account/team-members).

### Tools

Most tools are automatically installed for you. Most of them are Go binaries and use `go install`. There are a few that might cause issues and you might want to pre-install manually:

- [devbox](https://www.jetify.com/devbox) to be able to enter a sandbox development environment that includes necessary tools for the release process.
- [Docker](https://www.docker.com/) to be able to deal with containers.

## Release preparations (minimum n-1 weeks before the actual release)

At least **one** (1) week before the release the Kubernetes Version testing matrix has to be updated both in this repository and the CLI repository https://github.com/mongodb/mongodb-atlas-cli.

Please refer to the [CI documentation](ci.md#kubernetes-version-matrix) and submit a pull request, example: https://github.com/mongodb/mongodb-atlas-kubernetes/pull/2161 or https://github.com/mongodb/mongodb-atlas-kubernetes/pull/2082.

## Release Notes

- Create a draft of the release notes in a Google Document and share with Product and the Docs team.
  - In confluence, look for the `AKO Release Internal details` page for more details.
- Ensure as well that supporting documents for new features are in review.
- Wait for approval of the release notes and availability of the associated documents.

**DO NOT** start the release process until the release notes are approved and associated documentation is, at least, in review state. Always seek explicit approval by Product and/or Management.

The reason for this preparatory step is to avoid customers getting new or breaking changes before their supporting documentation.

## Create the release branch

Once the release notes and documentation are ready and got explicit approval to start the release:

- Use the [GitHub UI](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/release-branch.yml) to create the new "Create Release Branch" workflow.
- Specify the `version` to be released in the text box and the author MongoDB email or multiple emails in the case of multiple release authors involved in the release.
Do not specify the whole team but only the release authors to respect SSDLC compliance requirements.

The deployment scripts (K8s configs, OLM bundle) will be generated and PR will be created with new changes on behalf
of the `github-actions` bot.

Pass the version with the `X.Y.Z` eg. `1.2.3`, **without** the `v...` prefix.

See [Troubleshooting](#troubleshooting) in case of issues, such as [errors with the major version](#major-version-issues-when-create-release-branch).

Expect this branch to include the Software Security Development Lifecycle Policy Checklist (SSDLC) document at path `docs/releases/v${VERSION}/sdlc-compliance.md`.
Note the SBOM files cannot be generated yet, as they require the image to have been published already.

## Approve the Pull Request named "Release x.y.z"

1. Review the Pull Request.
1. Approve and merge it to `main`.

At this point `main` represents what would become the next release, cut the release by doing:

```shell
$ export VERSION=x.y.z
$ git checkout -b main origin/main
$ git tag v${VERSION}
$ git push origin v${VERSION}
```

A new job "Create Release" will be triggered and the following will be done:
* Atlas Operator image is built and pushed to DockerHub
* Draft Release will be created with all commits since the previous release
* A subsequent job will be triggered to create a SBOMs update PR

### SSDLC SBOMs PR

A new PR should have been created titled `Add SBOMs for version ...`. Please review all is as expected and merge. It should contain just a couple of new files at directory `docs/releases/v${VERSION}/`:
- `linux-amd64.sbom.json`
- `linux-arm64.sbom.json`

## Manual SSDLC steps

### Process Overview

The SSDLC process requirements are as follows:

1. Sign our images with a MongoDB owned signature.
1. Produce SBOM (Software Bill Of Materials) for each platform we support (`linux-amd64` and `linux-arm64`).
1. Upload the plain SBOMs to a MongoDB internal Kondukto service instance.
1. Produce the augmented SBOMS, including vulnerability metadata, from using Silkbomb 2.0.
1. Store both sets of SBOM files for internal reference.

The first two steps are semi-automated as documented here. The rest is fully manual.

Right now we are only using **one Kondukto branch per platform**:
- `main-linux-amd64`
- `main-linux-arm64`

This means only the latest version is tracked by Kondukto. Note each upload will replace the SBOM document tracked on each asset group.

For more details about credentials required, to to `MongoDB Confluence` and look for page:
`Kubernetes Atlas Operator SSDLC Compliance Manual`

What follows is a quick reference of the make rules involved, assuming the credential setup is already completed and the process is already familiar.

### Upload SBOMs to Kondukto and Augment SBOMs with Kondukto Scan results

Make sure that you have the credentials configured to handle SBOM artifacts.
Read through the wiki page "Kubernetes Atlas Operator SSDLC Compliance Manual" on how to retrieve them.

Update the local `main` branch to point to the commit which includes the merged SSDLC files from the previous step:

```shell
$ git checkout main
$ git pull
```

```shell
$ make augment-sbom SBOM_JSON_FILE="docs/releases/v${VERSION}/linux_amd64.sbom.json"
$ make augment-sbom SBOM_JSON_FILE="docs/releases/v${VERSION}/linux_arm64.sbom.json"
```

### Register SBOMs internally

To be able to store SBOMs in S3, you need special credentials.
Please advise the Wiki page "Kubernetes Atlas Operator SSDLC Compliance Manual".

```shell
$ make store-augmented-sboms VERSION=${VERSION} TARGET_ARCH=amd64
$ make store-augmented-sboms VERSION=${VERSION} TARGET_ARCH=arm64
```

## Edit the Release Notes and publish the release

Follow the format described in the [release-notes-template.md](../release-notes/release-notes-template.md) file.
Paste the release notes content approved before the release was started.
Once the image is out, publish the release notes draft as soon as possible.

## Synchronize configuration changes with the Helm Charts

Go to the [helm-chart repo](https://github.com/mongodb/helm-charts) and locate the [Pull Request](https://github.com/mongodb/helm-charts/pulls)
that is being automatically generated by the [GitHub "Create PR with Atlas Operator Release" action](https://github.com/mongodb/helm-charts/actions/workflows/post-atlas-operator-release.yaml).
It is named "Release Atlas Operator x.y.z.".

The will update two Helm charts:
* [atlas-operator-crds](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator-crds)
* [atlas-operator](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator)
    
Merge the PR - the chart will get released automatically.

## Create Pull Requests to publish OLM bundles

All bundles/package manifests for Operators for operatorhub.io reside in the following repositories:
* https://github.com/k8s-operatorhub/community-operators - Kubernetes Operators that appear on [OperatorHub.io](https://operatorhub.io/)
* https://github.com/redhat-openshift-ecosystem/community-operators-prod - Kubernetes Operators that appear on [OpenShift](https://openshift.com/) and [OKD](https://www.okd.io/)
* https://github.com/redhat-openshift-ecosystem/certified-operators - Red Hat certified Kubernetes Operators

### Fork/Update the community operators repositories

**Note**: this has to be done once only. 

First ensure your SSH keys in [https://github.com/settings/keys] are authorized for `mongodb-forks` MongoDB SSO.

Execute the following steps:

1. Clone each of the above forked OLM repositories from https://github.com/mongodb-forks
2. Add `upstream` remotes
3. Export each cloned repository directory in environment variables

#### community-operators repository
```
git clone git@github.com:mongodb-forks/community-operators.git
git remote add upstream https://github.com/k8s-operatorhub/community-operators.git
export RH_COMMUNITY_OPERATORHUB_REPO_PATH=$PWD/community-operators
```
#### community-operators-prod repository
```
git clone git@github.com:mongodb-forks/community-operators-prod.git
git remote add upstream https://github.com/redhat-openshift-ecosystem/community-operators-prod.git
export RH_COMMUNITY_OPENSHIFT_REPO_PATH=$PWD/community-operators-prod
```
#### certified-operators repository
```
git clone git@github.com:mongodb-forks/certified-operators.git
git remote add upstream https://github.com/redhat-openshift-ecosystem/certified-operators
export RH_CERTIFIED_OPENSHIFT_REPO_PATH=$PWD/certified-operators
```

### Create a Pull Request for the `community-operators` repository

1. Ensure the `RH_COMMUNITY_OPERATORHUB_REPO_PATH` environment variable is set.
2. Invoke the following script with `<version>` set to `1.0.0` (don't use a `v` prefix):
```
./scripts/release-redhat.sh <version>
```

You can see an [example fixed PR here on Community Operators for version 1.9.1](https://github.com/k8s-operatorhub/community-operators/pull/3457).

Create the PR to the main repository and wait until CI jobs get green. 
After the PR is approved and merged - it will soon get available on https://operatorhub.io

### Create a Pull Request for the `community-operators-prod` repository

1. Ensure the `RH_COMMUNITY_OPENSHIFT_REPO_PATH` environment variable is set.
2. Invoke the following script with `<version>` set to `1.0.0` (don't use a `v` prefix):
```
./scripts/release-redhat-openshift.sh <version>
```

Submit the PR to the upstream repository and wait until CI jobs get green.

**Note**: It is required that the PR consists of only one commit - you may need to do
`git rebase -i HEAD~2; git push origin +mongodb-atlas-operator-community-<version>` if you need to squash multiple commits into one and perform force push)

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)

### Create a Pull Request for the `certified-operators` repository

This is necessary for the Operator to appear on "operators" tab in Openshift clusters in the "certified" section.
Ensure the `RH_CERTIFIED_OPENSHIFT_REPO_PATH` environment variable is set.

Invoke the following script and ensure to have the `VERSION` variable set from above:
```
./scripts/release-redhat-certified.sh
```

Then go the GitHub and create a PR
from the `mongodb-fork` repository to https://github.com/redhat-openshift-ecosystem/certified-operators (`origin`).

Note: For some reason, the certified OpenShift metadata does not use the multi arch image reference at all, and only understand direct architecture image references.

You can see an [example fixed PR here for certified version 1.9.1](https://github.com/redhat-openshift-ecosystem/certified-operators/pull/3020).

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)

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
1. Use your RedHat account credentials to log in, see Pre-requisites on the RedHat Connect account you need to setup before this.
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

The release creation will fail if the file `major-version` contents does not match the major version to be released. This file explicitly means the upcoming release is for a particular major version, with potential breaking changes. This allows us to:

1. Notice if we forgot to update the `major-version` file before releasing the next major version.
2. Notice if we tried to re-release an older major version when the code is already prepared for the next major version.
3. Skip some tests, like `helm update`, when crossing from one major version to the next, as such test is not expected to work across incompatible major version upgrades.

If the create release branch job fails due an error such as `Bad major version for X... expected Y..`, review whether or not the `major-version` file was updated as expected. Check as well you are not trying to release a patch for the older major version from the new major version codebase.
