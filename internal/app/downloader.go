package app

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/convert"
)

// Holds parsed CLI arguments
type Options struct {
	Slug	string
	All	bool
	Range	[2]int
	ScanDir	string
	Cleanup bool
}

func Run(opts Options) {
	opts := ParseFlags()

	if !opts.All && opts.Range == [2]int{} {
		fmt.Println("Please use --all or --range to download chapters.")
		return
	}

	fmt.Println("Fetching available chapters for: ", opts.Slug)

	chapters, err := fetch.GetChapters(opts.Slug)
	if err != nil {
		log.Fatalf("Failed to fecth chapters: %v", err)
	}

	fmt.Printf("Found %d chapters.\n", len(chapters))

	// Handle chapter range filtering
	if opts.Range != [2]int{} {
		var filtered []string
		for _, ch := range chapters {
			var chapterNum int
			_, err := fmt.Sscanf(ch, "%d", &chapterNum)
			if err == nil && chapterNum >= opts.Range[0] && chapterNum <= opts.Range[1] {
				filtered = append(filtered, ch)
			}
		}
		chapters = filtered
		fmt.Printf("Filtered to %d chapters from range %d-%d.\n", len(chapters), opts.Range[0], opts.Range[1])
	}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("Failed to get user home directory: %v", err)
		}

	for _, chapter := range chapters {
		fmt.Printf("Downloading chapter %s...\n", chapter)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chapter) 
		err = fetch.DownloadChapter(opts.Slug, chapter, imageDir)
		if err != nil {
			log.Printf("Failed to download chapter %s: %v", chapter, err)
			continue
		}

		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", opts.Slug, chapter))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup)
		if err != nil {
			log.Printf("Failed to create PDF for chapter %s: %v\n", chapter, err)
			continue
		}

		fmt.Printf("Chapter %s downloaded and saved as %s\n", chapter, pdfPath)
	}

	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		fmt.Printf("Cleaning up images directory: %s\n", rootImagesDir)
		err := os.RemoveAll(rootImagesDir)
		if err != nil {
			log.Fatalf("Failed to remove dir: %w", err)
		}
	}
}
