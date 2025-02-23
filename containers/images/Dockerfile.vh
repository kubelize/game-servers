# Base Dockerfile
FROM kubelize/multipaper:0.3.0-base-ln

# Set game name
LABEL server="valheim"

# Multipaper Version
LABEL steam_app_id="896660"

# Set Maintainer
LABEL maintainer="kubelize@kubelize.com"

COPY /bin/valheim/entrypoint.sh /usr/local/bin/entrypoint.sh
COPY /bin/valheim/serverconfig.template /home/kubelize/steam/config-data/
COPY /bin/valheim/configs/ /home/kubelize/steam/config-data/

RUN chmod +x /usr/local/bin/entrypoint.sh && \
    chown -R kubelize:kubelize /home/kubelize/steam

WORKDIR /home/kubelize/steam
USER kubelize

ENTRYPOINT ["entrypoint.sh"]
