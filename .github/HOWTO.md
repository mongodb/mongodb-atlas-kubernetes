# GitHub Actions
GitHub Actions help automate and customize workflows. We deploy Atlas Operator to Kubernetes using [GitHub Actions](https://docs.github.com/en/actions).

## Using GitHub Actions locally
Tools for successfully running pipeline locally:
- `act` allows running GitHub actions without pushing changes to a repository, more information [here](https://github.com/nektos/act)
- `githubsecrets` helps us to change/create [Github secrets](https://github.com/unfor19/githubsecrets) from CLI

## Req. for running act
Put the file `.actrc` to the root project folder with used secrets in GitHub

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

sample how to get config:

```bash
KUBE_CONFIG_DATA=$(kubectl config view -o json --raw | jq -c '.')
```

## Ways to run

Calling push trigger will run all workflow with `push` trigger. This command will run build-image and push-test workflows:

```bash
act push
```

This sample runs a specific job - unit-test

```bash
act -j unit-test
```

Some workflows have `workflow_dispatch` trigger - manual launch with inputs. It is possible to run it with `act` too.

```json
{
	"action":"workflow_dispatch",
	"inputs": {
		"key1":"val1",
		"key2":"val2"
	}
}
```

Build-image triggered by Pull_request:

```bash
act pull_request -e event_pull_request.json
act pull_request -j name -e event_pull_request.json
```

event_pull_request.json:

```json
{
  "pull_request": {
    "head": {
      "ref": "branch_name"
    }
  }
}
```

## Samples steps

For example, we need to run image-build
1. Prepare `.actrc` file from `.actrc.sample`
2. Find the name of job or trigger (.github/workflows/). For image build and push, the name of the job is `build-push`
3. Run job

```bash
act -j build-push
```
