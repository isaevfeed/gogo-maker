package file

import (
	"errors"
	"fmt"
	"gogo-maker/internal/models"
	"os"
	"path/filepath"
	"runtime"

	_ "embed"
)

//go:embed config.example.yaml
var cfgExample []byte

func (f *File) Create() error {
	goOs := models.OperationSystem(runtime.GOOS)

	switch goOs {
	case models.OperationSystemDarwin, models.OperationSystemLinux:
		return f.createUnix()
	case models.OperationSystemWin:
		return f.createWin()
	default:
		return errors.New(fmt.Sprintf("operating system '%s' is not supported", goOs))
	}
}

func (f *File) createUnix() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("os.UserHomeDir: %w", err)
	}

	if err := f.createConfigFileSafe(home, cfgDirName); err != nil {
		return fmt.Errorf("createConfigFileSafe: %w", err)
	}

	f.log.Info("Edit the config file to add your repositories.")

	return nil
}

func (f *File) createConfigFileSafe(homePath, cfgDirName string) error {
	configDir := filepath.Join(homePath, cfgDirName)
	cfgPath := filepath.Join(configDir, cfgOriginalFileName)

	if _, err := os.Stat(cfgPath); err == nil {
		f.log.Infof("Config already exists at %s", cfgPath)

		return nil
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	if err := os.WriteFile(cfgPath, cfgExample, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	f.log.Successf("Config created at %s", cfgPath)

	return nil
}

func (f *File) createWin() error {
	return errors.New("windows is not supported now")
}
