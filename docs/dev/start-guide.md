# Onboarding to Atlas Operator
1. Install [Nix](https://nixos.org/download/) and set up with the IDE (instructions below)
2. Run the Nix shell (by entering ```shell.nix``` in the bash to initialize an environment with pre-installed dependencies). For more details, refer to the shell.nix file.
3. Install optional dependencies
    ```
    # on Mac
    brew install coreutils # or https://www.gnu.org/software/coreutils/
    brew install pre-commit # or https://pre-commit.com/index.html#install
    pre-commit install # from the root of the project
    ```
4. Install [Kind](https://kind.sigs.k8s.io/docs/user/quick-start/#installation)
4. Clone the project to your workspace (note, that this doesn't need to be `GOPATH` as the project uses Go Modules)
5. Copy the default Github Actions settings for local run: `cp .actrc.local.sample .actrc`
6. Copy the default Github Actions environment for local run: `cp dotenv.sample .env`
7. Update the .actrc - specify your Atlas connectivity data (orgId, keys)
8. Build and deploy the Operator into the K8s cluster: `make deploy`
9. Create an AtlasProject: `kubectl apply -f config/samples/atlas_v1_atlasproject.yaml` (note, that Atlas connection secrets are created during running `make deploy`)
10. Create an AtlasDeployment: `kubectl apply -f config/samples/atlas_v1_atlasdeployment.yaml`
11. Create an AtlasDatabaseUser: `kubectl apply -f config/samples/atlas_v1_atlasdatabaseuser.yaml`

Some more details about using `act` can be found in [HOWTO.md](../../.github/HOWTO.md)
### IDE setup with Nix
Using Nix Package Manager (Visual Studio Code)
1. Install the [Nix package manager](https://nixos.org/) by following the instructions on the Nix website.
2. Install the [Nix Environment Selector](https://marketplace.visualstudio.com/items?itemName=arrterian.nix-env-selector)  extension.
3. Restart VSCode to ensure that nix-shell is in the system PATH.
4. Create a Nix environment configuration file (default.nix or shell.nix) in the root of your project's workspace.
5. Use the keyboard shortcut `Ctrl + Shift + P` to open the Command Palette in VSCode
6. Run the Nix-Env: Select Environment command from the Command Palette and choose the Nix environment you'd like to apply.
7. Allow the environment to build, which may take some time depending on the dependencies.
8. Restart VSCode to apply the built environment and ensure all dependencies are properly loaded.

Using the terminal
1. Install the [Nix package manager](https://nixos.org/) by following the instructions on the Nix website.
2. Ensure that the environment's launch command is available in your command line interface (CLI) to start the IDE  from the terminal.
4. Use the cd command to navigate to the directory containing your `shell.nix` file, which should define the development environment for your project:
5. Activate the development environment by running: `nix-shell` (This command loads all necessary dependencies and prepares your environment for development.)
6. Launch your IDE directly from the command line. For example, if using Visual Studio Code, type: `code`
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
