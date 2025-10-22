# Custom Go Nix Flake

## Why a flake?

Nix tends to be several minor versions behind Go's official releases.

This posses chanllenges in two situations:

1. When Go moves to a new majoer version
1. When Go has a vulnerability on the current latest Nix version, fixed by a newer official release

For major versions it is usually fine to wait for Nix to have a major version compilation avaiable for devbox to use. This is because, no matter how early we may want to upgrade, many go tools we depend on, such as licence checking or linting or Kubernetes libraries such as `controller-runtime`, usually need some time to catch up with the major release anyways. By the time they support the new major version, there is usually a Nix compilation of the new Go release, at least in the unstable channel.

For minor versions, it can be more problematic. If the latest Nix available release is compromised, it mgiht take a few days or weeks for the new version to become available in Nix. On the other hand, Go only marks a vulnerable release after releasing the fixed version.

In other words, we need to be able to move to the latest Go release as needed, specially to avoid vulvnerabilities within the same major version.

## How

The current flake in this directory will download and install the pre-compiled binaries straight from https://go.dev/dl, that is the official Go downloads site. It only supports 2 platforms:
- `x86_64-linux` for the CI and Linux developers.
- `aarch64-darwin` for developers working on MacOS.

The flake derivation does not build anything, just unpacks and places the binaries where expected to be used by the resulting flake.

## Updating

The flake is pinned to a particular Go point release. To bump the downloaded binary you have to:

1. Bump the `goVersion` variable. E.g. `goVersion = "1.25.3";` -> `goVersion = "1.25.4";`
2. Replace both `sha256` variable values with the correct ones for the new downloaded file.

One easy way to read the expected sha 256 hash to be used for each `sha-256` setting is to using `nix-prefetch-url` or `nix store prefetch-file --json` to grab the file and hash it.

For example:

```shell
$ nix store prefetch-file --json https://go.dev/dl/go1.25.3.linux-amd64.tar.gz |jq -r .hash
sha256-AzXzFLbnv+CMPQz6p8GduWG3uZ+yC+YrCoJsmSrRTg8=
```

Make sure to use the correct architecture filename download to grab its corresponding sha 256 hash.

## Testing

Using `devbox shell` normally would already grab and build the flake, as referenced by devbox.json entry `"path:./flakes/go": {}`. Still if you want to test the flake buil in isolation you can run (in this directory):

``shell
nix build .
```

On success a `result` entry in teh directory soft links to the built flake.
