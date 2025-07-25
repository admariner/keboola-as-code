FROM golang:1.24.4

ENV HOME=/my-home
ENV GOCACHE=/tmp/cache/go
ENV GOMODCACHE=/tmp/cache/go-mod
ENV GOBIN=$GOPATH/bin
ENV GOFLAGS="-mod=mod"
ENV GOLANGCI_LINT_CACHE=/tmp/cache/golangci-lint
ENV PATH="$PATH:$GOBIN"

# Install packages
RUN apt-get update && apt-get install --no-install-recommends --yes zip unzip nano libprotobuf-dev graphviz build-essential
ENV EDITOR=nano

# Install Protobuf
RUN curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v30.2/protoc-30.2-linux-x86_64.zip && \
    unzip protoc-30.2-linux-x86_64.zip -d /usr/local && \
    rm -f protoc-30.2-linux-x86_64.zip

# Download toxiproxy
RUN curl -L https://github.com/Shopify/toxiproxy/releases/download/v2.12.0/toxiproxy-server-linux-amd64 -o /usr/local/bin/toxiproxy-server
RUN chmod +x /usr/local/bin/toxiproxy-server

# Install Task
RUN sh -c "$(curl --location https://taskfile.dev/install.sh)" -- -d

# Install tools
RUN mkdir -p /tmp/build
COPY Taskfile.yml /tmp/build/Taskfile.yml
COPY scripts  /tmp/build/scripts
RUN cd /tmp/build && \
    task install-golangci-lint && \
    task install-goreleaser && \
    task install-protoc-gen-go && \
    cd / && \
    rm -rf /tmp/build && \
    go clean -cache -modcache

# Install envsubstr and helm
RUN curl -L https://github.com/a8m/envsubst/releases/download/v1.4.3/envsubst-$(uname -s)-$(uname -m) -o /usr/local/bin/envsubst && \
    chmod +x /usr/local/bin/envsubst && \
    curl https://raw.githubusercontent.com/helm/helm/main/scripts/get-helm-3 | bash

# Set prompt
RUN mkdir -p ~ && \
    echo 'PS1="\w > "' > ~/.bashrc

# Fix permissions
RUN chmod -R 1777 $GOPATH && \
    mkdir -p $GOCACHE && chmod -R 1777 $GOCACHE && \
    mkdir -p $GOMODCACHE && chmod -R 1777 $GOMODCACHE && \
    mkdir -p $GOLANGCI_LINT_CACHE && chmod -R 1777 $GOLANGCI_LINT_CACHE && \
    mkdir -p $HOME && chmod -R 1777 $HOME

WORKDIR /code/
CMD ["/bin/bash"]
