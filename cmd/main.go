package main

import (
	"log"

	"github.com/EliasLd/scan-scraper/internal/tui"
)

func main() {
	if err := tui.Start(); err != nil {
		log.Fatalf("Failed to execute the interface : %v", err)
	}
}

