package main

import (
	"fmt"
	"os"

	"github.com/kubelize/game-servers/gamekeeper/cmd"
)

// Version information
var (
	Version   = "dev"
	GitCommit = "unknown"
	BuildDate = "unknown"
)

func main() {
	cmd.SetVersion(Version, GitCommit, BuildDate)
	
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
