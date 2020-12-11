# Onboarding to Atlas Operator

1. Install Go (1.15)
1. Install act (`brew install act` on Mac or check [instructions](https://github.com/nektos/act#installation))
1. Install Kind ([instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation))
1. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
1. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
1. Build and push Operator image to the local registry: `act -j build-push`
1. Deploy Operator into the K8s cluster: `make deploy`
1. Create a Secret with credentials: `kubectl create secret generic my-atlas-key --from-literal="orgId=5fa0197sgfsa90532f1457b7d" --from-literal="publicApiKey=abcdef" --from-literal="privateApiKey=45352345-38d8-4a64-84d3-3425346346" `
1. Create an AtlasProject: `kubectl apply -f config/samples/atlasproject.yaml`

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)
