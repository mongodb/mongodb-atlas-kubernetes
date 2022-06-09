# Onboarding to Atlas Operator

1. Install Go (1.16)
2. Install dependencies
```
# on Mac
brew install kubebuilder # or https://book.kubebuilder.io/quick-start.html
brew install kustomize # or https://kubectl.docs.kubernetes.io/installation/kustomize/
brew install act # or https://github.com/nektos/act#installation
brew install yq # or https://pypi.org/project/yq/
brew install shellcheck # or https://github.com/koalaman/shellcheck#installing
brew install pre-commit # or https://pre-commit.com/index.html#install
pre-commit install # from the root of the project
```
3. Install Kind ([instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))
4. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
5. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
6. Copy the default Github Actions environment for local run: `cp .env.sample .env`
7. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
8. Build and deploy the Operator into the K8s cluster: `make deploy`
9. Create an AtlasProject: `kubectl apply -f config/samples/atlas_v1_atlasproject.yaml` (note, that Atlas connection secrets are
 created during running `make deploy`)
10. Create an AtlasDeployment: `kubectl apply -f config/samples/atlas_v1_atlasdeployment.yaml`
11. Create an AtlasDatabaseUser: `kubectl apply -f config/samples/atlas_v1_atlasdatabaseuser.yaml`

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)

# How-To
## Run integration tests
### make
When running the tests using `make` the Atlas credentials from `.actrc` will be used automatically to export environment
variables
```bash
make int-test
```

### IDE
When running integration tests from an IDE the following environment variables need to be provided to `go test` / `ginkgo`:
`KUBEBUILDER_ASSETS=<path-to-project>/mongodb-atlas-kubernetes/testbin/bin`
`ATLAS_ORG_ID=<..>`
`ATLAS_PUBLIC_KEY=<..>`
`ATLAS_PRIVATE_KEY=<..>`
`GINKGO_EDITOR_INTEGRATION=true`
