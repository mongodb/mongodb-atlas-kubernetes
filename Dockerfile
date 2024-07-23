FROM nixos/nix:latest AS go-builder
COPY Makefile Makefile
COPY shell.nix shell.nix

# Set the working directory
WORKDIR /workspace

# Copy necessary files
COPY . .

# Set build arguments
ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ENV VERSION=${VERSION}
ENV TARGET_ARCH=${TARGETARCH}
ENV TARGET_OS=${TARGETOS}

# Install Go tools and build the manager binary using nix-shell
RUN nix-shell shell.nix --run "go mod download"
RUN nix-shell shell.nix --run "go install golang.org/x/tools/cmd/goimports@latest"
RUN nix-shell shell.nix --run "make manager"

# Stage 3: Final image
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

# Copy licenses
COPY hack/licenses /licenses

# Set user
USER 1001:0

# Set the entrypoint
ENTRYPOINT ["/manager"]