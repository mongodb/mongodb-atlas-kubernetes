# Atlas Operator Release Instructions

## Create the release branch
Use the GitHub UI to create the new "Release Branch" workflow. Specify the version to be released in the text box.
The deployment scripts (K8s configs, OLM bundle) will be generated and PR will be created with new changes on behalf
of the `github-actions` bot.

## Approve the Pull Request named "Release x.y.z"
Review the Pull Request. Approve and merge it to main.
The new job "Create Release" will be triggered and the following will be done:
* Atlas Operator image built and pushed to dockerhub
* Draft Release will be created with all commits since the previous release

## Edit the Release Notes and publish the release
Follow the format of the previous release notes (two main sections: "Features" and "Bug Fixes"). Publish the release.

## Synchronize configuration changes with the Helm Charts

Create a PR to https://github.com/mongodb/helm-charts to update two Helm charts:
* [atlas-operator-crds](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator-crds)
* [atlas-operator](https://github.com/mongodb/helm-charts/tree/main/charts/atlas-operator)
  * `Chart.yaml` - update the `AppVersion` to the new Operator version and increment the minor digit for `version`
  * any changes to `templates` configuration
    
Merge the PR - the chart will get released automatically.

## Create the Pull Request to publish the bundle to operatorhub.io

All bundles/package manifests for Operators for operatorhub.io reside in:
* `https://github.com/k8s-operatorhub/community-operators` - for public Operators from operatorhub.io
* `https://github.com/redhat-openshift-ecosystem/community-operators-prod` - for Operators from "internal" operatorhub that are synchronized with Openshift clusters

### Fork/Update the community operators repositories

**(First time only) Fork 2 separate repositories**

#### 1. OperatorHub

Fork the following repo into your own:
  https://github.com/k8s-operatorhub/community-operators

Clone the *fork* and not *upstream* to your maching.

Add the upstream repository as a remote one
```bash
git remote add upstream https://github.com/k8s-operatorhub/community-operators.git
```

Assign the path to the repo to `RH_COMMUNITY_OPERATORHUB_REPO_PATH` env variable.

#### 2. Openshift

Fork the following repo into your own:
  https://github.com/redhat-openshift-ecosystem/community-operators-prod

Clone the *fork* and not *upstream* to your maching.

Add the upstream repository as a remote one
```bash
git remote add upstream https://github.com/redhat-openshift-ecosystem/community-operators-prod.git
```

Assign the path to the repo to `RH_COMMUNITY_OPENSHIFT_REPO_PATH` env variable.

### Create a Pull Request for `operatorhub` with a new bundle

This is necessary for the Operator to appear on [operatorhub.io] site.
This step should be done after the previous PR is approved and merged.

Ensure you have the `RH_COMMUNITY_OPERATORHUB_REPO_PATH` environment variable exported in `~/.bashrc` or `~/.zshrc`
pointing to the directory where `operatorhub-operator` repository was cloned in the previous step.

For this PR the sources are copied from the `community-operators` folder instead of the one where the `mongodb-atlas-kubernetes` resides.

Invoke with <version> like `1.0.0`:
```
./scripts/release-redhat.sh <version>
```

Create the PR to the main repository and wait until CI jobs get green. 
After the PR is approved and merged - it will soon get available on https://operatorhub.io
Example PR: https://github.com/k8s-operatorhub/community-operators/pull/69

### Create a Pull Request for `openshift` with a new bundle

This is necessary for the Operator to appear on "operators" tab in Openshift clusters

Ensure you have the `RH_COMMUNITY_OPERATORHUB_REPO_PATH` environment variable exported in `~/.bashrc` or `~/.zshrc`
pointing to the directory where `community-operators-prod` repository was cloned in the previous step.

*(This is temporary, to be fixed)
Change the `mongodb-atlas-kubernetes.clusterserviceversion.yaml` file and change the `replaces:` setting the previous version

Invoke the following script with <version> like `1.0.0`:
```
./scripts/release-redhat-openshift.sh <version>
```

Create the PR to the main repository and wait until CI jobs get green.

(note, that it's required that the PR consists of only one commit - you may need to do
`git rebase -i HEAD~2; git push origin +mongodb-atlas-operator-community-<version>` if you need to squash multiple commits into one and perform force push)

After the PR is approved it will soon appear in the [Atlas Operator openshift cluster](https://console-openshift-console.apps.atlas.operator.mongokubernetes.com)
See https://github.com/redhat-openshift-ecosystem/community-operators-prod/pull/98 as an example

# Post install hook release

If changes have been made to the post install hook (mongodb-atlas-kubernetes/cmd/post-install/main.go).
You must also release this image. Run the "Release Post Install Hook" workflow manually specifying the desired 
release version. 
