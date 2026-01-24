# GameKeeper

**GameKeeper** is a unified game server management tool for Kubernetes-based game server deployments.

## Purpose

GameKeeper handles the lifecycle of game servers within containerized environments:

- **Server Installation**: Download and install game server files
- **Update Management**: Check for and apply game updates automatically
- **Mod Management**: Install, configure, and manage game mods from ConfigMaps
- **Configuration**: Render and place server configs in the correct locations
- **Lifecycle Control**: Start, monitor, and gracefully shutdown game servers

## Supported Games

- **Hytale** - Multiplayer sandbox RPG

## In Progress

- **Conan Exiles** - Survival game (Steam)
- **Seven Days to Die** - Survival horror (Steam)
- **Palworld** - Creature collection survival (Steam)
- **Valheim** - Viking survival (Steam)
- **Minecraft** - Block-building game

## Architecture

GameKeeper is designed to work with Helm charts:
- **Helm Chart**: Declares *what* should be installed (mods, configs, settings)
- **GameKeeper**: Implements *how* to install it (download, place, configure)

```
Kubernetes Pod
├── ConfigMap (mounted) → mods.yaml, server-config.yaml
├── PVC (mounted) → /home/kubelize/server/data
└── GameKeeper binary
    ├── Reads ConfigMap
    ├── Downloads game files
    ├── Installs mods
    ├── Renders configs
    └── Starts game server
```

## Usage

```bash
# Start a game server (reads config from mounted ConfigMaps)
gamekeeper start --game hytale

# Check for updates without starting
gamekeeper update --game hytale --check-only

# Validate configuration
gamekeeper validate --game hytale

# List installed mods
gamekeeper mods list

# Install a specific mod
gamekeeper mods install --name MyMod --version 1.2.3
```

## Configuration

GameKeeper reads configuration from:
- `/config-data/config-values.yaml` - Server settings (from ConfigMap)
- `/config-data/mods.yaml` - Mod specifications (from ConfigMap)
- Environment variables - Runtime overrides

Example mod configuration:
```yaml
mods:
  - name: "ExampleMod"
    version: "1.0.0"
    url: "https://example.com/mod.zip"
    checksum: "sha256:abc123..."
    installPath: "Mods/"
    loadOrder: 10
    configFiles:
      - source: "config.json"
        destination: "Config/ExampleMod/config.json"
        template: true
```

## Building

```bash
cd containers/gamekeeper
go build -o gamekeeper
```

## Development

See [docs/architecture.md](docs/architecture.md) for detailed design documentation.

## License

See LICENSE in the root of the repository.
