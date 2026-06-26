package console

import (
	"fmt"
)

func (c *Console) Init() error {
	c.log.Title("Initializing go-helper")

	if err := c.file.Create(); err != nil {
		return fmt.Errorf("file.Create: %w", err)
	}

	return nil
}
