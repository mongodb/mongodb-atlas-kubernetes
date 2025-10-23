FROM golang:1.25

ENV KUBECTL_VERSION=1.34.1
ENV KUSTOMIZE_VERSION=5.7.1
ENV CONTROLLER_GEN_VERSION=0.17.2
ENV OPERATOR_SDK_VERSION=1.39.2
ENV GO111MODULE=on

# --- System Packages ---
RUN apt-get update && \
    apt-get install -y jq && \
    rm -rf /var/lib/apt/lists/*

# --- Go Tools ---
# Install yq
RUN go install github.com/mikefarah/yq/v4@latest
# Install controller-gen
RUN CONTROLLER_GEN_TMP_DIR=$(mktemp -d); \
    cd $CONTROLLER_GEN_TMP_DIR; \
    go mod init tmp; \
    go install sigs.k8s.io/controller-tools/cmd/controller-gen@v${CONTROLLER_GEN_VERSION}; \
    rm -rf $CONTROLLER_GEN_TMP_DIR; \
    CONTROLLER_GEN=${GOBIN}/controller-gen

# --- Prebuilt Binaries ---
# Install kubectl
RUN curl -L https://storage.googleapis.com/kubernetes-release/release/v${KUBECTL_VERSION}/bin/linux/amd64/kubectl -o /usr/bin/kubectl && \
    chmod +x /usr/bin/kubectl
# Install kustomize
RUN curl -L https://github.com/kubernetes-sigs/kustomize/releases/download/kustomize%2Fv${KUSTOMIZE_VERSION}/kustomize_v${KUSTOMIZE_VERSION}_linux_amd64.tar.gz \
    -o /tmp/kustomize.tar.gz && \
    tar -xzvf /tmp/kustomize.tar.gz -C /usr/local/bin && \
    chmod +x /usr/local/bin/kustomize
# Install operator-sdk
RUN curl -LO https://github.com/operator-framework/operator-sdk/releases/download/v${OPERATOR_SDK_VERSION}/operator-sdk_linux_amd64 && \
    chmod +x operator-sdk_linux_amd64 && \
    mv operator-sdk_linux_amd64 /usr/local/bin/operator-sdk


COPY . /home
COPY .github/actions/gen-install-scripts/entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh

ENTRYPOINT ["/home/entrypoint.sh"]
