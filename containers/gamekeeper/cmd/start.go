package cmd

import (
	"fmt"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
	"github.com/kubelize/game-servers/gamekeeper/pkg/output"
	"github.com/kubelize/game-servers/gamekeeper/pkg/server"
	"github.com/spf13/cobra"
)

var (
	gameType      string
	configPath    string
	skipUpdate    bool
	forceUpdate   bool
)

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a game server",
	Long: `Start a game server with full lifecycle management:
	
- Download/install game files if missing
- Check for updates (unless --skip-update)
- Install and configure mods
- Render configuration files
- Start the game server process`,
	RunE: runStart,
}

func init() {
	startCmd.Flags().StringVar(&gameType, "game", "", "Game type (hytale, conan-exiles, sdtd, palworld, valheim, minecraft)")
	startCmd.Flags().StringVar(&configPath, "config", "/config-data/config-values.yaml", "Path to config values file")
	startCmd.Flags().BoolVar(&skipUpdate, "skip-update", false, "Skip checking for game updates")
	startCmd.Flags().BoolVar(&forceUpdate, "force-update", false, "Force update even if up to date")
	startCmd.MarkFlagRequired("game")
}

func runStart(cmd *cobra.Command, args []string) error {
	fmt.Printf("🎮 %sStarting GameKeeper for %s%s\n", output.Bold, gameType, output.Reset)

	// Load configuration
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Create server manager for the game type
	mgr, err := server.NewManager(gameType, cfg)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}

	// Setup phase
	output.Section("Setting up directories")
	output.Step("Creating directories")
	if err := mgr.Setup(); err != nil {
		output.Error(err.Error())
		return fmt.Errorf("setup failed: %w", err)
	}
	output.Success()

	// Download/update phase
	if !skipUpdate {
		autoUpdate := cfg.GetBool("HYTALE_AUTO_UPDATE", true)
		if autoUpdate || forceUpdate {
			output.Section("Checking for game updates")
			if err := mgr.Update(forceUpdate); err != nil {
				output.Error(err.Error())
				return fmt.Errorf("update failed: %w", err)
			}
		}
	}

	// Mod installation phase
	output.Section("Installing mods")
	if err := mgr.InstallMods(); err != nil {
		output.Error(err.Error())
		return fmt.Errorf("mod installation failed: %w", err)
	}

	// Configuration phase
	output.Section("Rendering configuration")
	if err := mgr.Configure(); err != nil {
		output.Error(err.Error())
		return fmt.Errorf("configuration failed: %w", err)
	}

	// Validation phase
	output.Section("Validating setup")
	if err := mgr.Validate(); err != nil {
		output.Error(err.Error())
		return fmt.Errorf("validation failed: %w", err)
	}

	// Start server
	output.Section("Launching Game Server")
	output.Launch("Starting game server...")
	if err := mgr.Start(); err != nil {
		return fmt.Errorf("server start failed: %w", err)
	}

	// This blocks until server exits
	return nil
}
