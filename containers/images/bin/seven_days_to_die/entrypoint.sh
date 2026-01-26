#!/bin/bash

echo "checking if the serverconfig already exists"
if [ ! -f /home/kubelize/server/sdtdconfig.xml ]; then

    echo "serverconfig.xml does not exist in server directory"

    gomplate -f /home/kubelize/config-data/serverconfig.template \
            -d config=/home/kubelize/config-data/config-values.yaml \
            -d password=/home/kubelize/config-data/serverpassword.yaml \
            -o /home/kubelize/config-data/serverconfig.xml

    echo "Rendered serverconfig"
    echo "adding the newly generated serverconfig.xml to server directory"

    cp /home/kubelize/config-data/serverconfig.xml /home/kubelize/server/sdtdconfig.xml 
    chmod +x /home/kubelize/server/sdtdconfig.xml
    echo "serverconfig.xml added to server directory"
fi

echo "installing the sdtd server"
/home/kubelize/steam/steamcmd.sh +force_install_dir /home/kubelize/server +login anonymous +app_update 294420 validate +quit
echo "installation complete"

echo "starting sdtd server"
/home/kubelize/server/startserver.sh -configfile=/home/kubelize/server/sdtdconfig.xml
