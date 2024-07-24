FROM nixos/nix:2.23.3 AS builder

# Set the working directory
WORKDIR /workspace

# Copy the Go source files and other required files
COPY cmd/manager/main.go cmd/manager/main.go
COPY .git/ .git/
COPY pkg/ pkg/
COPY internal/ internal/
COPY Makefile Makefile
COPY config/ config/
COPY hack/ hack/
COPY shell.nix shell.nix

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Set build arguments
ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ENV VERSION=${VERSION}
ENV TARGET_ARCH=${TARGETARCH}
ENV TARGET_OS=${TARGETOS}

# Update the nix channel
RUN nix-channel --update

# Install Go tools and build the manager binary using nix-shell
RUN nix-shell shell.nix --run "go mod download && go install golang.org/x/tools/cmd/goimports@latest && make manager"

# Final image
FROM registry.access.redhat.com/ubi9/ubi:9.2 as ubi-certs
FROM registry.access.redhat.com/ubi9/ubi-micro:9.2

# Metadata
LABEL name="MongoDB Atlas Operator" \
      maintainer="support@mongodb.com" \
      vendor="MongoDB" \
      release="1" \
      summary="MongoDB Atlas Operator Image" \
      description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources without leaving the Kubernetes cluster" \
      io.k8s.display-name="MongoDB Atlas Operator" \
      io.k8s.description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources without leaving the Kubernetes cluster" \
      io.openshift.tags="mongodb,atlas" \
      io.openshift.maintainer.product="MongoDB" \
      License="Apache-2.0"

# Set the working directory
WORKDIR /

# Copy the built binary from the appropriate builder stage
COPY --from=go-builder /workspace/bin/manager /manager

# Copy certificates
COPY --from=ubi-certs /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem

# Copy licenses
COPY hack/licenses /licenses

# Set user
USER 1001:0

# Set the entrypoint
ENTRYPOINT ["/manager"]
