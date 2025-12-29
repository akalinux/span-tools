FROM mcr.microsoft.com/devcontainers/base:ubuntu

RUN --mount=type=cache,target=/var/cache/apt \
    --mount=type=cache,target=/var/lib/apt \
    apt-get update && \
    apt-get install -y vim wget curl
 
#RUN --mount=type=cache,target=/tmp \
#   if [ ! -e /tmp/go.tgz ];then wget https://go.dev/dl/go1.25.5.linux-amd64.tar.gz -O /tmp/go.tgz;fi && \
#  rm -rf /usr/local/go && tar -C /usr/local -xzf /tmp/go.tgz
  
ENV PATH=/usr/local/go/bin:$PATH
CMD ["/bin/bash"]