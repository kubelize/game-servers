# Base Dockerfile
FROM kubelize/game-servers:0.2.3-ln-alpha

# Set game name
LABEL server="valheim"

# Set type
LABEL type="steam"

# Steam App ID
LABEL steam_app_id="896660"

# Set Maintainer
LABEL maintainer="kubelize@kubelize.com"

COPY /bin/valheim/entrypoint.sh /usr/local/bin/entrypoint.sh

COPY /bin/valheim/serverconfig.template /home/kubelize/server/config-data/

RUN chmod +x /usr/local/bin/entrypoint.sh && \
    chown -R kubelize:kubelize /home/kubelize/server

WORKDIR /home/kubelize/server
USER kubelize

ENTRYPOINT ["entrypoint.sh"]
