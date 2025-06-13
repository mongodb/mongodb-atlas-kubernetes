# TODO: Eventually replace main Dockerfile
FROM golang:1.24 AS certs-source
ARG GOTOOLCHAIN=auto

# Using rolling tag to stay on latest UBI 9
FROM registry.access.redhat.com/ubi9/ubi:latest AS ubi-certs
FROM registry.access.redhat.com/ubi9/ubi-micro:latest

ARG TARGETOS
ARG TARGETARCH
ENV TARGET_ARCH=${TARGETARCH}
ENV TARGET_OS=${TARGETOS}

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
COPY bin/${TARGET_OS}/${TARGET_ARCH}/manager .
COPY hack/licenses licenses
COPY --from=ubi-certs /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem

USER 1001:0
ENTRYPOINT ["/manager"]
