# Onboarding to Atlas Operator

1. Install Go (1.15)
1. Install act (`brew install act` on Mac or check [instructions](https://github.com/nektos/act#installation))
1. Install `yq`(`brew install yq` on Mac)
1. Install Kind ([instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))
1. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
1. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
1. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
1. Build and push Operator image to the local registry: `act -j build-push`
1. Deploy Operator into the K8s cluster: `make deploy`
1. Create an AtlasProject: `kubectl apply -f config/samples/atlasproject.yaml` (note, that the secret `my-atlas-key` is
 created during running `make deploy`)

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)
