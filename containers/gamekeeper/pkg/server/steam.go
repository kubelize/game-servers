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
	args := []string{}

	platformType := s.Config.GetString("STEAM_PLATFORM", "")
	if platformType == "" && s.GameType == "conan-exiles" {
		// Conan dedicated server ships as Windows binaries and needs Wine.
		platformType = "windows"
	}
	if platformType != "" {
		args = append(args, "+@sSteamCmdForcePlatformType", platformType)
	}

	args = append(args,
		"+force_install_dir", s.BaseDir,
		"+login", "anonymous",
		"+app_update", fmt.Sprintf("%d", s.appID),
		"validate",
		"+quit",
	)

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
	templatePath := s.Config.GetString("CONFIG_TEMPLATE_PATH", defaultTemplatePath)
	outputPath := s.Config.GetString("CONFIG_OUTPUT_PATH", s.defaultConfigOutputPath())
	configPath := s.Config.GetString("CONFIG_VALUES_PATH", defaultConfigPath)
	passwordPath := s.Config.GetString("CONFIG_PASSWORD_PATH", "")
	if passwordPath == "" {
		passwordPath = s.Config.GetString("PASSWORD_VALUES_PATH", defaultPasswordPath)
	}

	if outputPath == "" {
		fmt.Printf("  ℹ  No config template output configured for %s\n", s.GameType)
		return nil
	}

	if !fileExists(templatePath) {
		fmt.Printf("  ℹ  Template not found at %s, skipping render\n", templatePath)
		return nil
	}

	fmt.Printf("  → Rendering config template to %s...\n", outputPath)
	if err := renderTemplateWithGomplate(templateRenderSpec{
		TemplatePath: templatePath,
		OutputPath:   outputPath,
		ConfigPath:   configPath,
		PasswordPath: passwordPath,
	}); err != nil {
		return err
	}

	return nil
}

func (s *SteamManager) defaultConfigOutputPath() string {
	switch s.GameType {
	case "conan-exiles":
		return filepath.Join(s.BaseDir, "ConanSandbox", "Config", "DefaultServerSettings.ini")
	case "sdtd":
		return filepath.Join(s.BaseDir, "sdtdconfig.xml")
	default:
		return ""
	}
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

	// Get sensible game-specific defaults and allow overriding via config.
	defaultStartCommand, defaultStartArgs := s.defaultStartConfig()
	startCommand := s.Config.GetString("START_COMMAND", defaultStartCommand)
	startArgs := s.Config.GetString("START_ARGS", defaultStartArgs)

	var args []string
	if startArgs != "" {
		args = splitArgs(startArgs)
	}

	// Fall back to direct process start when tmux isn't available in the image.
	if _, err := exec.LookPath("tmux"); err != nil {
		fmt.Println("  ℹ  tmux not found, starting server process directly")
		cmd := exec.Command(startCommand, args...)
		cmd.Dir = s.BaseDir
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.Stdin = os.Stdin
		return cmd.Run()
	}

	return rcon.StartServerWithTmux(consolePort, sessionName, startCommand, args, s.BaseDir)
}

func (s *SteamManager) Stop() error {
	return nil
}

func (s *SteamManager) defaultStartConfig() (string, string) {
	switch s.GameType {
	case "conan-exiles":
		exePath := filepath.Join(s.BaseDir, "ConanSandbox", "Binaries", "Win64", "ConanSandboxServer-Win64-Shipping.exe")
		return "xvfb-run", fmt.Sprintf("-a wine %s -log", exePath)
	case "sdtd":
		return filepath.Join(s.BaseDir, "startserver.sh"), fmt.Sprintf("-configfile=%s", filepath.Join(s.BaseDir, "sdtdconfig.xml"))
	default:
		return filepath.Join(s.BaseDir, "startserver.sh"), ""
	}
}
