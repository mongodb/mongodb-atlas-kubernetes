# Onboarding to Atlas Operator
1. Install devbox by following the instructions provided on the official [Jetify website](https://www.jetify.com/devbox).
2. Run the Devbox shell (by entering ```devbox shell``` in the bash to initialize an environment with pre-installed dependencies). For more details, refer to the devbox.json file.
3. Install optional dependencies
    ```
    # on Mac
    brew install coreutils # or https://www.gnu.org/software/coreutils/
    brew install pre-commit # or https://pre-commit.com/index.html#install
    pre-commit install # from the root of the project
    ```
4. Install [Kind] (https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
4. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
5. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
6. Copy the default Github Actions environment for local run: `cp dotenv.sample .env`
7. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
8. Build and deploy the Operator into the K8s cluster: `make deploy`
9. Create an AtlasProject: `kubectl apply -f config/samples/atlas_v1_atlasproject.yaml` (note, that Atlas connection secrets are created during running `make deploy`)
10. Create an AtlasDeployment: `kubectl apply -f config/samples/atlas_v1_atlasdeployment.yaml`
11. Create an AtlasDatabaseUser: `kubectl apply -f config/samples/atlas_v1_atlasdatabaseuser.yaml`

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)
### IDE setup with Devbox
Using Direnv Environment Extension (Visual Studio Code)
1. Install devbox by following the instructions provided on the official [Jetify website](https://www.jetify.com/devbox).
2. Restart Visual Studio Code to ensure that Devbox is included in the system PATH.
3. Install the [Direnv Environment Extension](https://marketplace.visualstudio.com/items?itemName=mkhl.direnv) extension for Visual Studio Code
4. If not present, create a devbox.json environment configuration file in the root of your project's workspace.
5. Run the following command to generate a direnv configuration:  `devbox generate direnv`
6. Use the keyboard shortcut  `Ctrl + Shift + P` to open the Command Palette .
7. Run the direnv:  `Load .envrc file ` command from the Command Palette followed by  `direnv: Reload environment `.
8. Wait for the environment to build, as this may take some time depending on the dependencies.
9. Restart Visual Studio Code to apply the built environment and ensure all dependencies are properly loaded.

Using the terminal(Visual Studio Code)
1. Install devbox by following the instructions provided on the official [Jetify website](https://www.jetify.com/devbox).
2. Ensure that the environment's launch command is available in your command line interface (CLI) to start the IDE from the terminal.
3. Use the cd command to navigate to the directory containing your `devbox.json` file, which should define the development environment for your project:
4. Activate the development environment by running: `devbox shell` (This command loads all necessary dependencies and prepares your environment for development).
5. Launch your IDE directly from the command line. For example, if using Visual Studio Code, type: `code`.

Using a devcontainer
1. Install devbox by following the instructions provided on the official [Jetify website](https://www.jetify.com/devbox).
2. If not present, create a `devbox.json` environment configuration file in the root of your project's workspace.
3. Run the following command to generate a direnv configuration: `devbox generate devcontainer`
4. Install the IDE plugin for running in devcontainers.
5. Open the project using the generated devcontainer configuration. You can do this by selecting `Reopen in Container` when prompted.


# How-To
## Run integration tests

**IMPORTANT: Please ensure you are in a Devbox environment when running any make targets**

### make
When running the tests using `make` the Atlas credentials from `.actrc` will be used automatically to export environment
variables
```devbox shell
make int-test
```

### IDE
When running integration tests from an IDE the following environment variables need to be provided to `go test` / `ginkgo`:
`KUBEBUILDER_ASSETS=<path-to-project>/mongodb-atlas-kubernetes/testbin/bin`
`ATLAS_ORG_ID=<..>`
`ATLAS_PUBLIC_KEY=<..>`
`ATLAS_PRIVATE_KEY=<..>`
`GINKGO_EDITOR_INTEGRATION=true`
