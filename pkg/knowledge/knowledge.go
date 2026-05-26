// Package knowledge provides access to public/global knowledge sources.
// It integrates with uPlace (global knowledge), uConnect (network knowledge),
// and uScript (operational scripts) to provide a unified knowledge interface.
package knowledge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Source represents a knowledge source
type Source struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Path        string `json:"path"`
	Type        string `json:"type"` // global, user, script, network
}

// Knowledge provides unified access to knowledge sources
type Knowledge struct {
	basePath string
}

// New creates a new Knowledge instance
func New(basePath string) *Knowledge {
	if basePath == "" {
		basePath = guessBasePath()
	}
	return &Knowledge{basePath: basePath}
}

// ListSources returns all available knowledge sources
func (k *Knowledge) ListSources() ([]Source, error) {
	var sources []Source

	// uPlace - global/public knowledge
	uPlacePath := filepath.Join(k.basePath, "uPlace")
	if info, err := os.Stat(uPlacePath); err == nil && info.IsDir() {
		sources = append(sources, Source{
			Name:        "uPlace",
			Description: "Global/public knowledge and user data",
			Path:        uPlacePath,
			Type:        "global",
		})
		// List subdirectories
		entries, _ := os.ReadDir(uPlacePath)
		for _, entry := range entries {
			if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
				sources = append(sources, Source{
					Name:        "uPlace/" + entry.Name(),
					Description: "Knowledge domain: " + entry.Name(),
					Path:        filepath.Join(uPlacePath, entry.Name()),
					Type:        "global",
				})
			}
		}
	}

	// uConnect - network/connectivity knowledge
	uConnectPath := filepath.Join(k.basePath, "uConnect")
	if info, err := os.Stat(uConnectPath); err == nil && info.IsDir() {
		sources = append(sources, Source{
			Name:        "uConnect",
			Description: "Network connectivity and service knowledge",
			Path:        uConnectPath,
			Type:        "network",
		})
	}

	// uScript - operational scripts
	uScriptPath := filepath.Join(k.basePath, "uScript")
	if info, err := os.Stat(uScriptPath); err == nil && info.IsDir() {
		sources = append(sources, Source{
			Name:        "uScript",
			Description: "Operational scripts and automation",
			Path:        uScriptPath,
			Type:        "script",
		})
	}

	return sources, nil
}

// Query searches across all knowledge sources for a term
func (k *Knowledge) Query(term string) ([]string, error) {
	var results []string

	sources, err := k.ListSources()
	if err != nil {
		return nil, err
	}

	termLower := strings.ToLower(term)
	for _, source := range sources {
		// Search for files containing the term
		filepath.Walk(source.Path, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			// Only search text files
			ext := strings.ToLower(filepath.Ext(path))
			if ext == ".md" || ext == ".txt" || ext == ".yaml" || ext == ".yml" || ext == ".json" {
				data, err := os.ReadFile(path)
				if err == nil && strings.Contains(strings.ToLower(string(data)), termLower) {
					relPath, _ := filepath.Rel(source.Path, path)
					results = append(results, fmt.Sprintf("[%s] %s: %s", source.Name, source.Path, relPath))
				}
			}
			return nil
		})
	}

	return results, nil
}

func guessBasePath() string {
	candidates := []string{
		"/Users/fredbook/Code",
		os.Getenv("HOME") + "/Code",
		os.Getenv("HOME") + "/uDos",
	}
	for _, p := range candidates {
		if _, err := os.Stat(filepath.Join(p, "uPlace")); err == nil {
			return p
		}
	}
	return os.Getenv("HOME") + "/Code"
}
