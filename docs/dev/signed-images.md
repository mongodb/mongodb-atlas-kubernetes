# Image Signature Verification

Starting from version 2.2.0 and onward, all Atlas Kubernetes Operator images are signed at release time.

## Verification

The simplest way to verify an image is using the project makefile rule `verify`:

```bash
$ make
...
  verify                     Verify an AKO multi-architecture image's signature
```

It takes an `IMG` parameter and the `SIGNATURE_REPO` to use as parameters:

```bash
$ make verify IMG=mongodb/mongodb-atlas-kubernetes-operator:2.2.0 SIGNATURE_REPO=mongodb/signatures
...
VERIFIED OK
```

A successful verification run ends with a `VERIFIED OK` message. The most common error message instead is usually `Error: no matching signatures`.

Note we recommend to verify the signatures at MongoDB registry `mongodb/signatures`.

### Using cosign directly

#### Prerequisites

If you prefer to use [cosign](https://docs.sigstore.dev/signing/quickstart/) directly, you need to install it first. You can install with the *go tool* directly like this:

```bash
$ go install github.com/sigstore/cosign/v2/cmd/cosign@latest
```

Alternatively, use any of the [documented ways to install cosign](https://docs.sigstore.dev/system_config/installation/).

Another pre-requisite you need is the signing key from our team to verify the signatures against. You can get it like this:

```bash
$ curl -LO https://cosign.mongodb.com/atlas-kubernetes-operator.pem
```

The last thing you need to have is the image reference you want to verify. The `cosign` tool would prefer you pass it the image reference including the SHA to verify (eg. `mongodb/mongodb-atlas-kubernetes-operator@sha256:c7420df24f236831d21cd591c32aeafcd41787382eb093afcc2ce456c30f3a17`) but it would work also for the non SHA qualified image. Note that we sign the following SHA directly:
- The multi-architecture manifest image.
- Both the ARM and Intel (AMD64) image direct references.

#### Run verify with cosign

Once you have all ready, the verification of each image reference can be done like this:

```bash
$ COSIGN_REPOSITORY=mongodb/signatures cosign verify --insecure-ignore-tlog --key="${KEY_FILENAME}" "${IMG}" && echo PASS
```

Where:
- `KEY_FILENAME` is the name you downloaded the signature key PEM to, usually `atlas-kubernetes-operator.pem`.
- `IMG` is the image to verify, including the SHA 256 hash if possible.

Eg:
```bash
$ COSIGN_REPOSITORY=mongodb/signatures cosign verify --insecure-ignore-tlog --key=atlas-kubernetes-operator.pem mongodb/mongodb-atlas-kubernetes-operator@sha256:c7420df24f236831d21cd591c32aeafcd41787382eb093afcc2ce456c30f3a17 && echo PASS
...
PASS
```

We added the `PASS` echo on success to make it more clear when the verification passes. Again, the most common error message instead is usually `Error: no matching signatures`.
