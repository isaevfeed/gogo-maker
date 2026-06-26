package file

import (
	"fmt"
	"os"
)

func CreateWithPath(path string) error {
	_, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("os.Create: %w", err)
	}

	return nil
}
