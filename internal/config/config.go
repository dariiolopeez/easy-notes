package config

import (
	"os"
	"os/exec"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	BaseDir string `yaml:"base_dir"`
	Editor  string `yaml:"editor"`
}

func configPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "easy-notes", "config.yaml"), nil
}

func defaultBaseDir() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".local", "share", "easy-notes")
}

func Load() (*Config, error) {
	path, err := configPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return &Config{BaseDir: defaultBaseDir()}, nil
	}
	if err != nil {
		return nil, err
	}

	cfg := &Config{}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}

	if cfg.BaseDir == "" {
		cfg.BaseDir = defaultBaseDir()
	}

	return cfg, nil
}

func (c *Config) Save() error {
	path, err := configPath()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}

	data, err := yaml.Marshal(c)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func (c *Config) ResolveEditor() string {
	if c.Editor != "" {
		return c.Editor
	}
	if v := os.Getenv("EDITOR"); v != "" {
		return v
	}
	if v := os.Getenv("VISUAL"); v != "" {
		return v
	}
	for _, ed := range []string{"nano", "vim"} {
		if _, err := exec.LookPath(ed); err == nil {
			return ed
		}
	}
	return "nano"
}
