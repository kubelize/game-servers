#!/bin/bash
set -e

# Set default values
MINECRAFT_VERSION=${MINECRAFT_VERSION:-"1.20.4"}
SERVER_JAR=${SERVER_JAR:-"server.jar"}
MIN_RAM=${MIN_RAM:-"1024M"}
MAX_RAM=${MAX_RAM:-"2048M"}
SERVER_PORT=${SERVER_PORT:-"25565"}
ONLINE_MODE=${ONLINE_MODE:-"true"}
DIFFICULTY=${DIFFICULTY:-"normal"}
MAX_PLAYERS=${MAX_PLAYERS:-"20"}
MOTD=${MOTD:-"A Minecraft Server"}
LEVEL_NAME=${LEVEL_NAME:-"world"}
GAMEMODE=${GAMEMODE:-"survival"}
PVP=${PVP:-"true"}

GAME_DIR="/home/kubelize/gameserver"
cd $GAME_DIR

# Download server jar if it doesn't exist
if [ ! -f "$SERVER_JAR" ]; then
    echo "Downloading Minecraft server version $MINECRAFT_VERSION..."
    wget -O $SERVER_JAR "https://launcher.mojang.com/v1/objects/${MINECRAFT_VERSION}/${SERVER_JAR}" || \
    echo "Note: You may need to manually download the server jar and mount it"
fi

# Accept EULA
if [ ! -f "eula.txt" ]; then
    echo "eula=true" > eula.txt
fi

# Generate server.properties if template exists
if [ -f "/home/kubelize/gameserver/config-data/server.properties.template" ]; then
    echo "Generating server.properties from template..."
    gomplate -f /home/kubelize/gameserver/config-data/server.properties.template -o server.properties
fi

# Start the server
echo "Starting Minecraft server..."
exec java -Xms${MIN_RAM} -Xmx${MAX_RAM} -jar ${SERVER_JAR} nogui
