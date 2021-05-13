# Build the manager binary
FROM golang:1.18 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# Copy DBaaSProvider config
COPY config/dbaasprovider/dbaas_provider.yaml dbaas_provider.yaml

# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download
RUN go install golang.org/x/tools/cmd/goimports@latest

# Copy the go source & git info
COPY cmd/manager/main.go cmd/manager/main.go
COPY .git/ .git/
COPY pkg/ pkg/
COPY Makefile Makefile
COPY config/ config/
COPY hack/ hack/

ARG VERSION
ARG TARGETOS
ARG TARGETARCH
ENV VERSION=${VERSION}
ENV TARGET_ARCH=${TARGETARCH}
ENV TARGET_OS=${TARGETOS}

RUN make manager

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.7

RUN microdnf install yum &&\
    yum -y update &&\
    yum clean all &&\
    dnf remove python3-pip -y &&\
    microdnf clean all

#FROM registry.access.redhat.com/ubi8/ubi
#
#RUN dnf -y update-minimal --security --sec-severity=Important --sec-severity=Critical

LABEL name="MongoDB Atlas Operator" \
      maintainer="support@mongodb.com" \
      vendor="MongoDB" \
      release="1" \
      summary="MongoDB Atlas Operator Image" \
      description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources not leaving Kubernetes cluster" \
      io.k8s.display-name="MongoDB Atlas Operator" \
      io.k8s.description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources not leaving Kubernetes cluster" \
      io.openshift.tags="mongodb,atlas" \
      io.openshift.maintainer.product="MongoDB" \
      License="Apache-2.0"


WORKDIR /
COPY --from=builder /workspace/bin/manager .
COPY hack/licenses licenses
COPY --from=builder /workspace/dbaas_provider.yaml .

USER 1001:0
ENTRYPOINT ["/manager"]
