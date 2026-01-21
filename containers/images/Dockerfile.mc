# Base Dockerfile
FROM kubelize/game-servers:0.2.3-jv-alpha

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

COPY /bin/minecraft/entrypoint.sh /usr/local/bin/entrypoint.sh

COPY /bin/minecraft/serverconfig.template /home/kubelize/server/config-data/

RUN chmod +x /usr/local/bin/entrypoint.sh && \
    chown -R kubelize:kubelize /home/kubelize/server

WORKDIR /home/kubelize/server
USER kubelize

ENTRYPOINT ["entrypoint.sh"]
