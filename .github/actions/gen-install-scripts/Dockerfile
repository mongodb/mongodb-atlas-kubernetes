FROM golang:1.24

ENV KUBECTL_VERSION 1.18.12
ENV GO111MODULE on

RUN go install github.com/mikefarah/yq/v4@latest

# Install
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl

RUN cd /usr/local/bin &&\
    curl -L https://raw.githubusercontent.com/kubernetes-sigs/kustomize/master/hack/install_kustomize.sh | bash

RUN CONTROLLER_GEN_TMP_DIR=$(mktemp -d) && \
    cd $CONTROLLER_GEN_TMP_DIR && \
    go mod init tmp && \
    go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.17.2 && \
    rm -rf $CONTROLLER_GEN_TMP_DIR && \
    CONTROLLER_GEN=${GOBIN}/controller-gen

RUN curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.34.1/operator-sdk_linux_amd64 && \
    chmod +x operator-sdk_linux_amd64 && \
    mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
