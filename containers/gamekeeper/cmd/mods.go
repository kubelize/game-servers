package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var modsCmd = &cobra.Command{
	Use:   "mods",
	Short: "Manage game server mods",
	Long:  "Install, list, and manage mods for game servers",
}

var modsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List installed mods",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📦 Installed mods:")
		// TODO: Implement
		return nil
	},
}

var modsInstallCmd = &cobra.Command{
	Use:   "install",
	Short: "Install a mod",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("📥 Installing mod...")
		// TODO: Implement
		return nil
	},
}

func init() {
	modsCmd.AddCommand(modsListCmd)
	modsCmd.AddCommand(modsInstallCmd)
}
