FROM alpine:latest

# Copies your code file from your action repository to the filesystem path `/` of the container
COPY entrypoint.sh /home/entrypoint.sh
RUN apk update && \
    apk add jq && \
    apk add curl && \
    apk add bash && \
    apk add --update coreutils && \
    rm -rf /var/cache/apk/*
RUN chmod +x /home/entrypoint.sh
# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
