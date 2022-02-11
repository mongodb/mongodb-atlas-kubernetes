FROM golang:1.17

ENV GO111MODULE=on

# Download binaries for envtets (api-server, etcd)
RUN curl -Lo setup-envtest.sh https://raw.githubusercontent.com/kubernetes-sigs/controller-runtime/v0.8.0/hack/setup-envtest.sh && \
    mkdir -p /usr/local/kubebuilder/bin && \
    /bin/bash -c "source setup-envtest.sh && fetch_envtest_tools /usr/local/kubebuilder"

RUN go install github.com/onsi/ginkgo/v2/ginkgo@v2.1.1 && \
    go install github.com/onsi/gomega/...

COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh

ENTRYPOINT ["/home/entrypoint.sh"]
