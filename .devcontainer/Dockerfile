FROM mcr.microsoft.com/devcontainers/go:2-1.27-trixie

ARG AIR_VERSION=v1.67.3
ARG TARGETOS
ARG TARGETARCH

# Add package repositories
RUN curl -1sLf 'https://dl.cloudsmith.io/public/evilmartians/lefthook/setup.deb.sh' | bash \
    && curl -1sLf 'https://dl.cloudsmith.io/public/task/task/setup.deb.sh' | bash \
    && apt-get update \
    && apt-get install -y --no-install-recommends \
    lefthook \
    task \
    && rm -rf /var/lib/apt/lists/*

# Install Air
RUN curl -fsSL \
    "https://github.com/air-verse/air/releases/download/${AIR_VERSION}/air_${AIR_VERSION#v}_${TARGETOS}_${TARGETARCH}" \
    -o /usr/local/bin/air \
    && chmod +x /usr/local/bin/air

# -------------------------------------------------------------
# Uber Mock Setup for non-root vscode user
# -------------------------------------------------------------
USER vscode

RUN go install go.uber.org/mock/mockgen@latest
