package config

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
)

const (
	ProjectConfigFileName = "keysync.json"
	ProjectConfigDir      = ".keysync"
)

type ProjectConfig struct {
	ID   string   `json:"id"`
	Name string   `json:"name"`
	Keys []string `json:"keys"` // List of allowed SSH public keys
}

// LoadProjectConfig looks for keysync.json in the current working directory
func LoadProjectConfig(cwd string) (*ProjectConfig, error) {
	path := filepath.Join(cwd, ProjectConfigFileName)
	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return nil, nil // Not found
	}
	if err != nil {
		return nil, err
	}

	var cfg ProjectConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveProjectConfig saves the project config to keysync.json in the current directory
func SaveProjectConfig(cwd string, cfg *ProjectConfig) error {
	path := filepath.Join(cwd, ProjectConfigFileName)

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// IsProjectInitialized checks if keysync.json exists in the given directory
func IsProjectInitialized(dir string) (bool, error) {
	path := filepath.Join(dir, ProjectConfigFileName)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// AddKey adds a public key to the project if it doesn't already exist
func (p *ProjectConfig) AddKey(key string) error {
	for _, k := range p.Keys {
		if k == key {
			return errors.New("key already exists in project")
		}
	}
	p.Keys = append(p.Keys, key)
	return nil
}

// RemoveKey removes a public key from the project
func (p *ProjectConfig) RemoveKey(key string) error {
	for i, k := range p.Keys {
		if k == key {
			p.Keys = append(p.Keys[:i], p.Keys[i+1:]...)
			return nil
		}
	}
	return errors.New("key not found in project")
}
