package console

import (
	"gogo-maker/internal/config"
	"gogo-maker/internal/file"
	"gogo-maker/internal/gitrepo"
	"gogo-maker/internal/logger"
)

type (
	Console struct {
		log *logger.Logger

		file    *file.File
		gitRepo *gitrepo.GitRepo
		cfg     *config.Config
	}
)

func New(log *logger.Logger, file *file.File, gitRepo *gitrepo.GitRepo, cfg *config.Config) *Console {
	return &Console{
		log:     log,
		file:    file,
		gitRepo: gitRepo,
		cfg:     cfg,
	}
}
