package cmd

import (
	"fmt"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
	"github.com/kubelize/game-servers/gamekeeper/pkg/server"
	"github.com/spf13/cobra"
)

var (
	checkOnly bool
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Check for and apply game server updates",
	RunE:  runUpdate,
}

func init() {
	updateCmd.Flags().StringVar(&gameType, "game", "", "Game type")
	updateCmd.Flags().StringVar(&configPath, "config", "/home/kubelize/config-data/config-values.yaml", "Path to config values file")
	updateCmd.Flags().BoolVar(&checkOnly, "check-only", false, "Only check for updates, don't apply")
	updateCmd.MarkFlagRequired("game")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mgr, err := server.NewManager(gameType, cfg)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}

	if checkOnly {
		hasUpdate, version, err := mgr.CheckUpdate()
		if err != nil {
			return fmt.Errorf("update check failed: %w", err)
		}
		if hasUpdate {
			fmt.Printf("✨ Update available: %s\n", version)
		} else {
			fmt.Println("✅ Already up to date")
		}
		return nil
	}

	return mgr.Update(true)
}

var validateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate game server configuration",
	RunE:  runValidate,
}

func init() {
	validateCmd.Flags().StringVar(&gameType, "game", "", "Game type")
	validateCmd.Flags().StringVar(&configPath, "config", "/home/kubelize/config-data/config-values.yaml", "Path to config values file")
	validateCmd.MarkFlagRequired("game")
}

func runValidate(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	mgr, err := server.NewManager(gameType, cfg)
	if err != nil {
		return fmt.Errorf("failed to create server manager: %w", err)
	}

	if err := mgr.Validate(); err != nil {
		fmt.Println("❌ Validation failed:")
		return err
	}

	fmt.Println("✅ Configuration is valid")
	return nil
}
