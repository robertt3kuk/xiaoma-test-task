package main

import (
	"github.com/robertt3kuk/xiaoma-test-task/config"
	"github.com/robertt3kuk/xiaoma-test-task/internal/app"
)

func main() {
	cfg, err := config.NewConfig()
	if err != nil {
		panic(err)
	}
	app.Run(cfg)
}
