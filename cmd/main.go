package main

import (
	"os"

	"github.com/EliasLd/scan-scraper/internal/app"
)

func main() {
	opts := app.ParseFlags()
	app.Run(opts, os.Stdout)
}
