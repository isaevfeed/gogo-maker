package app

import (
	"errors"
	"gogo-maker/internal/config"
	"gogo-maker/internal/console"
	"gogo-maker/internal/logger"
	"gogo-maker/internal/models"
	"os"
)

type (
	App struct {
		version string

		cfg *config.Config
		log *logger.Logger

		console *console.Console
	}
)

func New(version string, cfg *config.Config, log *logger.Logger, console *console.Console) *App {
	return &App{
		version: version,
		cfg:     cfg,
		log:     log,
		console: console,
	}
}

func (a *App) Run() error {
	args := os.Args
	if len(args) < 2 {
		a.logCommands()
		return nil
	}

	cmd := models.Command(args[1])

	switch cmd {
	case models.CommandInit:
		return a.console.Init()
	case models.CommandCreate:
		return a.handleCreateCommand(args)
	}

	return nil
}

func (a *App) handleCreateCommand(args []string) error {
	if len(args) < 4 {
		a.log.Info("Usage: gogo-maker create <project-type> <name> [destination-dir]")
		a.console.ListAvailableProjects()
		return errors.New("wrong number of arguments")
	}

	projectType := models.ProjectType(args[2])
	name := args[3]

	var destDir string
	if len(args) >= 5 {
		destDir = args[4]
	}

	params := console.CreateProjectParams{
		ProjectType: projectType,
		Name:        name,
		DestDir:     destDir,
	}

	return a.console.CreateProject(params)
}

func (a *App) logCommands() {
	a.log.Commands("gogo-maker CLI", []logger.CommandInfo{
		{
			Cmd:  "init",
			Desc: "Initialize config file in ~/.gogo-maker",
		},
		{
			Cmd:  "create <type> <name> [dir]",
			Desc: "Create project from template",
			Examples: []string{
				"gogo-maker create http-server my-api",
				"gogo-maker create http-server my-api ./projects",
			},
		},
	})
}
