package rcon

import (
	"fmt"
	"os"
	"os/exec"
)

// StartGottyConsole starts a gotty web console that writes to the game server's stdin pipe
// This allows web-based interaction while keeping kubectl logs working
func StartGottyConsole(port, pipePath string) error {
	// Create a script that forwards input to the game server pipe
	script := fmt.Sprintf(`#!/bin/bash
echo "GameKeeper Web Console - Connected to game server"
echo "Type commands and press Enter to send to server"
echo ""

# Keep pipe open for writing
exec 3> %s

while IFS= read -r line; do
  echo "$line" >&3
  echo "[Sent: $line]" >&2
done

# Close pipe on exit
exec 3>&-
`, pipePath)

	// Create temp script file
	scriptPath := "/tmp/console-input.sh"
	if err := os.WriteFile(scriptPath, []byte(script), 0755); err != nil {
		return fmt.Errorf("failed to create console script: %w", err)
	}

	// Start gotty wrapping the script
	// -w: allow writes (input)
	// -p: port
	cmd := exec.Command("gotty", "-w", "-p", port, "/bin/bash", scriptPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start gotty: %w", err)
	}

	fmt.Printf("  ℹ Web console available on port %s\n", port)
	fmt.Printf("  ℹ Access via: kubectl port-forward <pod> %s:%s\n", port, port)
	
	return nil
}

// CreateConsolePipe creates a named pipe for console input
func CreateConsolePipe(pipePath string) error {
	// Remove existing pipe if it exists
	os.Remove(pipePath)
	
	// Create named pipe using mkfifo
	cmd := exec.Command("mkfifo", pipePath)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create pipe: %w", err)
	}
	
	// Make it world-writable so gotty can write to it
	if err := os.Chmod(pipePath, 0666); err != nil {
		return fmt.Errorf("failed to chmod pipe: %w", err)
	}
	
	return nil
}

// OpenConsolePipe opens the named pipe for reading (as stdin)
func OpenConsolePipe(pipePath string) (*os.File, error) {
	pipe, err := os.OpenFile(pipePath, os.O_RDONLY, os.ModeNamedPipe)
	if err != nil {
		return nil, fmt.Errorf("failed to open console pipe: %w", err)
	}
	return pipe, nil
}
