package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
	"github.com/kubelize/game-servers/gamekeeper/pkg/output"
	"github.com/kubelize/game-servers/gamekeeper/pkg/rcon"
)

// HytaleManager manages Hytale game servers
type HytaleManager struct {
	*BaseManager
	downloaderPath string
	serverJarPath  string
	assetsZipPath  string
}

// NewHytaleManager creates a new Hytale server manager
func NewHytaleManager(cfg *config.Config) *HytaleManager {
	baseDir := cfg.GetString("BASE_DIR", "/home/kubelize/server")
	dataDir := filepath.Join(baseDir, "data")

	return &HytaleManager{
		BaseManager: &BaseManager{
			GameType: "hytale",
			Config:   cfg,
			BaseDir:  baseDir,
			DataDir:  dataDir,
		},
		downloaderPath: filepath.Join(dataDir, "hytale-downloader"),
		serverJarPath:  filepath.Join(dataDir, "Server", "HytaleServer.jar"),
		assetsZipPath:  filepath.Join(dataDir, "Assets.zip"),
	}
}

func (h *HytaleManager) Setup() error {
	return h.ensureDirectories(
		h.BaseDir,
		h.DataDir,
		filepath.Join(h.BaseDir, "config-data"),
	)
}

func (h *HytaleManager) Update(force bool) error {
	// Check if files exist and auto-update is enabled
	filesExist := fileExists(h.serverJarPath) && fileExists(h.assetsZipPath)
	autoUpdate := h.Config.GetBool("HYTALE_AUTO_UPDATE", true)

	shouldDownload := !filesExist || (autoUpdate && force)

	if !shouldDownload && filesExist {
		output.Step("Server files")
		output.SuccessWithMessage("already downloaded (auto-update disabled)")
		return nil
	}

	// Ensure downloader exists
	if !fileExists(h.downloaderPath) {
		output.Step("Downloading hytale-downloader")
		url := h.Config.GetString("HYTALE_DOWNLOADER_URL", 
			"https://drive.kubelize.com/public.php/dav/files/HJqqWZx5522wnoT")
		
		if err := downloadFile(url, h.downloaderPath); err != nil {
			output.Error(err.Error())
			return fmt.Errorf("failed to download hytale-downloader: %w", err)
		}
		
		if err := os.Chmod(h.downloaderPath, 0755); err != nil {
			output.Error(err.Error())
			return fmt.Errorf("failed to make downloader executable: %w", err)
		}
		output.Success()
	}

	// Run the downloader
	output.Step("Running hytale-downloader")
	fmt.Println()
	output.Warning("Authentication may be required - follow prompts")
	fmt.Println()
	
	cmd := exec.Command(h.downloaderPath)
	cmd.Dir = h.DataDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	if err := cmd.Run(); err != nil {
		output.Error(err.Error())
		return fmt.Errorf("hytale-downloader failed: %w", err)
	}

	// Extract downloaded ZIP
	output.Step("Extracting server files")
	return h.extractLatestZip()
}

func (h *HytaleManager) CheckUpdate() (bool, string, error) {
	// TODO: Implement version checking logic
	return false, "", nil
}

func (h *HytaleManager) InstallMods() error {
	if !h.Config.Mods.Enabled || len(h.Config.Mods.Mods) == 0 {
		output.Info("No mods configured")
		return nil
	}

	modsDir := filepath.Join(h.DataDir, "Mods")
	if err := ensureDir(modsDir); err != nil {
		return err
	}

	for _, mod := range h.Config.Mods.Mods {
		fmt.Printf("  → Installing mod: %s (%s)\n", mod.Name, mod.Version)
		if err := h.installMod(mod); err != nil {
			return fmt.Errorf("failed to install mod %s: %w", mod.Name, err)
		}
	}

	return nil
}

func (h *HytaleManager) installMod(mod config.ModConfig) error {
	// TODO: Implement full mod installation
	// 1. Download from URL
	// 2. Verify checksum
	// 3. Extract to installPath
	// 4. Place config files
	return nil
}

func (h *HytaleManager) Configure() error {
	configPath := filepath.Join(h.BaseDir, "config.json")
	
	// Skip if already exists (unless force)
	if fileExists(configPath) {
		output.Step("config.json")
		output.SuccessWithMessage("already exists")
		return nil
	}

	output.Step("Rendering config.json")
	// TODO: Implement gomplate rendering or native Go templating
	return nil
}

func (h *HytaleManager) Validate() error {
	// Check required files exist
	required := []string{h.serverJarPath, h.assetsZipPath}
	for _, path := range required {
		if !fileExists(path) {
			return fmt.Errorf("required file missing: %s", path)
		}
	}
	return nil
}

func (h *HytaleManager) Start() error {
	// Build Java command
	javaArgs := h.Config.GetString("JAVA_ARGS", "")
	serverPort := h.Config.GetString("SERVER_PORT", "5520")
	serverIP := h.Config.GetString("SERVER_IP", "0.0.0.0")
	tz := h.Config.GetString("TZ", "UTC")

	// Build Hytale-specific options
	opts := h.buildServerOptions()

	args := []string{}
	if javaArgs != "" {
		// Parse Java args
		args = append(args, splitArgs(javaArgs)...)
	}
	args = append(args, 
		fmt.Sprintf("-Duser.timezone=%s", tz),
		"-Dterminal.jline=false",
		"-Dterminal.ansi=true",
		"-jar", h.serverJarPath,
	)
	args = append(args, splitArgs(opts)...)
	args = append(args,
		"--assets", h.assetsZipPath,
		"--bind", fmt.Sprintf("%s:%s", serverIP, serverPort),
	)

	// Create console pipe and start gotty web console
	consolePipe := "/tmp/gameserver-console.pipe"
	consolePort := h.Config.GetString("CONSOLE_PORT", "8080")
	
	if err := rcon.CreateConsolePipe(consolePipe); err != nil {
		return fmt.Errorf("failed to create console pipe: %w", err)
	}
	output.Info(fmt.Sprintf("Console pipe created at: %s", consolePipe))

	// Start gotty console in background
	if err := rcon.StartGottyConsole(consolePort, consolePipe); err != nil {
		output.Warning(fmt.Sprintf("Failed to start web console: %v", err))
	} else {
		output.Info(fmt.Sprintf("Web console available on port %s", consolePort))
		output.Info(fmt.Sprintf("Access via: kubectl port-forward <pod> %s:%s", consolePort, consolePort))
	}

	// Open the console pipe for stdin
	pipe, err := rcon.OpenConsolePipe(consolePipe)
	if err != nil {
		return fmt.Errorf("failed to open console pipe: %w", err)
	}

	cmd := exec.Command("java", args...)
	cmd.Dir = h.BaseDir
	cmd.Stdin = pipe       // Read from named pipe
	cmd.Stdout = os.Stdout // Normal logging - kubectl logs works!
	cmd.Stderr = os.Stderr // Normal logging - kubectl logs works!

	return cmd.Run()
}

func (h *HytaleManager) Stop() error {
	// TODO: Implement graceful shutdown
	return nil
}

func (h *HytaleManager) buildServerOptions() string {
	opts := ""
	
	// Map of config keys to command-line flags
	boolFlags := map[string]string{
		"HYTALE_ACCEPT_EARLY_PLUGINS":     "--accept-early-plugins",
		"HYTALE_ALLOW_OP":                 "--allow-op",
		"HYTALE_BACKUP":                   "--backup",
		"HYTALE_BARE":                     "--bare",
		"HYTALE_DISABLE_ASSET_COMPARE":    "--disable-asset-compare",
		"HYTALE_DISABLE_CPB_BUILD":        "--disable-cpb-build",
		"HYTALE_DISABLE_FILE_WATCHER":     "--disable-file-watcher",
		"HYTALE_DISABLE_SENTRY":           "--disable-sentry",
		"HYTALE_EVENT_DEBUG":              "--event-debug",
		"HYTALE_GENERATE_SCHEMA":          "--generate-schema",
		"HYTALE_SHUTDOWN_AFTER_VALIDATE":  "--shutdown-after-validate",
		"HYTALE_SINGLEPLAYER":             "--singleplayer",
		"HYTALE_VALIDATE_ASSETS":          "--validate-assets",
		"HYTALE_VALIDATE_WORLD_GEN":       "--validate-world-gen",
	}

	valueFlags := map[string]string{
		"HYTALE_AUTH_MODE":          "--auth-mode",
		"HYTALE_BACKUP_DIR":         "--backup-dir",
		"HYTALE_BACKUP_FREQUENCY":   "--backup-frequency",
		"HYTALE_BACKUP_MAX_COUNT":   "--backup-max-count",
		"HYTALE_BOOT_COMMAND":       "--boot-command",
		"HYTALE_IDENTITY_TOKEN":     "--identity-token",
		"HYTALE_SESSION_TOKEN":      "--session-token",
		"HYTALE_OWNER_NAME":         "--owner-name",
		"HYTALE_OWNER_UUID":         "--owner-uuid",
		"HYTALE_TRANSPORT":          "--transport",
		"HYTALE_UNIVERSE":           "--universe",
		"HYTALE_WORLD_GEN":          "--world-gen",
	}

	// Add boolean flags
	for key, flag := range boolFlags {
		if h.Config.GetBool(key, false) {
			opts += " " + flag
		}
	}

	// Add value flags
	for key, flag := range valueFlags {
		if val := h.Config.GetString(key, ""); val != "" {
			opts += fmt.Sprintf(" %s=%s", flag, val)
		}
	}

	return opts
}

func (h *HytaleManager) extractLatestZip() error {
	// Find the latest ZIP file
	pattern := filepath.Join(h.DataDir, "*.zip")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return err
	}
	
	if len(matches) == 0 {
		return fmt.Errorf("no ZIP file found to extract")
	}

	// Use the first (newest) match
	zipFile := matches[0]
	
	cmd := exec.Command("unzip", "-q", "-o", zipFile, "-d", h.DataDir)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to extract %s: %w", zipFile, err)
	}

	fmt.Printf("  ✓ Extracted: %s\n", filepath.Base(zipFile))
	return nil
}
