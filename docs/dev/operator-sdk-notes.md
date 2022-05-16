## Operator-SDK upgrade process.

Operator SDK has some dependencies:
- [kubebuilder](https://github.com/kubernetes-sigs/kubebuilder) (used to scaffold the project, APIs. It's called when the `operator-sdk init` or `operator-sdk create api` is used)
- controller-runtime - Go dependency that affects the API types used in the code. Upgrade of the library may lead to breaking changes in the Go code.

Sometimes upgrade of the Operator SDK may result in kubebuilder and controller-runtime versions upgrade.
In this case some actions may be necessary to be done: https://sdk.operatorframework.io/docs/upgrading-sdk-version/

The safest Operator SDK upgrade procedure is:
1. Upgrade Operator SDK: `brew upgrade operator-sdk`
1. Scaffold the new project in some temporary directory:
```
mkdir scaffolded
cd scaffolded
operator-sdk init --domain=mongodb.com --repo=github.com/mongodb/mongodb-atlas-kubernetes --license apache2 --owner "MongoDB"
```
1. Compare the scaffolded project with the existing one and manually apply changes.
1. Scaffold the APIs and apply changes manually. Scaffolding one API should be enough to find the difference and apply 
   it to all the other existing APIs
```
operator-sdk create api --group atlas --version v1 --kind AtlasDeployment --resource=true --controller=true
```

### Upgrade of controller-runtime and k8s.io Go dependencies

While dependabot constantly suggests upgrading of these dependencies it's not 
recommended to upgrade them manually. The best upgrade should be after the operator SDK upgrade
(see above) as this will make sure both configs and Go files are aligned.

