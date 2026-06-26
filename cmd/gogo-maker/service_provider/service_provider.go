package service_provider

import (
	"gogo-maker/internal/app"
	"gogo-maker/internal/config"
	"gogo-maker/internal/console"
	"gogo-maker/internal/file"
	"gogo-maker/internal/gitrepo"
	"gogo-maker/internal/logger"
)

const version = "v0.0.1"

type (
	ServiceProvider struct {
		cfg *config.Config
		log *logger.Logger

		app     *app.App
		console *console.Console
		file    *file.File
		gitRepo *gitrepo.GitRepo
	}
)

func New(cfg *config.Config, log *logger.Logger) *ServiceProvider {
	return &ServiceProvider{
		cfg: cfg,
		log: log,
	}
}

func (p *ServiceProvider) App() *app.App {
	if p.app == nil {
		p.app = app.New(version, p.cfg, p.log, p.Console())
	}

	return p.app
}

func (p *ServiceProvider) Console() *console.Console {
	if p.console == nil {
		p.console = console.New(p.log, p.File(), p.GitRepo(), p.cfg)
	}

	return p.console
}

func (p *ServiceProvider) File() *file.File {
	if p.file == nil {
		p.file = file.New(p.log)
	}

	return p.file
}

func (p *ServiceProvider) GitRepo() *gitrepo.GitRepo {
	if p.gitRepo == nil {
		p.gitRepo = gitrepo.New(p.log)
	}

	return p.gitRepo
}
