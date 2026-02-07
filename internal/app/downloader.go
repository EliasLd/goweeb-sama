package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/EliasLd/scan-scraper/internal/convert"
	"github.com/EliasLd/scan-scraper/internal/fetch"
)

func Run(opts Options, writer io.Writer) {
	if !opts.All && opts.Range == [2]int{} {
		fmt.Fprintln(writer, "[E] Please use --all or --range to download chapters.")
		return
	}

	fmt.Fprintf(writer, "[L] Fetching available chapters for: %s\n", opts.Slug)

	startChapter := 1
	endChapter := 0 // 0 means open-ended (search until no more chapters)

	if opts.Range != [2]int{} {
		startChapter = opts.Range[0]
		endChapter = opts.Range[1]

		if endChapter == 0 {
			fmt.Fprintf(writer, "[L] Searching from chapter %d to the end...\n", startChapter)
		} else {
			fmt.Fprintf(writer, "[L] Searching for chapters %d to %d...\n", startChapter, endChapter)
		}
	}

	chapters, err := fetch.GetChapters(opts.Slug, startChapter, endChapter, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to fetch chapters: %v\n", err)
		return
	}

	if len(chapters) == 0 {
		fmt.Fprintln(writer, "[W] No chapters found in the specified range.")
		return
	}

	fmt.Fprintf(writer, "[L] Found %d chapters.\n", len(chapters))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to get user home directory: %v\n", err)
		return
	}

	for _, chapter := range chapters {
		fmt.Fprintf(writer, "[L] Downloading chapter %s...\n", chapter)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chapter)
		err = fetch.DownloadChapter(opts.Slug, chapter, imageDir, writer)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to download chapter %s: %v\n", chapter, err)
			continue
		}

		// Create output path if not exists
		err = os.MkdirAll(opts.ScanDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to create output directory %s: %v\n", opts.ScanDir, err)
			return
		}

		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", opts.Slug, chapter))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup, writer)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to create PDF for chapter %s: %v\n", chapter, err)
			continue
		}

		fmt.Fprintf(writer, "[L] Chapter %s downloaded and saved as %s\n", chapter, pdfPath)
	}

	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		fmt.Fprintf(writer, "[L] Cleaning up images directory: %s\n", rootImagesDir)
		err := os.RemoveAll(rootImagesDir)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to remove dir: %v\n", err)
		}
	}
}
