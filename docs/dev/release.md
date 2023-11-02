# Atlas Operator Release Instructions

## Create the release branch

Use the GitHub UI to create the new "Release Branch" workflow. Specify the version to be released in the text box.
The deployment scripts (K8s configs, OLM bundle) will be generated and PR will be created with new changes on behalf
of the `github-actions` bot.

NOTE: The X- and Y- stream releases should only be launched using the workflow from the MAIN branch. Z-stream (patch)
releases can be launched from a separate branch

## Approve the Pull Request named "Release x.y.z"

Review the Pull Request. Approve and merge it to main.
The new job "Create Release" will be triggered and the following will be done:
* Atlas Operator image built and pushed to dockerhub
* Draft Release will be created with all commits since the previous release

## Edit the Release Notes and publish the release

Follow the format described in the release-notes-template.md file. Publish the release.

## Synchronize configuration changes with the Helm Charts

Go to the [helm-chart repo](https://github.com/mongodb/helm-charts) and use GitHub Action [Create PR with Atlas Operator Release](https://github.com/mongodb/helm-charts/actions/workflows/post-atlas-operator-release.yaml). Run the workflow on branch `main` setting the version that is being release, eg. `1.9.0`.

The will update two Helm charts:
* [atlas-operator-crds](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator-crds)
* [atlas-operator](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator)
    
Merge the PR - the chart will get released automatically.

## Create the Pull Request to publish the bundle to operatorhub.io

All bundles/package manifests for Operators for operatorhub.io reside in:
* `https://github.com/k8s-operatorhub/community-operators` - for public Operators from operatorhub.io
* `https://github.com/redhat-openshift-ecosystem/community-operators-prod` - for Operators from "internal" operatorhub that are synchronized with Openshift clusters

### Fork/Update the community operators repositories

**(First time only) Fork 2 separate repositories**

#### 1. OperatorHub

Clone, if not done before, the MongoDB fork of [the community operators repo](https://github.com/k8s-operatorhub/community-operators):

```bash
git clone git@github.com:mongodb-forks/community-operators.git
```

Add the upstream repository as a remote one:

```bash
git remote add upstream https://github.com/k8s-operatorhub/community-operators.git
```

Assign the repo path to `RH_COMMUNITY_OPERATORHUB_REPO_PATH` env variable.

#### 2. Openshift

Clone, if not done before, the MongoDB fork of [the OpenShift Community Operators repo](https://github.com/redhat-openshift-ecosystem/community-operators-prod):

```bash
git clone git@github.com:mongodb-forks/community-operators-prod.git
```

Add the upstream repository as a remote one:

```bash
git remote add upstream https://github.com/redhat-openshift-ecosystem/community-operators-prod.git
```

Assign the repo path to `RH_COMMUNITY_OPENSHIFT_REPO_PATH` env variable.

#### 3. OpenShift Certified

Clone, if not done before, the MongoDB fork of [the Red Hat certified operators production catalog repo](https://github.com/redhat-openshift-ecosystem/certified-operators):

```bash
git clone git@github.com:mongodb-forks/community-operators-prod.git
```

Add the upstream repository as a remote one:

```bash
git remote add upstream https://github.com/redhat-openshift-ecosystem/certified-operators
```

Assign the repo path to `RH_CERTIFIED_OPENSHIFT_REPO_PATH` env variable.

### Create a Pull Request for `operatorhub` with a new bundle

This is necessary for the Operator to appear on [operatorhub.io] site.
This step should be done after the previous PR is approved and merged.

Ensure you have the `RH_COMMUNITY_OPERATORHUB_REPO_PATH` environment variable exported in `~/.bashrc` or `~/.zshrc`
pointing to the directory where `operatorhub-operator` repository was cloned in the previous step.

For this PR the sources are copied from the `community-operators` folder instead of the one where the `mongodb-atlas-kubernetes` resides.

Invoke with <version> like `1.0.0` (never use the `v` prefix here, just the plain SEMVER version `x.y.z`):

```
./scripts/release-redhat.sh <version>
```

Before posting the PR there is a manual change you need to make:

* Ensure to add the `quay.io/` prefix in all Operator image references.

You can see an [example fixed PR here on Community Operators for version 1.9.1](https://github.com/k8s-operatorhub/community-operators/pull/3457).

Create the PR to the main repository and wait until CI jobs get green. 
After the PR is approved and merged - it will soon get available on https://operatorhub.io

### Create a Pull Request for `openshift` with a new bundle

This is necessary for the Operator to appear on "operators" tab in Openshift clusters

Ensure you have the `RH_COMMUNITY_OPERATORHUB_REPO_PATH` environment variable exported in `~/.bashrc` or `~/.zshrc`
pointing to the directory where `community-operators-prod` repository was cloned in the previous step.

*(This is temporary, to be fixed)
Change the `mongodb-atlas-kubernetes.clusterserviceversion.yaml` file and change the `replaces:` setting the previous version

Invoke the following script with <version> like `1.0.0` (no `v` prefix):
```
./scripts/release-redhat-openshift.sh <version>
```

Before posting the PR there is a manual change you need to make:

* Ensure to add the `quay.io/` prefix in all Operator image references.

You can see an [example fixed PR here on OpenShift for version 1.9.1](https://github.com/redhat-openshift-ecosystem/community-operators-prod/pull/3521).

Create the PR to the main repository and wait until CI jobs get green.

(note, that it's required that the PR consists of only one commit - you may need to do
`git rebase -i HEAD~2; git push origin +mongodb-atlas-operator-community-<version>` if you need to squash multiple commits into one and perform force push)

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)

### Create a Pull Request for `openshift-certified-operators` with a new bundle

This is necessary for the Operator to appear on "operators" tab in Openshift clusters in the "certified" section.

**Prerequisites**:
 - Ensure you have the `RH_CERTIFIED_OPENSHIFT_REPO_PATH` environment variable exported in `~/.bashrc` or `~/.zshrc`
pointing to the directory where `certified-operators` repository: https://github.com/redhat-openshift-ecosystem/certified-operators.
 - Download (and build locally, if you're running MacOS) https://github.com/redhat-openshift-ecosystem/openshift-preflight and put the binary to your `$PATH`
 - Use the image reference including the hash (`quay.io/mongodb/mongodb-atlas-kubernetes-operator:...@sha256:...`) from the [release process step "Push Atlas Operator to Quay.io"](https://github.com/mongodb/mongodb-atlas-kubernetes/actions/workflows/release-post-merge.yml) as `IMG_SHA`
 - Use the version of the release as `VERSION`, remember the SEMVER x.y.z version without NPO `v`prefix.

Invoke the following script:
```
IMG_SHA=<image hash pushed to scan.connect.redhat.com with sha had rather than tag> \
VERSION=<image-version> \
./scripts/release-redhat-certified.sh
```

If script successfully finishes, you should be able to see new tag (e.g. 1.2.0) here https://connect.redhat.com/projects/63568bb95612f26f8db42d7a/images

Then go the GitHub and create a PR
from the `mongodb-fork` repository to https://github.com/redhat-openshift-ecosystem/certified-operators (`origin`).

Before posting the PR there are manual changes you need to make:

1. Ensure to add the `quay.io/` prefix in all Operator image references.
1. Add missing a `com.redhat.openshift.versions: "v4.8"` line at the end of `metadata/annotations.yaml`.
1. Ensure all image references, including `containerImage`, do NOT use the version *tag*. They **should only use the SHA of the AMD image**, NEVER the multi arch SHA.
1. Add missing `spec.relatedImages` section in certified pinning all the images per architecture.

For some reason, the certified OpenShift metadata does not use the multi arch image reference at all, and only understand direct architecture image references.

You can see an [example fixed PR here for certified version 1.9.1](https://github.com/redhat-openshift-ecosystem/certified-operators/pull/3020).

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)

# Post install hook release

If changes have been made to the post install hook (mongodb-atlas-kubernetes/cmd/post-install/main.go).
You must also release this image. Run the "Release Post Install Hook" workflow manually specifying the desired 
release version. 

# Post Release actions

Following the instructions above the release is completed, but the AKO repo is left using a `helm-charts/` submodule which is not pointing to the latest released Helm Chart. Make a PR to update that submodule.
