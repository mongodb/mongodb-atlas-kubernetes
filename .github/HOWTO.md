## Req. for running act
Put the file `.actrc` to the root project folder with used secrets in GitHub

```
-s DOCKER_PASSWORD=<dd>
-s DOCKER_USERNAME=<dd>
-s REGISTRY=<docker.pkg.github.com or ghcr.io>
-s REGISTRY_USERNAME=<user-name>
-s REGISTRY_PASSWORD=<password>
```

Now simply call:
```bash
act delete
act push
act <trigger>
```