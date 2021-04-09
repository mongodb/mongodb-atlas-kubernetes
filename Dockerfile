# Build the manager binary
FROM golang:1.15 as builder

WORKDIR /workspace
# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum
# cache deps before building and copying source so that we don't need to re-download as much
# and so that source changes don't invalidate our downloaded layer
RUN go mod download

# Copy the go source
COPY main.go main.go
COPY pkg/ pkg/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -a -o manager main.go

FROM registry.access.redhat.com/ubi8/ubi

RUN yum -y update-minimal --security --sec-severity=Important --sec-severity=Critical

LABEL name="MongoDB Atlas Operator" \
      maintainer="support@mongodb.com" \
      vendor="MongoDB" \
      release="1" \
      summary="MongoDB Atlas Operator Image" \
      description="MongoDB Atlas Operator is a Kubernetes Operator allowing to manage MongoDB Atlas resources not leaving Kubernetes cluster"

WORKDIR /
COPY --from=builder /workspace/manager .
COPY hack/licenses licenses

USER nonroot:nonroot
ENTRYPOINT ["/manager"]
