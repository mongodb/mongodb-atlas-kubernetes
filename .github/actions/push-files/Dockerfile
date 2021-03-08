FROM alpine/git:latest

# Install GitHub CLI
RUN apk update && \
    apk add --no-cache libc6-compat bash
RUN mkdir ghcli && cd ghcli && \
    wget https://github.com/cli/cli/releases/download/v1.5.0/gh_1.5.0_linux_386.tar.gz -O ghcli.tar.gz --no-check-certificate  && \
    tar --strip-components=1 -xf ghcli.tar.gz -C /usr/local

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
