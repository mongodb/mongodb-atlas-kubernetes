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

## Create the Pull Request to publish the bundle to operatorhub.io

All bundles/package manifests for Operators for operatorhub.io reside in `https://github.com/operator-framework/community-operators/tree/master/upstream-community-operators`

### (First time) Fork the repository
Fork the following repo into your own:

    https://github.com/operator-framework/community-operators/tree/master/upstream-community-operators

Make sure you clone the *fork* and not *upstream*.

Add the upstream repository as a remote one

```bash
git remote add upstream git@github.com:operator-framework/community-operators.git
```

### Update the forked repository
Pull changes from the upstream:

```bash
git fetch upstream
git checkout master
git merge upstream/master
```