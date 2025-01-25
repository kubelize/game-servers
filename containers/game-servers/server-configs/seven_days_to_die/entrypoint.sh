#!/bin/bash
/home/steam/steamcmd/steamcmd.sh +force_install_dir /home/steam/sdtd/ +login anonymous +app_update 294420 +quit
/home/steam/sdtd/startserver.sh -configfile=/home/steam/kubelize/serverconfig.xml