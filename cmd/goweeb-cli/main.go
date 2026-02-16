package main

import (
	"os"

	"github.com/EliasLd/scan-scraper/internal/app"
	"github.com/EliasLd/scan-scraper/internal/logger"
)

func main() {
	opts := app.ParseFlags()

	logLevel := logger.LevelInfo
	if opts.Debug {
		logLevel = logger.LevelDebug
	}

	log := logger.New(os.Stdout, logLevel)

	app.Run(opts, log)
}
