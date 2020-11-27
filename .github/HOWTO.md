# GitHub Actions
GitHub Actions help to automate and customize workflows. We deploy Atlas Broker to Cloud Foundry using [GitHub Actions](https://docs.github.com/en/actions). Also, they are used for other tasks like release Atlas-OSB, cleaning our Atlas organization, testing, and demo-runs.

## Using GitHub Actions locally
Tools for successfully running pipeline locally:
- `act` allows running GitHub actions without pushing changes to a repository, more information [here](https://github.com/nektos/act)
- `githubsecrets` helps us to change/create [Github secrets](https://github.com/unfor19/githubsecrets) from CLI

## Req. for running act
Put the file `.actrc` to the root project folder with used secrets in GitHub

```
-s DOCKER_PASSWORD=<dd>
-s DOCKER_USERNAME=<dd>
-s DOCKER_REPO=<OWNER/IMAGE_NAME>
```

## Ways to run

Simply call trigger:

```bash
act push
```

This sample runs a specific job

```bash
act -j <job name>
```

Additionally, we can run workflows/jobs with different runs-on images:

```bash
act -j build -P ubuntu-latest=leori/atlas-ci:v3
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
