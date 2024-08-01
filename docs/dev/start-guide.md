# Onboarding to Atlas Operator
1. Install Nix [[Instructions](https://nixos.org/download/)]
2. Run the Nix shell (by entering ```shell.nix``` in the bash to initialize an environment with pre-installed dependencies). For more details, refer to the shell.nix file.
3. Install optional dependencies
    ```
    # on Mac
    brew install coreutils # or https://www.gnu.org/software/coreutils/
    brew install pre-commit # or https://pre-commit.com/index.html#install
    pre-commit install # from the root of the project
    ```
4. Install Kind [[Instructions](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)]
4. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
5. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
6. Copy the default Github Actions environment for local run: `cp dotenv.sample .env`
7. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
8. Build and deploy the Operator into the K8s cluster: `make deploy`
9. Create an AtlasProject: `kubectl apply -f config/samples/atlas_v1_atlasproject.yaml` (note, that Atlas connection secrets are created during running `make deploy`)
10. Create an AtlasDeployment: `kubectl apply -f config/samples/atlas_v1_atlasdeployment.yaml`
11. Create an AtlasDatabaseUser: `kubectl apply -f config/samples/atlas_v1_atlasdatabaseuser.yaml`

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)

# How-To
## Run integration tests

**IMPORTANT: Please ensure you are in a Nix environment when running any make targets**

### make
When running the tests using `make` the Atlas credentials from `.actrc` will be used automatically to export environment
variables
```nix-bash
make int-test
```

### IDE
When running integration tests from an IDE the following environment variables need to be provided to `go test` / `ginkgo`:
`KUBEBUILDER_ASSETS=<path-to-project>/mongodb-atlas-kubernetes/testbin/bin`
`ATLAS_ORG_ID=<..>`
`ATLAS_PUBLIC_KEY=<..>`
`ATLAS_PRIVATE_KEY=<..>`
`GINKGO_EDITOR_INTEGRATION=true`
