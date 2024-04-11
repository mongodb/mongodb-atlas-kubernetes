# Atlas Operator Release Instructions

## Create the release branch

Use the GitHub UI to create the new "Release Branch" workflow. Specify the version to be released in the text box.
The deployment scripts (K8s configs, OLM bundle) will be generated and PR will be created with new changes on behalf
of the `github-actions` bot.

NOTE: The X- and Y- stream releases should only be launched using the workflow from the MAIN branch. Z-stream (patch)
releases can be launched from a separate branch

The release creation will fail if the file `major-version` contents does not match the major version to be released. This file explicitly means the upcoming release is for a particular major version, with potential breaking changes. This allows us to:

1. Notice if we forgot to update the `major-version` file before releasing the next major version.
2. Notice if we tried to re-release an older major version when the code is already prepared for the next major version.
3. Skip some tests, like `helm update`, when crossing from one major version to the next, as such test is not expected to work across incompatible major version upgrades.

If the create release branch job fails due an error such as `Bad major version for X... expected Y..`, review whether or not the `major-version` file was updated as expected. Check as well you are not trying to release a patch for the older major version from the new major version codebase.

## Approve the Pull Request named "Release x.y.z"

Review the Pull Request. Approve and merge it to main.
The new job "Create Release" will be triggered and the following will be done:
* Atlas Operator image built and pushed to dockerhub
* Draft Release will be created with all commits since the previous release

Once the Pull Request is approved, a tag is created out of the branch, which can then be discarded. A branch `release/X.Y.Z` will imply a tag `vX.Y.Z`, and `pre-release/X.Y.Z-...` will imply `vX.Y.Z-...`. The `tag.yml` workflow is the one responsible for creating such tag and triggering the release process workflow (`release-post-merge.yml`).

## Edit the Release Notes and publish the release

Follow the format described in the release-notes-template.md file. Publish the release.

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

**Note**: this has to be done once only:
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

# Post install hook release

If changes have been made to the post install hook (mongodb-atlas-kubernetes/cmd/post-install/main.go).
You must also release this image. Run the "Release Post Install Hook" workflow manually specifying the desired 
release version. 

# Post Release actions

Following the instructions above the release is completed, but the AKO repo is left using a `helm-charts/` submodule which is not pointing to the latest released Helm Chart. Make a PR to update that submodule.
