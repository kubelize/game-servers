package server

import (
	"fmt"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
)

// Manager defines the interface for game server management
type Manager interface {
	// Setup prepares directories and initial environment
	Setup() error

	// Update checks for and applies game updates
	Update(force bool) error

	// CheckUpdate checks if an update is available
	CheckUpdate() (hasUpdate bool, version string, err error)

	// InstallMods downloads and installs mods from configuration
	InstallMods() error

	// Configure renders and applies configuration files
	Configure() error

	// Validate checks if the server is properly configured
	Validate() error

	// Start launches the game server (blocking)
	Start() error

	// Stop gracefully shuts down the game server
	Stop() error
}

// BaseManager provides common functionality for all game servers
type BaseManager struct {
	GameType string
	Config   *config.Config
	BaseDir  string
	DataDir  string
}

// NewManager creates a server manager for the specified game type
func NewManager(gameType string, cfg *config.Config) (Manager, error) {
	switch gameType {
	case "hytale":
		return NewHytaleManager(cfg), nil
	case "conan-exiles", "ce":
		return NewSteamManager(cfg, "conan-exiles", 443030), nil
	case "seven-days-to-die", "sdtd":
		return NewSteamManager(cfg, "sdtd", 294420), nil
	case "palworld":
		return NewSteamManager(cfg, "palworld", 2394010), nil
	case "valheim":
		return NewSteamManager(cfg, "valheim", 896660), nil
	case "minecraft":
		return NewMinecraftManager(cfg), nil
	default:
		return nil, fmt.Errorf("unsupported game type: %s", gameType)
	}
}

// Common helper methods
func (b *BaseManager) ensureDirectories(dirs ...string) error {
	for _, dir := range dirs {
		if err := ensureDir(dir); err != nil {
			return err
		}
	}
	return nil
}
