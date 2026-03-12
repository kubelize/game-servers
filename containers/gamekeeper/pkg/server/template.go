package server

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

const (
	defaultTemplatePath = "/usr/local/share/game-templates/serverconfig.template"
	defaultConfigPath   = "/home/kubelize/config-data/config-values.yaml"
	defaultPasswordPath = "/home/kubelize/config-data/serverpassword.yaml"
)

type templateRenderSpec struct {
	TemplatePath string
	OutputPath   string
	ConfigPath   string
	PasswordPath string
}

// renderTemplateWithGomplate renders a template to the target output path.
// It supports "config" and "password" datasources used by existing templates.
func renderTemplateWithGomplate(spec templateRenderSpec) error {
	if spec.TemplatePath == "" || spec.OutputPath == "" {
		return fmt.Errorf("template and output paths are required")
	}
	if !fileExists(spec.TemplatePath) {
		return fmt.Errorf("template file does not exist: %s", spec.TemplatePath)
	}
	if !fileExists(spec.ConfigPath) {
		return fmt.Errorf("config values file does not exist: %s", spec.ConfigPath)
	}

	if err := ensureDir(filepath.Dir(spec.OutputPath)); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	passwordPath := spec.PasswordPath
	if passwordPath == "" {
		passwordPath = defaultPasswordPath
	}
	if !fileExists(passwordPath) {
		// Keep templates renderable even when password secret isn't mounted.
		fallback, err := os.CreateTemp("", "gamekeeper-password-*.yaml")
		if err != nil {
			return fmt.Errorf("failed to create temporary password file: %w", err)
		}
		defer os.Remove(fallback.Name())

		if _, err := fallback.WriteString("ServerPassword: \"\"\nWebUIPassword: \"\"\n"); err != nil {
			fallback.Close()
			return fmt.Errorf("failed writing temporary password file: %w", err)
		}
		if err := fallback.Close(); err != nil {
			return fmt.Errorf("failed closing temporary password file: %w", err)
		}
		passwordPath = fallback.Name()
	}

	tmpOutput := spec.OutputPath + ".tmp"

	args := []string{
		"-f", spec.TemplatePath,
		"-d", "config=" + spec.ConfigPath,
		"-d", "password=" + passwordPath,
		"-o", tmpOutput,
	}

	cmd := exec.Command("gomplate", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		_ = os.Remove(tmpOutput)
		return fmt.Errorf("gomplate render failed: %w", err)
	}

	if err := os.Rename(tmpOutput, spec.OutputPath); err != nil {
		_ = os.Remove(tmpOutput)
		return fmt.Errorf("failed to move rendered config into place: %w", err)
	}

	return nil
}
