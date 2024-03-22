# TODO: Eventually replace main Dockerfile
FROM golang:1.22 as certs-source

FROM registry.access.redhat.com/ubi9/ubi-micro:9.2

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
COPY --from=certs-source /etc/ssl/certs/ca-certificates.crt /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem

USER 1001:0
ENTRYPOINT ["/manager"]
