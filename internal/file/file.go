package file

import (
	"gogo-maker/internal/logger"
)

type (
	File struct {
		log *logger.Logger
	}
)

func New(log *logger.Logger) *File {
	return &File{
		log: log,
	}
}
