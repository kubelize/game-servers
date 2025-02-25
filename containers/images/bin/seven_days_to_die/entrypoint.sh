#!/bin/bash
echo "Rendering serverconfig.xml"
gomplate -f /home/kubelize/steam/config-data/serverconfig.template \
         -d config=/home/kubelize/steam/config-data/config-values.yaml \
         -d password=/home/kubelize/steam/config-data/serverpassword.yaml \
         -o /home/kubelize/steam/config-data/serverconfig.xml
echo "Rendering of the serverconfig complete"

echo "Checking if the serverconfig already exists"
if [ ! -f /home/kubelize/server/serverconfig.xml ]; then
    echo "Moving the newly generated serverconfig.xml to server directory"
    mv /home/kubelize/steam/config-data/serverconfig.xml /home/kubelize/server/
fi

echo "Installing Seven Days to Die Dedicated Server"
/home/kubelize/steam/steamcmd.sh +force_install_dir /home/kubelize/server +login anonymous +app_update 294420 validate +quit
echo "Installation of Seven Days to Die Dedicated Server complete"

echo "Starting Seven Days to Die Dedicated Server"
/home/kubelize/server/startserver.sh -configfile=/home/kubelize/server/serverconfig.xml
