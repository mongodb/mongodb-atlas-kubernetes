#TODO change to alpine
FROM ubuntu:18.04

# Install GitHub CLI
RUN apt-get update && \
    apt-get install -y --no-install-recommends git wget
RUN wget https://github.com/cli/cli/releases/download/v1.0.0/gh_1.0.0_linux_386.tar.gz -O ghcli.tar.gz --no-check-certificate  && \
    tar --strip-components=1 -xf ghcli.tar.gz

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
