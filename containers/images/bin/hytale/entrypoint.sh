#!/bin/bash
set -e

GAME_DIR="/home/kubelize/gameserver"
cd $GAME_DIR

# Download the hytale downloader binary if it doesn't exist
if [ ! -f "./hytale_downloader" ]; then
    echo "Downloading Hytale downloader..."
    HYTALE_DOWNLOADER_URL="${HYTALE_DOWNLOADER_URL:-https://drive.kubelize.com/public.php/dav/files/HJqqWZx5522wnoT}"
    wget -O ./hytale_downloader "$HYTALE_DOWNLOADER_URL"
    chmod +x ./hytale_downloader
fi