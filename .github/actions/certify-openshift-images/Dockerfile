FROM registry.access.redhat.com/ubi8/ubi:latest

RUN yum install -y \
  bzip2 \
  gzip \
  tar \
  iptables \
  yum-utils \
  jq

RUN yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo && \
  yum list docker-ce --showduplicates | sort -r && \
  yum install -y \
  docker-ce \
  docker-ce-cli \
  containerd.io

RUN yum clean all

RUN curl -LO https://github.com/redhat-openshift-ecosystem/openshift-preflight/releases/latest/download/preflight-linux-amd64  && \
  chmod +x ./preflight-linux-amd64 && \
  mv ./preflight-linux-amd64 /usr/local/bin/preflight

COPY ./entrypoint.sh /home/entrypoint.sh
RUN chmod +x /home/entrypoint.sh

# Code file to execute when the docker container starts up (`entrypoint.sh`)
ENTRYPOINT ["/home/entrypoint.sh"]
