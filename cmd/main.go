package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/convert"
)

// Holds parsed CLI arguments
type Options struct {
	Slug	string
	All	bool
	Range	[2]int
	ScanDir	string
	Cleanup	bool
}

func ParseFlags() Options {
	// Define flags
	allFlag := flag.Bool("all", false, "Download all available chapters")
	allShort := flag.Bool("a", false, "Download all available chapters (Shorthand)")

	rangeFlag := flag.String("range", "", "Range of chapters to download, e.g., 10-77")
	rangeShort := flag.String("r", "", "Shorthand for --range")

	var scanDir string
	flag.StringVar(&scanDir, "scan-dir", "pdf", "Directory to save the generated PDF files")
	flag.StringVar(&scanDir, "d", "pdf", "Shorthand for --scan-dir")

	cleanupFlag := flag.Bool("cleanup", false, "Delete images after PDF creation")
	cleanupShort := flag.Bool("c", false, "Delete images after PDF creation (shorthand)")

	flag.Parse()

	// Expecting slug (manga title) as a positional argument
	args := flag.Args()
	if len(args) < 1 {
		fmt.Println("Usage: scan-scraper [options] <manga-slug>")
		flag.PrintDefaults()
		os.Exit(1)
	}
	slug := args[0]

	// Resolve final values
	all := *allFlag || *allShort
	dir := scanDir
	cleanup := *cleanupFlag || *cleanupShort

	// Parse chapters range
	rangeStr := *rangeFlag
	if *rangeShort != "" {
		rangeStr = *rangeShort
	}

	var chapterRange [2]int
	if rangeStr != "" {
		var start, end int
		n, err := fmt.Sscanf(rangeStr, "%d-%d", &start, &end)
		if err != nil || n != 2 || start > end {
			fmt.Printf("Invalid range format: %s. Use format: <start>-<end>\n", rangeStr)
			os.Exit(1)
		}
		chapterRange[0] = start
		chapterRange[1] = end
	}

	// Create scanDir if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		const defaultDirPerm = 0755
		err := os.MkdirAll(dir, defaultDirPerm)
		if err != nil {
			log.Fatalf("Failed to create scan-dir (%s): %v", dir, err)
		}
	}

	return Options {
		Slug:	slug,
		All:	all,
		Range:	chapterRange,
		ScanDir:dir,
		Cleanup:cleanup,
	}
}

func main() {
	opts := ParseFlags()

	if opts.All || opts.Range != [2]int{} {
		fmt.Println("Fetching available chapters for: ", opts.Slug)

		chapters, err := fetch.GetChapters(opts.Slug)
		if err != nil {
			log.Fatalf("Failed to fetch chapters: %v", err)
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

		for _, chapter := range chapters {
			fmt.Printf("Downloading chapter %s...\n", chapter)

			imageDir := fmt.Sprintf("/home/guts/Images/%s/%s", opts.Slug, chapter)
			err := fetch.DownloadChapter(opts.Slug, chapter, imageDir)
			if err != nil {
				log.Printf("Failed to download chapter %s: %v", chapter, err)
				continue
			}

			pdfPath := fmt.Sprintf("%s/%s_%s.pdf", opts.ScanDir, opts.Slug, chapter)
			err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup)
			if err != nil {
				log.Printf("Failed to create PDF for chapter %s: %v\n", chapter, err)
				continue
			}

			fmt.Printf("Chapter %s downloaded and saved as %s\n", chapter, pdfPath)
		}
	} else {
		fmt.Println("Please use --all or --range to download chapters.")
	}
}
