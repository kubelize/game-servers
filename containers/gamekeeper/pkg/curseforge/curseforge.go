package curseforge

import (
	"crypto/md5"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/kubelize/game-servers/gamekeeper/pkg/config"
	"github.com/kubelize/game-servers/gamekeeper/pkg/output"
)

const (
	cfAPIBase = "https://api.curseforge.com"
	cfAPIHost = "api.curseforge.com"
)

// ReleaseType represents CurseForge release types
type ReleaseType int

const (
	ReleaseTypeRelease ReleaseType = 1
	ReleaseTypeBeta    ReleaseType = 2
	ReleaseTypeAlpha   ReleaseType = 3
)

// Client handles CurseForge API interactions
type Client struct {
	apiKey       string
	httpClient   *http.Client
	cacheAPIURL  string
	cacheDownURL string
}

// ModFile represents a CurseForge mod file
type ModFile struct {
	ID          int       `json:"id"`
	FileName    string    `json:"fileName"`
	DisplayName string    `json:"displayName"`
	FileDate    string    `json:"fileDate"`
	ReleaseType int       `json:"releaseType"`
	IsAvailable bool      `json:"isAvailable"`
	Hashes      []Hash    `json:"hashes"`
	GameVersions []string `json:"gameVersions"`
}

// Hash represents a file hash
type Hash struct {
	Algo  int    `json:"algo"` // 1 = SHA1, 2 = MD5
	Value string `json:"value"`
}

// Manifest tracks installed mods
type Manifest struct {
	SchemaVersion  int                    `json:"schemaVersion"`
	LastCheckEpoch int64                  `json:"lastCheckEpoch"`
	Mods           map[string]ManifestMod `json:"mods"`
}

// ManifestMod represents a mod in the manifest
type ManifestMod struct {
	Reference string          `json:"reference"`
	Resolved  *ResolvedFile   `json:"resolved,omitempty"`
	Installed *InstalledFile  `json:"installed,omitempty"`
}

// ResolvedFile represents resolved file info
type ResolvedFile struct {
	FileID      int    `json:"fileId"`
	FileName    string `json:"fileName"`
	DisplayName string `json:"displayName"`
	FileDate    string `json:"fileDate"`
	ReleaseType int    `json:"releaseType"`
	DownloadURL string `json:"downloadUrl"`
	Hash        *HashInfo `json:"hash,omitempty"`
}

// HashInfo represents hash information
type HashInfo struct {
	Algo  string `json:"algo"`
	Value string `json:"value"`
}

// InstalledFile represents installed file info
type InstalledFile struct {
	FileID          int    `json:"fileId"`
	Path            string `json:"path"`
	InstalledAtEpoch int64 `json:"installedAtEpoch"`
}

// NewClient creates a new CurseForge client
func NewClient(cfg *config.Config) (*Client, error) {
	apiKey := cfg.GetString("HYTALE_CURSEFORGE_API_KEY", "")
	
	// Check for API key from file
	apiKeySrc := cfg.GetString("HYTALE_CURSEFORGE_API_KEY_SRC", "")
	if apiKeySrc != "" {
		if data, err := os.ReadFile(apiKeySrc); err == nil {
			apiKey = strings.TrimSpace(string(data))
		}
	}

	if apiKey == "" {
		return nil, fmt.Errorf("HYTALE_CURSEFORGE_API_KEY is required")
	}

	// Validate API key format
	if !strings.HasPrefix(apiKey, "$2a$10$") {
		return nil, fmt.Errorf("invalid API key format - should start with '$2a$10$'")
	}

	client := &Client{
		apiKey: apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		cacheAPIURL:  cfg.GetString("HYTALE_CURSEFORGE_HTTP_CACHE_API_URL", ""),
		cacheDownURL: cfg.GetString("HYTALE_CURSEFORGE_HTTP_CACHE_DOWNLOAD_URL", ""),
	}

	// Test API key
	if err := client.testAPIKey(); err != nil {
		return nil, fmt.Errorf("API key validation failed: %w", err)
	}

	return client, nil
}

// testAPIKey validates the API key
func (c *Client) testAPIKey() error {
	_, err := c.apiGet("/v1/games")
	return err
}

// apiGet makes a GET request to the CurseForge API
func (c *Client) apiGet(path string) ([]byte, error) {
	url := cfAPIBase + path
	hostHeader := ""
	
	if c.cacheAPIURL != "" {
		url = strings.TrimSuffix(c.cacheAPIURL, "/") + path
		hostHeader = cfAPIHost
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.apiKey)
	if hostHeader != "" {
		req.Header.Set("Host", hostHeader)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

// SearchModBySlug searches for a mod by its slug name
func (c *Client) SearchModBySlug(slug string) (int, error) {
	// Hytale game ID is 70216
	data, err := c.apiGet(fmt.Sprintf("/v1/mods/search?gameId=70216&slug=%s", slug))
	if err != nil {
		return 0, err
	}

	var resp struct {
		Data []struct {
			ID   int    `json:"id"`
			Slug string `json:"slug"`
		} `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return 0, err
	}

	for _, mod := range resp.Data {
		if strings.EqualFold(mod.Slug, slug) {
			return mod.ID, nil
		}
	}

	if len(resp.Data) > 0 {
		return resp.Data[0].ID, nil
	}

	return 0, fmt.Errorf("mod not found: %s", slug)
}

// GetModFile fetches a specific mod file
func (c *Client) GetModFile(modID, fileID int) (*ModFile, error) {
	data, err := c.apiGet(fmt.Sprintf("/v1/mods/%d/files/%d", modID, fileID))
	if err != nil {
		return nil, err
	}

	var resp struct {
		Data ModFile `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return nil, err
	}

	return &resp.Data, nil
}

// GetDownloadURL fetches the download URL for a mod file
func (c *Client) GetDownloadURL(modID, fileID int) (string, error) {
	data, err := c.apiGet(fmt.Sprintf("/v1/mods/%d/files/%d/download-url", modID, fileID))
	if err != nil {
		return "", err
	}

	var resp struct {
		Data string `json:"data"`
	}
	if err := json.Unmarshal(data, &resp); err != nil {
		return "", err
	}

	return resp.Data, nil
}

// ResolveBestFile finds the best matching file for a mod
func (c *Client) ResolveBestFile(modID int, partial string, releaseChannel string, gameVersionFilter string) (*ModFile, error) {
	allowedTypes := getAllowedReleaseTypes(releaseChannel)
	
	var bestFile *ModFile
	var bestDate string
	index := 0
	pageSize := 50

	for {
		data, err := c.apiGet(fmt.Sprintf("/v1/mods/%d/files?index=%d&pageSize=%d", modID, index, pageSize))
		if err != nil {
			break
		}

		var resp struct {
			Data []ModFile `json:"data"`
		}
		if err := json.Unmarshal(data, &resp); err != nil {
			break
		}

		if len(resp.Data) == 0 {
			break
		}

		for _, file := range resp.Data {
			if !file.IsAvailable {
				continue
			}
			if !isReleaseTypeAllowed(file.ReleaseType, allowedTypes) {
				continue
			}
			if gameVersionFilter != "" && !containsVersion(file.GameVersions, gameVersionFilter) {
				continue
			}
			if partial != "" && !matchesPartial(file, partial) {
				continue
			}
			if bestDate == "" || file.FileDate > bestDate {
				fileCopy := file
				bestFile = &fileCopy
				bestDate = file.FileDate
			}
		}

		index += pageSize
		if index >= 10000 {
			break
		}
	}

	if bestFile == nil {
		return nil, fmt.Errorf("no matching file found")
	}

	return bestFile, nil
}

// Manager handles mod installation
type Manager struct {
	client          *Client
	cfg             *config.Config
	modsDir         string
	stateDir        string
	downloadsDir    string
	filesDir        string
	manifestPath    string
	releaseChannel  string
	autoUpdate      bool
	failOnError     bool
	gameVersionFilter string
	prune           bool
}

// NewManager creates a new CurseForge mod manager
// baseDir is where mods should be installed (e.g., /home/kubelize/server)
// dataDir is where state/cache files are stored (e.g., /home/kubelize/server/data)
func NewManager(cfg *config.Config, baseDir, dataDir string) (*Manager, error) {
	client, err := NewClient(cfg)
	if err != nil {
		return nil, err
	}

	// Hytale expects mods in ./mods relative to the server working directory
	modsPath := cfg.GetString("HYTALE_MODS_PATH", filepath.Join(baseDir, "mods"))
	stateDir := filepath.Join(dataDir, ".hytale-curseforge-mods")

	m := &Manager{
		client:          client,
		cfg:             cfg,
		modsDir:         modsPath,
		stateDir:        stateDir,
		downloadsDir:    filepath.Join(stateDir, "downloads"),
		filesDir:        filepath.Join(stateDir, "files"),
		manifestPath:    filepath.Join(stateDir, "manifest.json"),
		releaseChannel:  strings.ToLower(cfg.GetString("HYTALE_CURSEFORGE_RELEASE_CHANNEL", "release")),
		autoUpdate:      cfg.GetBool("HYTALE_CURSEFORGE_AUTO_UPDATE", true),
		failOnError:     cfg.GetBool("HYTALE_CURSEFORGE_FAIL_ON_ERROR", false),
		gameVersionFilter: cfg.GetString("HYTALE_CURSEFORGE_GAME_VERSION_FILTER", ""),
		prune:           cfg.GetBool("HYTALE_CURSEFORGE_PRUNE", false),
	}

	// Ensure directories exist
	for _, dir := range []string{m.modsDir, m.stateDir, m.downloadsDir, m.filesDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return m, nil
}

// InstallMods installs all configured CurseForge mods
func (m *Manager) InstallMods(modRefs string) error {
	if modRefs == "" {
		output.Info("No CurseForge mods configured")
		return nil
	}

	manifest := m.loadManifest()
	installedModIDs := make(map[string]bool)
	errors := 0

	refs := m.expandRefs(modRefs)
	for _, ref := range refs {
		ref = strings.TrimSpace(ref)
		if ref == "" || strings.HasPrefix(ref, "#") {
			continue
		}

		output.Step(fmt.Sprintf("Processing mod: %s", ref))

		modIDOrSlug, fileID, partial, isSlug, err := parseModRef(ref)
		if err != nil {
			output.Warning(err.Error())
			errors++
			continue
		}

		// Resolve slug to mod ID if needed
		var modID string
		var mid int
		if isSlug {
			resolvedID, err := m.client.SearchModBySlug(modIDOrSlug)
			if err != nil {
				output.Warning(fmt.Sprintf("could not find mod '%s': %v", modIDOrSlug, err))
				errors++
				continue
			}
			mid = resolvedID
			modID = strconv.Itoa(resolvedID)
			output.Info(fmt.Sprintf("Resolved '%s' to mod ID %d", modIDOrSlug, resolvedID))
		} else {
			modID = modIDOrSlug
			mid, _ = strconv.Atoi(modID)
		}

		manifest.Mods[modID] = ManifestMod{Reference: ref}

		var modFile *ModFile
		if fileID != "" {
			// Specific file ID requested
			fid, _ := strconv.Atoi(fileID)
			modFile, err = m.client.GetModFile(mid, fid)
			if err != nil {
				output.Warning(fmt.Sprintf("could not resolve %s: %v", ref, err))
				errors++
				continue
			}
		} else {
			// Find best matching file
			modFile, err = m.client.ResolveBestFile(mid, partial, m.releaseChannel, m.gameVersionFilter)
			if err != nil {
				output.Warning(fmt.Sprintf("could not resolve %s: %v", ref, err))
				errors++
				continue
			}
		}

		// Check if already installed
		if existing, ok := manifest.Mods[modID]; ok && existing.Installed != nil {
			if existing.Installed.FileID == modFile.ID {
				expectedPath := filepath.Join(m.modsDir, existing.Installed.Path)
				if _, err := os.Stat(expectedPath); err == nil {
					installedModIDs[modID] = true
					output.SuccessWithMessage("already installed")
					continue
				}
			}
		}

		// Download and install
		if err := m.downloadAndInstall(mid, modFile, &manifest); err != nil {
			output.Warning(fmt.Sprintf("failed to install mod %s: %v", modID, err))
			errors++
			continue
		}

		installedModIDs[modID] = true
		output.Success()
	}

	// Prune removed mods
	if m.prune {
		m.pruneMods(&manifest, installedModIDs)
	}

	// Save manifest
	m.saveManifest(&manifest)

	if errors > 0 && m.failOnError {
		return fmt.Errorf("%d mod(s) failed to install", errors)
	}

	return nil
}

// downloadAndInstall downloads and installs a mod file
func (m *Manager) downloadAndInstall(modID int, file *ModFile, manifest *Manifest) error {
	downloadURL, err := m.client.GetDownloadURL(modID, file.ID)
	if err != nil {
		return err
	}

	// Download file
	destDir := filepath.Join(m.filesDir, strconv.Itoa(modID), strconv.Itoa(file.ID))
	if err := os.MkdirAll(destDir, 0755); err != nil {
		return err
	}

	destPath := filepath.Join(destDir, file.FileName)
	tmpPath := filepath.Join(m.downloadsDir, fmt.Sprintf("%d-%d.tmp", modID, file.ID))

	if err := m.downloadFile(downloadURL, tmpPath); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// Verify checksum
	if err := m.verifyChecksum(tmpPath, file.Hashes); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// Move to final location
	if err := os.Rename(tmpPath, destPath); err != nil {
		os.Remove(tmpPath)
		return err
	}

	// Create symlink in mods directory
	safeName := safeFilename(file.FileName)
	visibleName := fmt.Sprintf("cf-%d-%d-%s", modID, file.ID, safeName)
	visiblePath := filepath.Join(m.modsDir, visibleName)

	// Remove old symlink if exists
	os.Remove(visiblePath)

	// Try symlink first, fall back to copy
	if err := os.Symlink(destPath, visiblePath); err != nil {
		// Copy instead
		if err := copyFile(destPath, visiblePath); err != nil {
			return err
		}
	}

	// Update manifest
	modIDStr := strconv.Itoa(modID)
	entry := manifest.Mods[modIDStr]
	
	var hashInfo *HashInfo
	for _, h := range file.Hashes {
		if h.Algo == 1 {
			hashInfo = &HashInfo{Algo: "sha1", Value: h.Value}
			break
		} else if h.Algo == 2 {
			hashInfo = &HashInfo{Algo: "md5", Value: h.Value}
		}
	}

	entry.Resolved = &ResolvedFile{
		FileID:      file.ID,
		FileName:    file.FileName,
		DisplayName: file.DisplayName,
		FileDate:    file.FileDate,
		ReleaseType: file.ReleaseType,
		DownloadURL: downloadURL,
		Hash:        hashInfo,
	}
	entry.Installed = &InstalledFile{
		FileID:           file.ID,
		Path:             visibleName,
		InstalledAtEpoch: time.Now().Unix(),
	}
	manifest.Mods[modIDStr] = entry

	return nil
}

// downloadFile downloads a file from URL
func (m *Manager) downloadFile(url, destPath string) error {
	resp, err := m.client.httpClient.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

// verifyChecksum verifies file checksum
func (m *Manager) verifyChecksum(path string, hashes []Hash) error {
	if len(hashes) == 0 {
		return nil // No checksum to verify
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	for _, h := range hashes {
		var computed string
		switch h.Algo {
		case 1: // SHA1
			sum := sha1.Sum(data)
			computed = hex.EncodeToString(sum[:])
		case 2: // MD5
			sum := md5.Sum(data)
			computed = hex.EncodeToString(sum[:])
		default:
			continue
		}

		if !strings.EqualFold(computed, h.Value) {
			return fmt.Errorf("checksum mismatch: expected %s, got %s", h.Value, computed)
		}
		return nil
	}

	return nil
}

// loadManifest loads or creates the manifest
func (m *Manager) loadManifest() Manifest {
	manifest := Manifest{
		SchemaVersion:  1,
		LastCheckEpoch: 0,
		Mods:           make(map[string]ManifestMod),
	}

	data, err := os.ReadFile(m.manifestPath)
	if err != nil {
		return manifest
	}

	if err := json.Unmarshal(data, &manifest); err != nil {
		return manifest
	}

	if manifest.Mods == nil {
		manifest.Mods = make(map[string]ManifestMod)
	}

	return manifest
}

// saveManifest saves the manifest
func (m *Manager) saveManifest(manifest *Manifest) {
	manifest.LastCheckEpoch = time.Now().Unix()
	
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return
	}

	os.WriteFile(m.manifestPath, data, 0644)
}

// expandRefs expands mod references (including @file references)
func (m *Manager) expandRefs(input string) []string {
	var refs []string
	
	for _, line := range strings.Split(input, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Handle @file references
		if strings.HasPrefix(line, "@") {
			filePath := strings.TrimPrefix(line, "@")
			if data, err := os.ReadFile(filePath); err == nil {
				for _, fline := range strings.Split(string(data), "\n") {
					fline = strings.TrimSpace(fline)
					if fline != "" && !strings.HasPrefix(fline, "#") {
						refs = append(refs, fline)
					}
				}
			}
			continue
		}

		// Split by spaces for multiple refs on one line
		for _, token := range strings.Fields(line) {
			refs = append(refs, token)
		}
	}

	return refs
}

// pruneMods removes mods no longer in the config
func (m *Manager) pruneMods(manifest *Manifest, desired map[string]bool) {
	for modID, mod := range manifest.Mods {
		if !desired[modID] {
			if mod.Installed != nil && mod.Installed.Path != "" {
				os.Remove(filepath.Join(m.modsDir, mod.Installed.Path))
			}
			os.RemoveAll(filepath.Join(m.filesDir, modID))
			delete(manifest.Mods, modID)
		}
	}
}

// Helper functions

func parseModRef(ref string) (modIDOrSlug, fileID, partial string, isSlug bool, err error) {
	// Format: modID, modID:fileID, modID@partial, slug, slug:fileID, slug@partial
	if strings.Contains(ref, ":") {
		parts := strings.SplitN(ref, ":", 2)
		modIDOrSlug = parts[0]
		fileID = parts[1]
	} else if strings.Contains(ref, "@") {
		parts := strings.SplitN(ref, "@", 2)
		modIDOrSlug = parts[0]
		partial = parts[1]
	} else {
		modIDOrSlug = ref
	}

	// Check if modID is numeric or a slug
	if _, numErr := strconv.Atoi(modIDOrSlug); numErr != nil {
		// It's a slug (non-numeric)
		isSlug = true
	}

	if fileID != "" {
		if _, err := strconv.Atoi(fileID); err != nil {
			return "", "", "", false, fmt.Errorf("invalid mod reference: %s (fileID must be numeric)", ref)
		}
	}

	return modIDOrSlug, fileID, partial, isSlug, nil
}

func getAllowedReleaseTypes(channel string) []ReleaseType {
	switch channel {
	case "beta":
		return []ReleaseType{ReleaseTypeRelease, ReleaseTypeBeta}
	case "alpha", "any":
		return []ReleaseType{ReleaseTypeRelease, ReleaseTypeBeta, ReleaseTypeAlpha}
	default: // release
		return []ReleaseType{ReleaseTypeRelease}
	}
}

func isReleaseTypeAllowed(releaseType int, allowed []ReleaseType) bool {
	for _, t := range allowed {
		if int(t) == releaseType {
			return true
		}
	}
	return false
}

func containsVersion(versions []string, filter string) bool {
	for _, v := range versions {
		if v == filter {
			return true
		}
	}
	return false
}

func matchesPartial(file ModFile, partial string) bool {
	partial = strings.ToLower(partial)
	return strings.Contains(strings.ToLower(file.FileName), partial) ||
		strings.Contains(strings.ToLower(file.DisplayName), partial)
}

func safeFilename(name string) string {
	re := regexp.MustCompile(`[^A-Za-z0-9._-]`)
	return re.ReplaceAllString(strings.ReplaceAll(name, " ", "_"), "_")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}
