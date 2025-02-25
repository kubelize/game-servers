#!/bin/bash
gomplate -f /Users/dan/Git/kubelize/game-servers/test/serverconfig.template \
         -d config=/Users/dan/Git/kubelize/game-servers/test/config-test.yaml \
         -d password=/Users/dan/Git/kubelize/game-servers/test/serverpassword.yaml \
         -o /Users/dan/Git/kubelize/game-servers/test/serverconfig.xml