package app

import (
	"fmt"
	"log"
	"os"
	"io"
	"path/filepath"

	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/convert"
)

func Run(opts Options, writer io.Writer) {
	//opts := ParseFlags()
	logger := log.New(writer, "", log.LstdFlags)
	if !opts.All && opts.Range == [2]int{} {
		fmt.Fprintln(writer, "[E] Please use --all or --range to download chapters.")
		return
	}

	fmt.Fprintf(writer, "[L] Fetching available chapters for: %s\n", opts.Slug)

	chapters, err := fetch.GetChapters(opts.Slug, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to fecth chapters: %v", err)
	}

	fmt.Fprintf(writer, "[L] Found %d chapters.\n", len(chapters))

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
		fmt.Fprintf(writer, "[L] Filtered to %d chapters from range %d-%d.\n", len(chapters), opts.Range[0], opts.Range[1])
	}

		homeDir, err := os.UserHomeDir()
		if err != nil {
			logger.Fatalf("[L] Failed to get user home directory: %v", err)
		}

	for _, chapter := range chapters {
		fmt.Fprintf(writer, "[L] Downloading chapter %s...\n", chapter)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chapter) 
		err = fetch.DownloadChapter(opts.Slug, chapter, imageDir, writer)
		if err != nil {
			logger.Printf("Failed to download chapter %s: %v", chapter, err)
			continue
		}

		// Create output path if not exists
		err := os.MkdirAll(opts.ScanDir, os.ModePerm)
		if err != nil {
			logger.Printf("Failed to create output directory %s: %v", opts.ScanDir, err)
			return
		}

		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", opts.Slug, chapter))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup, writer)
		if err != nil {
			logger.Printf("Failed to create PDF for chapter %s: %v\n", chapter, err)
			continue
		}

		fmt.Fprintf(writer, "[L] Chapter %s downloaded and saved as %s\n", chapter, pdfPath)
	}

	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		fmt.Fprintf(writer, "[L] Cleaning up images directory: %s\n", rootImagesDir)
		err := os.RemoveAll(rootImagesDir)
		if err != nil {
			logger.Fatalf("[L] Failed to remove dir: %w", err)
		}
	}
}
