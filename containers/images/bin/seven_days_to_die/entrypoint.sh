#!/bin/sh

# Define paths to the template and final config
TEMPLATE_PATH="/home/kubelize/steam/config-data/serverconfig.template"
OUTPUT_PATH="/home/kubelize/steam/config-data/serverconfig.xml"

# Load values directly from the ConfigMap and Secret into temporary files
CONFIG_FILE="/home/kubelize/steam/config-data/config-values.yaml"
PASSWORD_FILE="/home/kubelize/steam/config-data/ServerPassword"

# Generate the configuration using dockerize
echo "Generating configuration from template..."

dockerize -template "$TEMPLATE_PATH:$OUTPUT_PATH" \
  # ConfigMap from the Helm-Chart
  -template="$CONFIG_FILE:/home/kubelize/steam/config-data/config-values.yaml" \
  # Secret from the Helm-Chart
  -template="$PASSWORD_FILE:/home/kubelize/steam/config-data/ServerPassword"

# Check if the config was created successfully
if [ -f "$OUTPUT_PATH" ]; then
  echo "Configuration generated successfully."
else
  echo "Failed to generate configuration."
  exit 1
fi

# Start the server or application (replace with actual command)

