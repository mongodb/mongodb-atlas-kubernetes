# Atlas Operator Release Instructions

For the various PRs involved, seek at least someone else to approve. In case of doubts, engage the team member(s) who might be able to clarify and seek their review as well.

## Create the release branch

Use the GitHub UI to create the new "Create Release Branch" workflow. Specify the version to be released in the text box.
The deployment scripts (K8s configs, OLM bundle) will be generated and PR will be created with new changes on behalf
of the `github-actions` bot.

Pass the version with the `X.Y.Z` eg. `1.2.3`, **without** the `v...` prefix.

See [Troubleshooting](#troubleshooting) in case of issues, such as [errors with the major version](#major-version-issues-when-create-release-branch).

## Approve the Pull Request named "Release x.y.z"

Review the Pull Request. Approve and merge it to `main`.
The new job "Create Release" will be triggered and the following will be done:
* Atlas Operator image is built and pushed to DockerHub
* Draft Release will be created with all commits since the previous release

Once the Pull Request is approved, a tag is created out of the branch, so such branch could be discarded. A branch `release/X.Y.Z` turns into tag `vX.Y.Z`, and `pre-release/X.Y.Z-...` into `vX.Y.Z-...`. The `tag.yml` workflow is the one responsible for creating a tag from the branch closed, and triggering the release process workflow (`release-post-merge.yml`).

## Edit the Release Notes and publish the release

Follow the format described in the [release-notes-template.md](../release-notes/release-notes-template.md) file. Publish the release.

## Synchronize configuration changes with the Helm Charts

Go to the [helm-chart repo](https://github.com/mongodb/helm-charts) and use GitHub Action [Create PR with Atlas Operator Release](https://github.com/mongodb/helm-charts/actions/workflows/post-atlas-operator-release.yaml). Run the workflow on branch `main` setting the version that is being release, eg. `1.9.0`.

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

Then, from your usual directory holding your git repository copies, do the following:

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

### Before creating PRs against RedHat repositories

To get PRs to be auto-committed for RedHat community & Openshift you need to make sure you are listed in the [team members ci.yaml list (community-operators)](https://github.com/k8s-operatorhub/community-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml). Link for [community-operators-prod](https://github.com/redhat-openshift-ecosystem/community-operators-prod/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml). This is not a requirement for [Certified Operators](https://github.com/redhat-openshift-ecosystem/certified-operators/blob/main/operators/mongodb-atlas-kubernetes/ci.yaml).

All this is in addition to holding a RedHat Connect account adn being a [team member with org administrator role in the team list](https://connect.redhat.com/account/team-members), which is part of the *team onboarding*.

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

1. Ensure the `RH_CERTIFIED_OPENSHIFT_REPO_PATH` environment variable is set.
2. Set the image SHA environment variables of the **certified** images. 
To get the SHAs:

1. Go to https://connect.redhat.com/manage/components
2. search for "mongodb-atlas-kubernetes-operator"
3. select "[Quay] mongodb-atlas-kubernetes-operator"
4. filter for the given tag, i.e. "2.2.0"

The direct link, at the time of writing is https://connect.redhat.com/component/view/63568bb95612f26f8db42d7a/images.

Copy the **certified** image SHAs of the **amd64** and the **arm64** image:

![img.png](certified-image-sha.png)

```
export IMG_SHA_AMD64=sha256:c997f8ab49ed5680c258ee4a3e6a9e5bbd8d8d0eef26574345d4c78a4f728186
export IMG_SHA_ARM64=sha256:aa3ed7b73f8409dda9ac32375dfddb25ee52d7ea172e08a54ecd144d52fe44da
```

 - Use the version of the release as `VERSION`, remember the SEMVER x.y.z version without the `v`prefix.

```
export VERSION=<image-version>
```

Invoke the following script:
```
./scripts/release-redhat-certified.sh
```

Then go the GitHub and create a PR
from the `mongodb-fork` repository to https://github.com/redhat-openshift-ecosystem/certified-operators (`origin`).

Note: For some reason, the certified OpenShift metadata does not use the multi arch image reference at all, and only understand direct architecture image references.

You can see an [example fixed PR here for certified version 1.9.1](https://github.com/redhat-openshift-ecosystem/certified-operators/pull/3020).

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)

# SSDLC checklist publishing

For the time being, preparing the SSDLC checklist for each release is a manual process. Use this [past PR as a starting point](https://github.com/mongodb/mongodb-atlas-kubernetes/pull/1524).

Copy the closest [sdlc-compliance.md](../releases/v2.2.1/sdlc-compliance.md) file and:
- Update the **version** references to the one being released.
- Update dates and release creators to match the reality of the current release.
- Usually no more changes will be needed. Only if you actually skipped some CI check, like allowed the release to happen even when some report flagged a dependency that had no fix yet, you might need to mention such case.

There is no need to make any changes about the image signature part other than updating the instructions to the current version.

For SBOMs, you will have to *generate* the files and place them in the same directory as the compliance doc `docs/releases/vX.Y.Z`, with the expected names; `linux-amd64.sbom.json` & `linux-arm64.sbom.json`. To generate each f the files you need to run:

```shell
docker sbom --platform "linux/${arch}" -o "docs/releases/v${version}/linux-${arch}.sbom.json" --format "cyclonedx-json" "$image"
```

Where:
- `${arch}` is the architecture to generate, either `amd64` or `arm64`.
- `${version}` is the current version released in `X.Y.Z` format, without the **v** prefix.
- `${image}` is the image reference released, usually something like `mongodb/mongodb-atlas-kubernetes-operator:${version}`.

Once all such information is in place, create a PR with it and merge it as close to the release as possible.

# Post install hook release

If changes have been made to the post install hook (mongodb-atlas-kubernetes/cmd/post-install/main.go).
You must also release this image. Run the "Release Post Install Hook" workflow manually specifying the desired 
release version. 

# Post Release actions

Following the instructions above the release is completed, but the AKO repo is left using a `helm-charts/` submodule which is not pointing to the latest released Helm Chart. Make a PR to update that submodule:
```
cd helm-charts
git pull
```
Checkout to the branch for the release:
```
git checkout mongodb-atlas-kubernetes-<VERSION>
```
Return up to the AKO repo, and create a new branch & PR with the newer submodule:
```
cd ..
git checkout -b update-helm-submodule-<VERSION>
git add helm-charts
git push
```

If the release is a new minor version, then the CLI must be updated with the new version (and any new CRDs) [here](https://github.com/mongodb/mongodb-atlas-cli/blob/master/internal/kubernetes/operator/features/crds.go).

## Troubleshooting

### Major version issues when Create Release Branch

The release creation will fail if the file `major-version` contents does not match the major version to be released. This file explicitly means the upcoming release is for a particular major version, with potential breaking changes. This allows us to:

1. Notice if we forgot to update the `major-version` file before releasing the next major version.
2. Notice if we tried to re-release an older major version when the code is already prepared for the next major version.
3. Skip some tests, like `helm update`, when crossing from one major version to the next, as such test is not expected to work across incompatible major version upgrades.

If the create release branch job fails due an error such as `Bad major version for X... expected Y..`, review whether or not the `major-version` file was updated as expected. Check as well you are not trying to release a patch for the older major version from the new major version codebase.