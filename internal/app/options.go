package app

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

// Holds parsed CLI arguments
type Options struct {
	Slug    string
	All     bool
	Range   [2]int // [0] = start, [1] = end (0 means open-ended)
	ScanDir string
	Cleanup bool
}

func ParseFlags() Options {
	// Define flags
	allFlag := flag.Bool("all", false, "Download all available chapters")
	allShort := flag.Bool("a", false, "Shortand for --all)")

	rangeFlag := flag.String("range", "", "Range of chapters to download, e.g., 10-77, 14-")
	rangeShort := flag.String("r", "", "Shorthand for --range")

	var scanDir string
	flag.StringVar(&scanDir, "scan-dir", "pdf", "Directory to save the generated PDF files")
	flag.StringVar(&scanDir, "d", "pdf", "Shorthand for --scan-dir")

	keepImagesFlag := flag.Bool("keep-images", false, "Keep images after PDF creation")
	keepImagesShort := flag.Bool("k", false, "Shorthand for --keep-images")

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
	keepImages := *keepImagesFlag || *keepImagesShort

	// Parse chapters range
	rangeStr := *rangeFlag
	if *rangeShort != "" {
		rangeStr = *rangeShort
	}

	var chapterRange [2]int
	if rangeStr != "" {
		// Handle open-ended range like "10-"
		if strings.HasSuffix(rangeStr, "-") {
			var start int
			trimmed := strings.TrimSuffix(rangeStr, "-")
			n, err := fmt.Sscanf(trimmed, "%d", &start)
			if err != nil || n != 1 {
				fmt.Printf("Invalid range format: %s. Use format: <start>-<end> or <start>-\n", rangeStr)
				os.Exit(1)
			}
			chapterRange[0] = start
			chapterRange[1] = 0 // 0 means open-ended
		} else {
			// Normal range like "10-20"
			var start, end int
			n, err := fmt.Sscanf(rangeStr, "%d-%d", &start, &end)
			if err != nil || n != 2 || start > end {
				fmt.Printf("Invalid range format: %s. Use format: <start>-<end> or <start>-\n", rangeStr)
				os.Exit(1)
			}
			chapterRange[0] = start
			chapterRange[1] = end
		}
	}

	// Create scanDir if it doesn't exist
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		const defaultDirPerm = 0755
		err := os.MkdirAll(dir, defaultDirPerm)
		if err != nil {
			log.Fatalf("Failed to create scan-dir (%s): %v", dir, err)
		}
	}

	return Options{
		Slug:    slug,
		All:     all,
		Range:   chapterRange,
		ScanDir: dir,
		Cleanup: !keepImages,
	}
}
