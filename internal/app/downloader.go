package app

import (
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/EliasLd/scan-scraper/internal/convert"
	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/scraper"
)

func Run(opts Options, writer io.Writer) {
	if !opts.All && opts.Range == [2]int{} {
		fmt.Fprintln(writer, "[E] Please use --all or --range to download chapters.")
		return
	}

	// Get active domain
	activeDomain := fetch.GetActiveDomain(opts.CustomDomain, writer)

	fmt.Fprintf(writer, "[L] Searching for manga: %s\n", opts.Slug)

	results, err := scraper.SearchCatalog(activeDomain, opts.Slug, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to search catalog: %v\n", err)
		return
	}

	if len(results) == 0 {
		fmt.Fprintf(writer, "[E] Manga '%s' not found in catalog\n", opts.Slug)
		return
	}

	mangaURL, err := scraper.PromptUserToSelectManga(results, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to select manga: %v\n", err)
		return
	}

	if mangaURL == "" {
		fmt.Fprintln(writer, "[L] No manga selected. Exiting.")
		return
	}

	scanPaths, err := scraper.GetAllScanPaths(mangaURL, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to get scan paths: %v\n", err)
		return
	}

	scanPath, err := scraper.PromptUserToSelectScanPath(scanPaths, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to select scan path: %v\n", err)
		return
	}

	if scanPath == "" {
		fmt.Fprintln(writer, "[L] No scan version selected. Exiting.")
		return
	}

	scanPageURL := scraper.CleanURL(mangaURL, scanPath)

	mangaName, err := scraper.ExtractMangaName(scanPageURL, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to extract manga name: %v\n", err)
		return
	}

	scanInfo, err := scraper.GetScanInfo(activeDomain, mangaName, writer)
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to get scan info: %v\n", err)
		return
	}

	// Filter chapters based on user's range
	chapters := scanInfo.Chapters

	if opts.Range != [2]int{} {
		var filtered []int
		for _, ch := range chapters {
			if opts.Range[1] == 0 {
				// Open-ended range (e.g., 10-)
				if ch >= opts.Range[0] {
					filtered = append(filtered, ch)
				}
			} else {
				// Closed range (e.g., 1-5)
				if ch >= opts.Range[0] && ch <= opts.Range[1] {
					filtered = append(filtered, ch)
				}
			}
		}
		chapters = filtered

		if opts.Range[1] == 0 {
			fmt.Fprintf(writer, "[L] Filtered to %d chapters from %d onwards\n", len(chapters), opts.Range[0])
		} else {
			fmt.Fprintf(writer, "[L] Filtered to %d chapters (%d-%d)\n", len(chapters), opts.Range[0], opts.Range[1])
		}
	}

	if len(chapters) == 0 {
		fmt.Fprintln(writer, "[W] No chapters found in range.")
		return
	}

	fmt.Fprintf(writer, "[L] Downloading %d chapters...\n", len(chapters))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(writer, "[E] Failed to get home dir: %v\n", err)
		return
	}

	// Build base URL with the EXACT manga name
	baseURL := fmt.Sprintf("%s/s2/scans/%s", activeDomain, scanInfo.MangaName)

	// Download each chapter
	for _, chNum := range chapters {
		chStr := fmt.Sprintf("%d", chNum)
		fmt.Fprintf(writer, "[L] Downloading chapter %s...\n", chStr)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chStr)

		// Use the exact base URL
		err = fetch.DownloadChapterFromBaseURL(baseURL, chStr, imageDir, writer)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to download chapter %s: %v\n", chStr, err)
			continue
		}

		// Create output directory
		err = os.MkdirAll(opts.ScanDir, os.ModePerm)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to create output dir: %v\n", err)
			return
		}

		// Convert to PDF
		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", mangaName, chStr))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup, writer)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to create PDF: %v\n", err)
			continue
		}

		fmt.Fprintf(writer, "[L] Chapter %s saved as %s\n", chStr, pdfPath)
	}

	// Cleanup
	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		fmt.Fprintf(writer, "[L] Cleaning up: %s\n", rootImagesDir)
		os.RemoveAll(rootImagesDir)
	}
}
