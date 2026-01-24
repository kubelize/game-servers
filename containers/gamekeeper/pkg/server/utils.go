package server

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// ensureDir creates a directory if it doesn't exist
func ensureDir(path string) error {
	return os.MkdirAll(path, 0755)
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// downloadFile downloads a file from a URL
func downloadFile(url, dest string) error {
	// Create destination directory
	if err := ensureDir(filepath.Dir(dest)); err != nil {
		return err
	}

	// Create the file
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	// Download
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Write to file
	_, err = io.Copy(out, resp.Body)
	return err
}

// splitArgs splits a string into arguments (respects quotes)
func splitArgs(s string) []string {
	if s == "" {
		return []string{}
	}
	
	var args []string
	var current strings.Builder
	inQuote := false
	
	for _, r := range s {
		switch r {
		case ' ':
			if inQuote {
				current.WriteRune(r)
			} else if current.Len() > 0 {
				args = append(args, current.String())
				current.Reset()
			}
		case '"':
			inQuote = !inQuote
		default:
			current.WriteRune(r)
		}
	}
	
	if current.Len() > 0 {
		args = append(args, current.String())
	}
	
	return args
}
