package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version   string
	gitCommit string
	buildDate string
)

var rootCmd = &cobra.Command{
	Use:   "gamekeeper",
	Short: "GameKeeper - Unified game server management for Kubernetes",
	Long: `GameKeeper manages game server lifecycle in containerized environments.
	
It handles installation, updates, mod management, configuration, and server startup
for multiple game types including Hytale, Conan Exiles, Seven Days to Die, and more.`,
	SilenceUsage: true,
}

// Execute runs the root command
func Execute() error {
	return rootCmd.Execute()
}

// SetVersion sets version information
func SetVersion(v, commit, date string) {
	version = v
	gitCommit = commit
	buildDate = date
}

func init() {
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(updateCmd)
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(modsCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("GameKeeper %s\n", version)
		fmt.Printf("Git Commit: %s\n", gitCommit)
		fmt.Printf("Build Date: %s\n", buildDate)
	},
}
