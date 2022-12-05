FROM alpine:latest

RUN apk update && \
    apk add --no-cache libc6-compat curl tar git sed

RUN curl -Lo tmp/ghcli.tar.gz https://github.com/cli/cli/releases/download/v2.20.2/gh_2.20.2_linux_386.tar.gz && \
    tar --strip-components=1 -xf tmp/ghcli.tar.gz -C /usr/local && \
    rm tmp/ghcli.tar.gz

COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]