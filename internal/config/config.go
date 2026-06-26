package config

import (
	"fmt"
	"gogo-maker/internal/file"
	"gopkg.in/yaml.v3"
	"os"
)

type (
	Config struct {
		Reps []Repo `yaml:"reps"`
	}

	Repo struct {
		Name string `yaml:"name"`
		Url  string `yaml:"url"`
		Type string `yaml:"type"`
	}
)

func Make() (*Config, error) {
	cfgPath, err := file.GetConfigPath()
	if err != nil {
		return nil, fmt.Errorf("file.GetConfigPath: %w", err)
	}

	cfgContent, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile: %w", err)
	}

	var cfg *Config
	if err := yaml.Unmarshal(cfgContent, &cfg); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal: %w", err)
	}

	return cfg, nil
}
