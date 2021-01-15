# Onboarding to Atlas Operator

1. Install Go (1.15)
1. Install act (`brew install act` on Mac or check [instructions](https://github.com/nektos/act#installation))
1. Install `yq`(`brew install yq` on Mac)
1. Install Kind ([instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))
1. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
1. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
1. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
1. Build and deploy the Operator into the K8s cluster: `make deploy`
1. Create an AtlasProject: `kubectl apply -f config/samples/atlasproject.yaml` (note, that the secret `my-atlas-key` is
 created during running `make deploy`)
   
Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)

#How-To
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

