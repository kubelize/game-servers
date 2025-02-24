#!/bin/bash
echo "Rendering serverconfig.xml"
gomplate -f /home/kubelize/steam/config-data/serverconfig.template \
         -d config=/home/kubelize/steam/config-data/config-values.yaml \
         -o /home/kubelize/steam/config-data/serverconfig.xml
