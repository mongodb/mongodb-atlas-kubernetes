# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source & git info
COPY main.go main.go
COPY .git/ .git/
COPY pkg/ pkg/

ARG VERSION
ENV PRODUCT_VERSION=${VERSION}

# Build
RUN if [ -z $PRODUCT_VERSION ]; then PRODUCT_VERSION=$(git describe --tags); fi; \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on \
    go build -a -ldflags="-X main.version=$PRODUCT_VERSION" -o manager main.go

FROM registry.access.redhat.com/ubi8/ubi-minimal:8.4-205.1626828526

RUN microdnf install yum &&\
    yum -y update-minimal --security --sec-severity=Important --sec-severity=Critical &&\
    yum clean all &&\
    microdnf remove yum &&\
    microdnf clean all

#FROM registry.access.redhat.com/ubi8/ubi
#
#RUN dnf -y update-minimal --security --sec-severity=Important --sec-severity=Critical

LABEL name="MongoDB Atlas Operator" \
      maintainer="support@mongodb.com" \
      vendor="MongoDB" \
      release="1" \
      summary="MongoDB Atlas Operator Image" \
      description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources not leaving Kubernetes cluster"

WORKDIR /
COPY --from=builder /workspace/manager .
COPY hack/licenses licenses

USER 1001:0
ENTRYPOINT ["/manager"]
