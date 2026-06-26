package file

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetConfigPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("os.UserHomeDir: %w", err)
	}

	cfgPath := filepath.Join(homeDir, cfgDirName, cfgOriginalFileName)

	_, err = os.Stat(cfgPath)
	if err != nil && !os.IsNotExist(err) {
		return "", fmt.Errorf("os.Stat: %w", err)
	}

	return cfgPath, nil
}
