package rcon

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// StartServerWithTmux runs the server in a tmux session with output piped to stdout
// and gotty attached for web console access
func StartServerWithTmux(port, sessionName, command string, args []string, workdir string) error {
	// Build the command string
	fullCmd := command + " " + strings.Join(args, " ")
	
	// Create log file for tmux output
	logFile := filepath.Join(workdir, "console.log")
	
	// Create tmux session with the server
	tmuxCmd := exec.Command("tmux", "new-session", "-d", "-s", sessionName,
		"-c", workdir,  // working directory
		fullCmd,
	)
	
	if err := tmuxCmd.Run(); err != nil {
		return fmt.Errorf("failed to start tmux session: %w", err)
	}
	
	// Pipe tmux output to log file for kubectl logs
	pipeCmd := exec.Command("tmux", "pipe-pane", "-t", sessionName, 
		fmt.Sprintf("cat >> %s", logFile))
	if err := pipeCmd.Run(); err != nil {
		return fmt.Errorf("failed to setup tmux pipe-pane: %w", err)
	}
	
	// Start tailing the log file to stdout in background
	go tailLogToStdout(logFile)
	
	// Start gotty that attaches to the tmux session
	// This provides full bidirectional terminal access via web
	gottyCmd := exec.Command("gotty",
		"-w",              // allow writes (permit-write)
		"-p", port,        // port
		"--reconnect",
		"--title-format", "Game Server Console",
		"tmux", "attach-session", "-t", sessionName,
	)
	gottyCmd.Stdout = os.Stdout
	gottyCmd.Stderr = os.Stderr
	
	if err := gottyCmd.Start(); err != nil {
		return fmt.Errorf("failed to start gotty: %w", err)
	}

	fmt.Printf("  ℹ Web console available on port %s\n", port)
	fmt.Printf("  ℹ Access via: kubectl port-forward <pod> %s:%s\n", port, port)
	fmt.Printf("  ℹ Manual attach: kubectl exec -it <pod> -- tmux attach -t %s\n", sessionName)
	
	// Wait for tmux session to end
	for {
		checkCmd := exec.Command("tmux", "has-session", "-t", sessionName)
		if err := checkCmd.Run(); err != nil {
			// Session ended
			break
		}
		// Check every 5 seconds
		exec.Command("sleep", "5").Run()
	}
	
	return nil
}

// tailLogToStdout continuously reads from log file and prints to stdout
func tailLogToStdout(logFile string) {
	// Wait for file to exist
	for {
		if _, err := os.Stat(logFile); err == nil {
			break
		}
		exec.Command("sleep", "1").Run()
	}
	
	// Open and tail the file
	file, err := os.Open(logFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open log file: %v\n", err)
		return
	}
	defer file.Close()
	
	// Start from beginning to capture all output
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			// No new content, wait a bit
			exec.Command("sleep", "0").Run() // Small delay
			continue
		}
		fmt.Print(line)
	}
}
