# Base Dockerfile
FROM kubelize/game-servers:0.1.0-baseln

# Set game name
LABEL game="valheim"

# Set Maintainer
LABEL maintainer="kubelize@kubelize.com"

# Set Steam App ID
LABEL steam_app_id="896660"

COPY /server-configs/valheim/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY /server-configs/valheim/kubelize/ /home/steam/kubelize/

RUN chmod +x /usr/local/bin/entrypoint.sh && \
    chown -R steam:steam /home/steam/kubelize/

WORKDIR /home/steam/steamcmd
USER steam

ENTRYPOINT ["entrypoint.sh"]
