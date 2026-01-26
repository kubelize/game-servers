# Base Dockerfile
FROM kubelize/game-servers:0.3.8-ln-alpha

# Set game name
LABEL server="valheim"

# Set type
LABEL type="steam"

# Steam App ID
LABEL steam_app_id="896660"

# Set Maintainer
LABEL maintainer="kubelize@kubelize.com"

# Copy server config template to protected location (not overwritten by PVC mounts)
COPY images/bin/valheim/serverconfig.template /usr/local/share/game-templates/serverconfig.template

RUN mkdir -p /home/kubelize/config-data && \
    chown -R kubelize:kubelize /home/kubelize/server /home/kubelize/config-data

WORKDIR /home/kubelize/server
USER kubelize

ENTRYPOINT ["gamekeeper", "start", "--game", "valheim"]
