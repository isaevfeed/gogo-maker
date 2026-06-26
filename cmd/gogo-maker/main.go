package main

import (
	"gogo-maker/cmd/go-helper/service_provider"
	"gogo-maker/internal/config"
	"gogo-maker/internal/logger"
	"gogo-maker/internal/models"
	"os"
)

func main() {
	log := logger.New()

	// Check if this is init command - config not needed
	if len(os.Args) >= 2 && models.Command(os.Args[1]) == models.CommandInit {
		sp := service_provider.New(nil, log)
		if err := sp.App().Run(); err != nil {
			log.Fatal(err)
		}
		return
	}

	cfg, err := config.Make()
	if err != nil {
		log.Fatal(err)
		return
	}

	sp := service_provider.New(cfg, log)
	if err := sp.App().Run(); err != nil {
		log.Fatal(err)
	}
}
