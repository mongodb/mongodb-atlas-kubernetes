FROM alpine
ARG ARCH
ENTRYPOINT ["/usr/bin/mongodb-atlas-kubernetes"]
COPY mongodb-atlas-kubernetes /usr/bin/mongodb-atlas-kubernetes
