FROM golang:1.24

ENV KUBECTL_VERSION=1.33.2
ENV GO111MODULE=on

# Install kubectl
RUN curl -L "https://dl.k8s.io/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl" -o /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl

# Install controller-gen
RUN CONTROLLER_GEN_TMP_DIR=$(mktemp -d) && \
    cd $CONTROLLER_GEN_TMP_DIR && \
    go mod init tmp && \
    go install sigs.k8s.io/controller-tools/cmd/controller-gen@v0.17.2 && \
    rm -rf $CONTROLLER_GEN_TMP_DIR && \
    CONTROLLER_GEN=${GOBIN}/controller-gen

# Install operator-sdk
RUN curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v1.34.1/operator-sdk_linux_amd64 && \
    chmod +x operator-sdk_linux_amd64 && \
    mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk

COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh

ENTRYPOINT ["/home/entrypoint.sh"]
