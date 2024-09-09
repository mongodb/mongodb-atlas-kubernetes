FROM alpine
FROM registry.access.redhat.com/ubi9/ubi:9.2 as ubi-certs
FROM registry.access.redhat.com/ubi9/ubi-micro:9.2
ARG ARCH
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
COPY --from=ubi-certs /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem /etc/pki/ca-trust/extracted/pem/tls-ca-bundle.pem
USER 1001:0
ENTRYPOINT ["/usr/bin/mongodb-atlas-kubernetes"]
COPY mongodb-atlas-kubernetes /usr/bin/mongodb-atlas-kubernetes
