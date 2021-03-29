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
* `https://github.com/operator-framework/community-operators/tree/master/upstream-community-operators` - for public Operators from operatorhub.io
* `https://github.com/operator-framework/community-operators/tree/master/community-operators` - for Operators from "internal" operatorhub that are synchronized with Openshift clusters

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

### Create a Pull Request to `community-operators` with a new bundle

This is necessary for the Operator to appear on "operators" tab in Openshift clusters

```
version=0.5.0
mkdir <community-operators-repo>/community-operators/mongodb-atlas-kubernetes/${version}
cp bundle.Dockerfile bundle/manifests bundle/metadata <community-operators-repo>/community-operators/mongodb-atlas-kubernetes/${version}
cd <community-operators-repo>/community-operators/mongodb-atlas-kubernetes/${version}

# replace the move instructions in the docker file
sed -i .bak 's/COPY bundle\/manifests/COPY manifests/' bundle.Dockerfile
sed -i .bak 's/COPY bundle\/metadata/COPY metadata/' bundle.Dockerfile
sed -i .bak '/COPY bundle\/tests\/scorecard \/tests\/scorecard\//d' bundle.Dockerfile
rm bundle.Dockerfile.bak

# commit
git checkout -b "mongodb-atlas-operator-community-${version}"
git commit -m "MongoDB Atlas Operator ${version}" --signoff * 
git push origin mongodb-atlas-operator-community-${version}
```

*(This should be automated in the future by the correct layout of kustomize configs)*
Remove the following lines for `Deployment` from the CSV file and leave the `securityContext` empty:

```
          runAsNonRoot: true
          runAsUser: 2000
```

Create the PR to the main repository and wait until CI jobs get green.

(note, that it's required that the PR consists of only one commit - you may need to do
`git rebase -i HEAD~2; git push origin +master` if you need to squash multiple commits into one and perform force push)

See https://github.com/operator-framework/community-operators/pull/3343 as an example

### Create a Pull Request to `upstream-community-operators` with a new bundle

This is necessary for the Operator to appear on [operatorhub.io] site. 
For this PR the sources are copied from the `community-operators` folder instead of the one where the `mongodb-atlas-kubernetes` resides.

```
# cd <community-operators-repo>
cp -r community-operators/mongodb-atlas-kubernetes/${version} upstream-community-operators/mongodb-atlas-kubernetes
git add upstream-community-operators/mongodb-atlas-kubernetes/${version}
git checkout -b "mongodb-atlas-operator-community-${version}"
git commit -m "[community] MongoDB Atlas Operator ${version}" --signoff upstream-community-operators/mongodb-atlas-kubernetes/${version}
git push origin mongodb-atlas-operator-community-${version}
git push
```

Create the PR to the main repository and wait until CI jobs get green. 
After the PR is approved and merged - it will soon get available on https://operatorhub.io
Example PR: https://github.com/operator-framework/community-operators/pull/3281