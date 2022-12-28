FROM alpine:latest

RUN apk update && \
    apk add --no-cache curl jq

RUN curl -LO https://github.com/redhat-openshift-ecosystem/openshift-preflight/releases/download/1.4.3/preflight-linux-amd64  && \
    chmod +x ./preflight-linux-amd64 && \
    sudo mv ./preflight-linux-amd64 /usr/local/bin/preflight

COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
