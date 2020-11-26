## Req. for running act
Put the file `.actrc` to the root project folder with used secrets in GitHub

```
-s DOCKER_PASSWORD=<dd>
-s DOCKER_USERNAME=<dd>
-s DOCKER_REPO=<OWNER/IMAGE_NAME>
```

Now simply call trigger:
```bash
act delete
act push
act <trigger>
```

This sample runs a specific job

```bash
act -j <job name> #call job
```

Additionally, we can run workflows/jobs with different runs-on images:

```bash
act -j build -P ubuntu-latest=leori/atlas-ci:v2
