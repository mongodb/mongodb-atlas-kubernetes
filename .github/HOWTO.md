# GitHub Actions
GitHub Actions help automate and customize workflows. We deploy Atlas Operator to Kubernetes using [GitHub Actions](https://docs.github.com/en/actions).

## Using GitHub Actions locally (Optional)

Developers can optionally use `act` to run GitHub Actions workflows locally. This is not required for development and is maintained on a per-developer basis.

### Setting up act (Optional)

If you want to use `act` for local testing:

1. Install `act` from [nektos/act](https://github.com/nektos/act)
2. Create your own `.actrc` file in the project root with your secrets:

```
# Update this data with your cloud-qa custom data
-s ATLAS_ORG_ID=<id>
-s ATLAS_PUBLIC_KEY=<public_key>
-s ATLAS_PRIVATE_KEY=<private_key>
# Push to Docker Registry
-s DOCKER_USERNAME=<username>
-s DOCKER_PASSWORD=<password>
-s DOCKER_REPO=owner/repo_name
-s DOCKER_REGISTRY=docker.io
-s KUBE_CONFIG_DATA=<copy of kubeconfig>
```

Sample how to get config:

```bash
KUBE_CONFIG_DATA=$(kubectl config view -o json --raw | jq -c '.')
```

3. Add `.actrc` to your `.gitignore` to avoid committing secrets

### Running workflows with act

```bash
# Run all workflows with push trigger
act push

# Run a specific job
act -j unit-test

# Run with workflow_dispatch trigger
act workflow_dispatch -e event.json
```

**Note**: The project does not maintain `.actrc` files or provide templates. Each developer is responsible for their own `act` setup if they choose to use it.
