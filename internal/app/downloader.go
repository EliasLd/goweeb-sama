package app

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/EliasLd/scan-scraper/internal/convert"
	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/logger"
	"github.com/EliasLd/scan-scraper/internal/scraper"
)

func Run(opts Options, log *logger.Logger) {
	if !opts.All && opts.Range == [2]int{} {
		log.Warn("Please use --all or --range to download chapters.")
		return
	}

	// Get active domain
	activeDomain := fetch.GetActiveDomain(opts.CustomDomain, log)

	log.Info("Searching for manga: %s\n", opts.Slug)

	results, err := scraper.SearchCatalog(activeDomain, opts.Slug, log)
	if err != nil {
		log.Error("Failed to search catalog: %v\n", err)
		return
	}

	if len(results) == 0 {
		log.Error("Manga '%s' not found in catalog\n", opts.Slug)
		return
	}

	mangaURL, err := scraper.PromptUserToSelectManga(results, log)
	if err != nil {
		log.Error("Failed to select manga: %v\n", err)
		return
	}

	if mangaURL == "" {
		log.Warn("No manga selected. Exiting.")
		return
	}

	scanPaths, err := scraper.GetAllScanPaths(mangaURL, log)
	if err != nil {
		log.Error("Failed to get scan paths: %v\n", err)
		return
	}

	scanPath, err := scraper.PromptUserToSelectScanPath(scanPaths, log)
	if err != nil {
		log.Error("Failed to select scan path: %v\n", err)
		return
	}

	if scanPath == "" {
		log.Warn("No scan version selected. Exiting.")
		return
	}

	scanPageURL := scraper.CleanURL(mangaURL, scanPath)

	mangaName, err := scraper.ExtractMangaName(scanPageURL, log)
	if err != nil {
		log.Error("Failed to extract manga name: %v\n", err)
		return
	}

	scanInfo, err := scraper.GetScanInfo(activeDomain, mangaName, log)
	if err != nil {
		log.Error("Failed to get scan info: %v\n", err)
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
			log.Info("Filtered to %d chapters from %d onwards\n", len(chapters), opts.Range[0])
		} else {
			log.Info("Filtered to %d chapters (%d-%d)\n", len(chapters), opts.Range[0], opts.Range[1])
		}
	}

	if len(chapters) == 0 {
		log.Warn("No chapters found in range.")
		return
	}

	log.Info("Downloading %d chapters...\n", len(chapters))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home dir: %v\n", err)
		return
	}

	// Build base URL with the EXACT manga name
	baseURL := fmt.Sprintf("%s/s2/scans/%s", activeDomain, scanInfo.MangaName)

	// Download each chapter
	for _, chNum := range chapters {
		chStr := fmt.Sprintf("%d", chNum)
		log.Info("Downloading chapter %s...\n", chStr)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chStr)

		// Use the exact base URL
		err = fetch.DownloadChapterFromBaseURL(baseURL, chStr, imageDir, log)
		if err != nil {
			log.Error("Failed to download chapter %s: %v\n", chStr, err)
			continue
		}

		// Create output directory
		err = os.MkdirAll(opts.ScanDir, os.ModePerm)
		if err != nil {
			log.Error("Failed to create output dir: %v\n", err)
			return
		}

		// Convert to PDF
		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", mangaName, chStr))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup, log)
		if err != nil {
			log.Error("Failed to create PDF: %v\n", err)
			continue
		}

		log.Info("Chapter %s saved as %s\n", chStr, pdfPath)
	}

	// Cleanup
	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		log.Debug("Cleaning up: %s\n", rootImagesDir)
		os.RemoveAll(rootImagesDir)
	}
}

// Runs the download with pre-selected manga URL and scan path (used by TUI)
func RunWithWorkflow(opts Options, log *logger.Logger) {
	if !opts.All && opts.Range == [2]int{} {
		log.Warn("Please use --all or --range to download chapters.")
		return
	}

	activeDomain := fetch.GetActiveDomain(opts.CustomDomain, log)

	// Use pre-selected manga URL and scan path (no prompts)
	scanPageURL := scraper.CleanURL(opts.MangaURL, opts.ScanPath)

	mangaName, err := scraper.ExtractMangaName(scanPageURL, log)
	if err != nil {
		log.Error("Failed to extract manga name: %v\n", err)
		return
	}

	scanInfo, err := scraper.GetScanInfo(activeDomain, mangaName, log)
	if err != nil {
		log.Error("Failed to get scan info: %v\n", err)
		return
	}

	// Filter chapters based on user's range
	chapters := scanInfo.Chapters

	if opts.Range != [2]int{} {
		var filtered []int
		for _, ch := range chapters {
			if opts.Range[1] == 0 {
				if ch >= opts.Range[0] {
					filtered = append(filtered, ch)
				}
			} else {
				if ch >= opts.Range[0] && ch <= opts.Range[1] {
					filtered = append(filtered, ch)
				}
			}
		}
		chapters = filtered

		if opts.Range[1] == 0 {
			log.Info("Filtered to %d chapters from %d onwards\n", len(chapters), opts.Range[0])
		} else {
			log.Info("Filtered to %d chapters (%d-%d)\n", len(chapters), opts.Range[0], opts.Range[1])
		}
	}

	if len(chapters) == 0 {
		log.Warn("No chapters found in range.")
		return
	}

	log.Info("Downloading %d chapters...\n", len(chapters))

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Error("Failed to get home dir: %v\n", err)
		return
	}

	baseURL := fmt.Sprintf("%s/s2/scans/%s", activeDomain, scanInfo.MangaName)

	for _, chNum := range chapters {
		chStr := fmt.Sprintf("%d", chNum)
		log.Info("Downloading chapter %s...\n", chStr)

		imageDir := filepath.Join(homeDir, "Images", opts.Slug, chStr)

		err = fetch.DownloadChapterFromBaseURL(baseURL, chStr, imageDir, log)
		if err != nil {
			log.Error("Failed to download chapter %s: %v\n", chStr, err)
			continue
		}

		err = os.MkdirAll(opts.ScanDir, os.ModePerm)
		if err != nil {
			log.Error("Failed to create output dir: %v\n", err)
			return
		}

		pdfPath := filepath.Join(opts.ScanDir, fmt.Sprintf("%s_%s.pdf", mangaName, chStr))
		err = convert.ImagesToPDF(imageDir, pdfPath, opts.Cleanup, log)
		if err != nil {
			log.Error("Failed to create PDF: %v\n", err)
			continue
		}

		log.Info("Chapter %s saved as %s\n", chStr, pdfPath)
	}

	if opts.Cleanup {
		rootImagesDir := filepath.Join(homeDir, "Images", opts.Slug)
		log.Debug("Cleaning up: %s\n", rootImagesDir)
		os.RemoveAll(rootImagesDir)
	}

	log.Info("Finished downloading")
}
