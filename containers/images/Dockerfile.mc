# Base Dockerfile
FROM kubelize/game-servers:0.3.7-jv-alpha

# Set game name
LABEL server="minecraft"

# Set type
LABEL type="java"

# Set Maintainer
LABEL maintainer="kubelize@kubelize.com"

# Install dependencies
RUN apt update && apt upgrade -y && \
    apt install --no-install-recommends --no-install-suggests -y \
    openjdk-17-jre-headless \
    wget

RUN apt remove --purge -y curl && \
    apt autoremove -y && \
    apt clean && \
    rm -rf /var/lib/apt/lists/* /tmp/* /var/tmp/*

# Copy server config template to protected location (not overwritten by PVC mounts)
COPY images/bin/minecraft/serverconfig.template /usr/local/share/game-templates/serverconfig.template

RUN chown -R kubelize:kubelize /home/kubelize/server

WORKDIR /home/kubelize/server
USER kubelize

ENTRYPOINT ["gamekeeper", "start", "--game", "minecraft"]
