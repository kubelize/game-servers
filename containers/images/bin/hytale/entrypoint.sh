#!/bin/bash
set -e

# ============================================
# Hytale Server Entrypoint Script
# Using gomplate for configuration rendering
# ============================================

# Force line buffering for better logging in Kubernetes
export PYTHONUNBUFFERED=1

# Set directory paths
BASE_DIR="${BASE_DIR:-/home/kubelize/server}"
DATA_DIR="$BASE_DIR/data"
GAME_DIR="$DATA_DIR"
CONFIG_DATA_DIR="$BASE_DIR/config-data"
SERVER_JAR_PATH="$GAME_DIR/Server/HytaleServer.jar"
HYTALE_DOWNLOADER="$GAME_DIR/hytale-downloader"
ASSETS_ZIP="$GAME_DIR/Assets.zip"
CONFIG_VALUES_FILE="$CONFIG_DATA_DIR/config-values.yaml"

# Source environment variables from ConfigMap for server options (not config.json)
if [ -f "$CONFIG_VALUES_FILE" ]; then
    echo "Loading server options from ConfigMap..."
    # Export vars needed for server command-line options
    while IFS=': ' read -r key value; do
        key=$(echo "$key" | tr -d ' ')
        value=$(echo "$value" | tr -d '"' | tr -d "'" | xargs)
        [[ -z "$key" || "$key" =~ ^# ]] && continue
        export "$key=$value"
    done < "$CONFIG_VALUES_FILE"
fi

# Color codes for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color
BOLD='\033[1m'
DIM='\033[2m'

log_section() {
    echo ""
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
    echo -e "${BOLD}${CYAN}  $1${NC}"
    echo -e "${BOLD}${CYAN}═══════════════════════════════════════════════════════════${NC}"
}

log_step() {
    echo -ne "  ${CYAN}→${NC} $1: "
}

log_success() {
    echo -e "${GREEN}✓${NC}"
}

log_error() {
    echo -e "${RED}✗ $1${NC}"
}

log_info() {
    echo -e "  ${CYAN}ℹ${NC} $1"
}

# ============================================
# Step 1: Directory Setup
# ============================================
log_section "Environment Setup"

log_step "Creating directories"
mkdir -p "$DATA_DIR"
mkdir -p "$CONFIG_DATA_DIR"
cd "$BASE_DIR"
log_success

log_step "Working directory"
echo -e "${GREEN}$GAME_DIR${NC}"

# ============================================
# Step 2: Download Hytale Downloader
# ============================================
log_section "Hytale Downloader Setup"

if [ ! -f "$HYTALE_DOWNLOADER" ]; then
    log_step "Downloading hytale-downloader"
    
    # Download from official Hytale CDN
    HYTALE_DOWNLOADER_URL="${HYTALE_DOWNLOADER_URL:-https://drive.kubelize.com/public.php/dav/files/HJqqWZx5522wnoT}"
    
    if wget -q -O "$HYTALE_DOWNLOADER" "$HYTALE_DOWNLOADER_URL"; then
        chmod +x "$HYTALE_DOWNLOADER"
        log_success
    else
        log_error "Failed to download hytale-downloader"
        exit 1
    fi
else
    log_step "hytale-downloader binary"
    echo -e "${GREEN}already exists${NC}"
fi

# ============================================
# Step 3: Download/Update Hytale Server
# ============================================
log_section "Hytale Server Installation"

# Enable auto-update by default
HYTALE_AUTO_UPDATE="${HYTALE_AUTO_UPDATE:-TRUE}"

# Determine if we need to run the downloader
SHOULD_DOWNLOAD=false

if [ ! -f "$SERVER_JAR_PATH" ] || [ ! -f "$ASSETS_ZIP" ]; then
    log_info "Server files not found. Initiating download..."
    SHOULD_DOWNLOAD=true
elif [ "$HYTALE_AUTO_UPDATE" = "TRUE" ]; then
    log_info "Checking for game file updates..."
    SHOULD_DOWNLOAD=true
else
    log_step "Hytale server files"
    echo -e "${GREEN}already downloaded (auto-update disabled)${NC}"
fi

if [ "$SHOULD_DOWNLOAD" = true ]; then
    echo ""
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo -e "${YELLOW}  AUTHENTICATION REQUIRED${NC}"
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    echo -e "The Hytale server requires authentication to download game files."
    echo -e "Please follow the instructions below:"
    echo ""
    echo -e "1. The downloader will provide an authentication URL"
    echo -e "2. Open the URL in your browser"
    echo -e "3. Log in with your Hytale account"
    echo -e "4. Complete the authentication process"
    echo ""
    echo -e "${YELLOW}━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━${NC}"
    echo ""
    
    # Run the downloader from the GAME_DIR so files download to the right place
    cd "$GAME_DIR"
    if ! "$HYTALE_DOWNLOADER"; then
        log_error "Failed to download/update Hytale server files"
        echo ""
        echo -e "${YELLOW}Troubleshooting Tips:${NC}"
        echo -e "  • Ensure you have a valid Hytale account"
        echo -e "  • Check your internet connection"
        echo -e "  • Mount /etc/machine-id for persistent authentication"
        echo -e "    Example: -v /etc/machine-id:/etc/machine-id:ro"
        echo -e "  • To disable auto-updates, set HYTALE_AUTO_UPDATE=FALSE"
        exit 1
    fi
    
    # Extract the downloaded ZIP file (only if a new one was downloaded)
    log_step "Extracting server files"
    DOWNLOADED_ZIP=$(ls -t "$GAME_DIR"/*.zip 2>/dev/null | head -1)
    if [ -n "$DOWNLOADED_ZIP" ] && [ -f "$DOWNLOADED_ZIP" ]; then
        if unzip -q -o "$DOWNLOADED_ZIP" -d "$GAME_DIR"; then
            log_success
            log_info "Extracted: $(basename "$DOWNLOADED_ZIP")"
            
            # Check if this was an update
            if [ -f "$SERVER_JAR_PATH" ]; then
                log_info "Server files updated to latest version"
            fi
        else
            log_error "Failed to extract server files"
            exit 1
        fi
    elif [ ! -f "$SERVER_JAR_PATH" ]; then
        log_error "No ZIP file found to extract"
        exit 1
    else
        log_success
        log_info "Server files are already up to date"
    fi
fi

# ============================================
# Step 4: Configuration Management
# ============================================
log_section "Server Configuration"

CONFIG_JSON="$BASE_DIR/config.json"
CONFIG_TEMPLATE="/usr/local/share/game-templates/serverconfig.template"

# Check if config.json already exists
if [ ! -f "$CONFIG_JSON" ]; then
    log_step "Rendering config.json with gomplate"
    
    # Render the config.json from template using datasource
    gomplate -f "$CONFIG_TEMPLATE" \
            -d config="$CONFIG_VALUES_FILE" \
            -o "$CONFIG_JSON"
    
    if [ $? -eq 0 ]; then
        log_success
        echo "  config.json rendered successfully"
    else
        log_error "Failed to render config.json"
        exit 1
    fi
else
    log_step "config.json"
    echo -e "${GREEN}already exists${NC}"
fi

# ============================================
# Step 5: Build Server Options
# ============================================
log_section "Building Server Options"

# Initialize all option variables
OPTS=""

# Add JVM options (must come before -jar)
if [ "${HYTALE_CACHE}" = "TRUE" ]; then
    JAVA_ARGS="$JAVA_ARGS -XX:AOTCache=$HYTALE_CACHE_DIR"
    log_info "AOT Cache: ${GREEN}enabled${NC}"
fi

# Add Hytale server options (come after -jar)
[ "${HYTALE_ACCEPT_EARLY_PLUGINS}" = "TRUE" ] && OPTS="$OPTS --accept-early-plugins"
[ "${HYTALE_ALLOW_OP}" = "TRUE" ] && OPTS="$OPTS --allow-op"
[ -n "${HYTALE_AUTH_MODE}" ] && OPTS="$OPTS --auth-mode=$HYTALE_AUTH_MODE"
[ "${HYTALE_BACKUP}" = "TRUE" ] && OPTS="$OPTS --backup"
[ -n "${HYTALE_BACKUP_DIR}" ] && OPTS="$OPTS --backup-dir=$HYTALE_BACKUP_DIR"
[ -n "${HYTALE_BACKUP_FREQUENCY}" ] && OPTS="$OPTS --backup-frequency=$HYTALE_BACKUP_FREQUENCY"
[ -n "${HYTALE_BACKUP_MAX_COUNT}" ] && OPTS="$OPTS --backup-max-count=$HYTALE_BACKUP_MAX_COUNT"
[ "${HYTALE_BARE}" = "TRUE" ] && OPTS="$OPTS --bare"
[ -n "${HYTALE_BOOT_COMMAND}" ] && OPTS="$OPTS --boot-command=$HYTALE_BOOT_COMMAND"
[ -n "${HYTALE_CLIENT_PID}" ] && OPTS="$OPTS --client-pid=$HYTALE_CLIENT_PID"
[ "${HYTALE_DISABLE_ASSET_COMPARE}" = "TRUE" ] && OPTS="$OPTS --disable-asset-compare"
[ "${HYTALE_DISABLE_CPB_BUILD}" = "TRUE" ] && OPTS="$OPTS --disable-cpb-build"
[ "${HYTALE_DISABLE_FILE_WATCHER}" = "TRUE" ] && OPTS="$OPTS --disable-file-watcher"
[ "${HYTALE_DISABLE_SENTRY}" = "TRUE" ] && OPTS="$OPTS --disable-sentry"
[ -n "${HYTALE_EARLY_PLUGINS}" ] && OPTS="$OPTS --early-plugins=$HYTALE_EARLY_PLUGINS"
[ "${HYTALE_EVENT_DEBUG}" = "TRUE" ] && OPTS="$OPTS --event-debug"
[ -n "${HYTALE_FORCE_NETWORK_FLUSH}" ] && OPTS="$OPTS --force-network-flush=$HYTALE_FORCE_NETWORK_FLUSH"
[ "${HYTALE_GENERATE_SCHEMA}" = "TRUE" ] && OPTS="$OPTS --generate-schema"
[ -n "${HYTALE_IDENTITY_TOKEN}" ] && OPTS="$OPTS --identity-token=$HYTALE_IDENTITY_TOKEN"
[ -n "${HYTALE_LOG}" ] && OPTS="$OPTS --log=$HYTALE_LOG"
[ -n "${HYTALE_MIGRATE_WORLDS}" ] && OPTS="$OPTS --migrate-worlds=$HYTALE_MIGRATE_WORLDS"
[ -n "${HYTALE_MIGRATIONS}" ] && OPTS="$OPTS --migrations=$HYTALE_MIGRATIONS"
[ -n "${HYTALE_MODS}" ] && OPTS="$OPTS --mods=$HYTALE_MODS"
[ -n "${HYTALE_OWNER_NAME}" ] && OPTS="$OPTS --owner-name=$HYTALE_OWNER_NAME"
[ -n "${HYTALE_OWNER_UUID}" ] && OPTS="$OPTS --owner-uuid=$HYTALE_OWNER_UUID"
[ -n "${HYTALE_PREFAB_CACHE}" ] && OPTS="$OPTS --prefab-cache=$HYTALE_PREFAB_CACHE"
[ -n "${HYTALE_SESSION_TOKEN}" ] && OPTS="$OPTS --session-token=$HYTALE_SESSION_TOKEN"
[ "${HYTALE_SHUTDOWN_AFTER_VALIDATE}" = "TRUE" ] && OPTS="$OPTS --shutdown-after-validate"
[ "${HYTALE_SINGLEPLAYER}" = "TRUE" ] && OPTS="$OPTS --singleplayer"
[ -n "${HYTALE_TRANSPORT}" ] && OPTS="$OPTS --transport=$HYTALE_TRANSPORT"
[ -n "${HYTALE_UNIVERSE}" ] && OPTS="$OPTS --universe=$HYTALE_UNIVERSE"
[ "${HYTALE_VALIDATE_ASSETS}" = "TRUE" ] && OPTS="$OPTS --validate-assets"
[ -n "${HYTALE_VALIDATE_PREFABS}" ] && OPTS="$OPTS --validate-prefabs=$HYTALE_VALIDATE_PREFABS"
[ "${HYTALE_VALIDATE_WORLD_GEN}" = "TRUE" ] && OPTS="$OPTS --validate-world-gen"
[ "${HYTALE_VERSION}" = "TRUE" ] && OPTS="$OPTS --version"
[ -n "${HYTALE_WORLD_GEN}" ] && OPTS="$OPTS --world-gen=$HYTALE_WORLD_GEN"

log_step "Server options"
echo -e "${GREEN}configured${NC}"

# ============================================
# Step 6: Launch Server
# ============================================
log_section "Launching Hytale Server"

echo ""
echo -e "${BOLD}${GREEN}🚀 Starting Hytale Server...${NC}"
echo ""

# Change to the server directory
cd "$BASE_DIR"

# Build the Java command
JAVA_CMD="java $JAVA_ARGS -Duser.timezone=\"$TZ\" -Dterminal.jline=false -Dterminal.ansi=true -jar \"$SERVER_JAR_PATH\" $OPTS --assets \"$ASSETS_ZIP\" --bind \"$SERVER_IP:$SERVER_PORT\""

# Execute the server
eval exec $JAVA_CMD