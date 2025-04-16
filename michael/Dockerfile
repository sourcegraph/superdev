FROM golang:1.23

# Install git and repo requirements
RUN apt-get update && apt-get install -y git ffmpeg && apt-get clean

# Install Node.js 22 and npm
RUN apt-get update && \
    apt-get install -y curl gnupg2 && \
    mkdir -p /etc/apt/keyrings && \
    curl -fsSL https://deb.nodesource.com/gpgkey/nodesource-repo.gpg.key | gpg --dearmor -o /etc/apt/keyrings/nodesource.gpg && \
    echo "deb [signed-by=/etc/apt/keyrings/nodesource.gpg] https://deb.nodesource.com/node_22.x nodistro main" | tee /etc/apt/sources.list.d/nodesource.list && \
    apt-get update && \
    apt-get install -y nodejs && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Install pnpm via npm
RUN npm install -g pnpm

# Install ripgrep
RUN apt-get update && \
    apt-get install -y ripgrep && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/*

# Setup PNPM environment
RUN mkdir -p /usr/local/pnpm-global
ENV PNPM_HOME=/usr/local/pnpm-global
ENV PATH="$PNPM_HOME:$PNPM_HOME/node_modules/.bin:$PATH"
ENV SHELL=/bin/bash

# Install Sourcegraph AMP
RUN pnpm setup && pnpm add -g @sourcegraph/amp

# Set environment variables
ENV PNPM_HOME=/usr/local/pnpm-global
ENV PATH="$PNPM_HOME:$PNPM_HOME/node_modules/.bin:$PATH"
ENV SHELL=/bin/bash

RUN mkdir -p /workdir/repo
VOLUME /workdir/repo

RUN mkdir -p /workdir/guidance
VOLUME /workdir/guidance

WORKDIR /workdir

CMD [ "sh", "-c", "echo $PROMPT | amp" ]
