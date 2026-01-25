package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
	"github.com/kubelize/game-servers/gamekeeper/pkg/rcon"
)

// SteamManager manages Steam-based game servers
type SteamManager struct {
	*BaseManager
	appID int
}

// NewSteamManager creates a new Steam game server manager
func NewSteamManager(cfg *config.Config, gameType string, appID int) *SteamManager {
	baseDir := cfg.GetString("BASE_DIR", "/home/kubelize/server")
	
	return &SteamManager{
		BaseManager: &BaseManager{
			GameType: gameType,
			Config:   cfg,
			BaseDir:  baseDir,
			DataDir:  baseDir,
		},
		appID: appID,
	}
}

func (s *SteamManager) Setup() error {
	return s.ensureDirectories(s.BaseDir)
}

func (s *SteamManager) Update(force bool) error {
	fmt.Printf("  → Installing/updating %s (Steam AppID: %d)...\n", s.GameType, s.appID)
	
	steamCmd := "/home/kubelize/steam/steamcmd.sh"
	args := []string{
		"+force_install_dir", s.BaseDir,
		"+login", "anonymous",
		"+app_update", fmt.Sprintf("%d", s.appID),
		"validate",
		"+quit",
	}

	cmd := exec.Command(steamCmd, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	return cmd.Run()
}

func (s *SteamManager) CheckUpdate() (bool, string, error) {
	// SteamCMD handles updates automatically
	return false, "", nil
}

func (s *SteamManager) InstallMods() error {
	// TODO: Implement Steam Workshop mod installation
	fmt.Println("  ℹ  Steam mod management not yet implemented")
	return nil
}

func (s *SteamManager) Configure() error {
	// Game-specific configuration handled by subclasses or config files
	return nil
}

func (s *SteamManager) Validate() error {
	// Check if server files exist
	if !fileExists(s.BaseDir) {
		return fmt.Errorf("server directory does not exist: %s", s.BaseDir)
	}
	return nil
}

func (s *SteamManager) Start() error {
	consolePort := s.Config.GetString("CONSOLE_PORT", "8080")
	sessionName := s.Config.GetString("TMUX_SESSION_NAME", fmt.Sprintf("%s-server", s.GameType))
	
	// Get the start command and args from config
	startCommand := s.Config.GetString("START_COMMAND", filepath.Join(s.BaseDir, "startserver.sh"))
	startArgs := s.Config.GetString("START_ARGS", "")
	
	var args []string
	if startArgs != "" {
		args = splitArgs(startArgs)
	}
	
	return rcon.StartServerWithTmux(consolePort, sessionName, startCommand, args, s.BaseDir)
}

func (s *SteamManager) Stop() error {
	return nil
}
