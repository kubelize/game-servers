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
	// Game-specific configuration
	// TODO: Implement per-game config rendering
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
	// Game-specific startup logic
	switch s.GameType {
	case "sdtd":
		return s.startSDTD()
	default:
		return fmt.Errorf("start not implemented for %s", s.GameType)
	}
}

func (s *SteamManager) Stop() error {
	return nil
}

func (s *SteamManager) startSDTD() error {
	// Create console pipe and start gotty web console
	consolePipe := "/tmp/gameserver-console.pipe"
	consolePort := s.Config.GetString("CONSOLE_PORT", "8080")
	
	if err := rcon.CreateConsolePipe(consolePipe); err != nil {
		return fmt.Errorf("failed to create console pipe: %w", err)
	}
	fmt.Printf("  ℹ Console pipe created at: %s\n", consolePipe)

	// Start gotty console in background
	if err := rcon.StartGottyConsole(consolePort, consolePipe); err != nil {
		fmt.Printf("  ⚠ Warning: Failed to start web console: %v\n", err)
	}

	configPath := filepath.Join(s.BaseDir, "sdtdconfig.xml")
	startScript := filepath.Join(s.BaseDir, "startserver.sh")
	
	// Open the console pipe for stdin
	pipe, err := rcon.OpenConsolePipe(consolePipe)
	if err != nil {
		return fmt.Errorf("failed to open console pipe: %w", err)
	}

	cmd := exec.Command(startScript, fmt.Sprintf("-configfile=%s", configPath))
	cmd.Dir = s.BaseDir
	cmd.Stdin = pipe       // Read from named pipe
	cmd.Stdout = os.Stdout // Normal logging - kubectl logs works!
	cmd.Stderr = os.Stderr // Normal logging - kubectl logs works!
	
	return cmd.Run()
}

// MinecraftManager manages Minecraft servers
type MinecraftManager struct {
	*BaseManager
}

func NewMinecraftManager(cfg *config.Config) *MinecraftManager {
	baseDir := cfg.GetString("BASE_DIR", "/home/kubelize/server")
	
	return &MinecraftManager{
		BaseManager: &BaseManager{
			GameType: "minecraft",
			Config:   cfg,
			BaseDir:  baseDir,
			DataDir:  baseDir,
		},
	}
}

func (m *MinecraftManager) Setup() error {
	return m.ensureDirectories(m.BaseDir)
}

func (m *MinecraftManager) Update(force bool) error {
	// TODO: Implement Minecraft server download
	fmt.Println("  ℹ  Minecraft update not yet implemented")
	return nil
}

func (m *MinecraftManager) CheckUpdate() (bool, string, error) {
	return false, "", nil
}

func (m *MinecraftManager) InstallMods() error {
	fmt.Println("  ℹ  Minecraft mod management not yet implemented")
	return nil
}

func (m *MinecraftManager) Configure() error {
	return nil
}

func (m *MinecraftManager) Validate() error {
	return nil
}

func (m *MinecraftManager) Start() error {
	return fmt.Errorf("minecraft start not yet implemented")
}

func (m *MinecraftManager) Stop() error {
	return nil
}
