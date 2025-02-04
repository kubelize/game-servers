#!/bin/sh

# Define paths to the template and final config
TEMPLATE_PATH="/config/serverconfig.template"
OUTPUT_PATH="/config//serverconfig.xml"

# Load values directly from the ConfigMap and Secret into temporary files
CONFIG_FILE="/config/config-values.yaml"
PASSWORD_FILE="/config/ServerPassword"

# Generate the configuration using dockerize
echo "Generating configuration from template..."

dockerize -template "$TEMPLATE_PATH:$OUTPUT_PATH" \
  -template="$CONFIG_FILE:/config/config-values.yaml" \
  -template="$PASSWORD_FILE:/config/ServerPassword"

# Check if the config was created successfully
if [ -f "$OUTPUT_PATH" ]; then
  echo "Configuration generated successfully."
else
  echo "Failed to generate configuration."
  exit 1
fi

# Start the server or application (replace with actual command)

