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

### Fork/Update the community operators repository
**(First time) Fork the repository**
Fork the following repo into your own:

    https://github.com/operator-framework/community-operators/tree/master/upstream-community-operators

Make sure you clone the *fork* and not *upstream*.

Add the upstream repository as a remote one

```bash
git remote add upstream git@github.com:operator-framework/community-operators.git
```

**(Not the first time) Update the forked repository**
Pull changes from the upstream:

```bash
git fetch upstream
git checkout master
git merge upstream/master
```

### Create a Pull Request with a new bundle

```
version=0.5.0
mkdir <community-operators-repo>/upstream-community-operators/mongodb-atlas-kubernetes/${version}
cp bundle.Dockerfile bundle/manifests bundle/metadata <community-operators-repo>/upstream-community-operators/mongodb-atlas-kubernetes/${version}
cd <community-operators-repo>
git checkout -b "mongodb-atlas-operator-${version}"
git commit -m "MongoDB Atlas Operator ${version}" --signoff * 
git push origin mongodb-atlas-operator-${version}
```

(note, that it's required that the PR consists of only one commit - you may need to do 
`git rebase -i HEAD~2; git push origin +master` if you need to squash multiple commits into one and perform force push)

Create the PR to the main repository and wait until CI jobs get green. 
After the PR is approved and merged - it will soon get available on https://operatorhub.io

Example PR: https://github.com/operator-framework/community-operators/pull/3281
