# Lets download kubectl first
FROM debian:stable AS build

# Lets install kubectl
RUN apt-get update && \
    apt-get install -y curl && \
    curl -L -o /usr/local/bin/kubectl https://storage.googleapis.com/kubernetes-release/release/v1.18.0/bin/linux/amd64/kubectl && \
    chmod +x /usr/local/bin/kubectl

# Then we will install the next tool that is kind.
RUN curl -L -o /usr/local/bin/kind https://github.com/kubernetes-sigs/kind/releases/download/v0.7.0/kind-$(uname)-amd64 && \
    chmod +x /usr/local/bin/kind

# And docker as a dependency for kind!
RUN curl -fsSL https://get.docker.com |bash

# Next lets switch to  debian slim image and only copy what we need here.
# Result image which is named stupid-kubectl will be a slim one.
FROM debian:stable AS final
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*
COPY --from=build /usr/local/bin/kubectl /usr/local/bin/kubectl
COPY --from=build /usr/local/bin/kind /usr/local/bin/kind
COPY --from=build /usr/bin/docker /usr/bin/docker
WORKDIR /workspace
ENTRYPOINT [ "/bin/bash" ]
