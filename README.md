# Work in progress

## Status

sdtd

When using a config copied into /home/steam/kubelize the server installs successfully but finally fails with the message: "Segmentation fault".

-> Todo: Figure out what causes this

Currently working on template in containers/game-servers/server-configs/seven_days_to_die/kubelize/serverconfig.template

## Changed templating tool

Switched from dockerize to gomplate. Currently there's a wierd error at "line 48, key no found".
Download the tool and test locally (it's the most effiecient way)