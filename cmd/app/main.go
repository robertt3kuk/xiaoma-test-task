package main

import (
	"log"

	"github.com/robertt3kuk/xiaoma-test-task/config"
	"github.com/robertt3kuk/xiaoma-test-task/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %s", err)
	}
	app.Run(cfg)
}
