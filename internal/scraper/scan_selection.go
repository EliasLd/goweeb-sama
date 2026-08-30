package scraper

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/EliasLd/scan-scraper/internal/logger"

	"github.com/PuerkitoBio/goquery"
)

type ScanPathResult struct {
	Label string
	Path  string
}

// Fetches the manga page, extracts all available scan paths
// from panneauScan() function calls.
// Returns the path(s) (e.g., "scan/vf", "scan/va", "Scans (couleur)", etc...")
func GetAllScanPaths(mangaURL string, log *logger.Logger) ([]ScanPathResult, error) {
	log.Debug("Fetching scan path from: %s\n", mangaURL)

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", mangaURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch manga page: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("manga page returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read manga page: %w", err)
	}
	html := string(body)

	var results []ScanPathResult
	seen := map[string]struct{}{}

	// DOM parsing: <a href="...scan/...">
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err == nil {
		doc.Find("a[href]").Each(func(i int, link *goquery.Selection) {
			href, ok := link.Attr("href")
			if !ok {
				return
			}
			h := strings.TrimSpace(strings.ToLower(href))
			if !(strings.Contains(h, "/scan/") || strings.HasPrefix(h, "scan/")) {
				return
			}
			if _, exists := seen[href]; exists {
				return
			}
			seen[href] = struct{}{}

			label := strings.TrimSpace(link.Text())
			if label == "" {
				label = href
			}

			results = append(results, ScanPathResult{
				Label: label,
				Path:  href,
			})
		})
	}

	// Fallback: parse panneauScan("Label","path")
	if len(results) == 0 {
		// capture: panneauScan("Scans", "scan/vf")
		// groups: [1]=label, [2]=path
		re := regexp.MustCompile(`panneauScan\(\s*["']([^"']+)["']\s*,\s*["']([^"']+)["']\s*\)`)
		matches := re.FindAllStringSubmatch(html, -1)

		for _, m := range matches {
			if len(m) < 3 {
				continue
			}

			label := strings.TrimSpace(m[1])
			path := strings.TrimSpace(m[2])

			p := strings.ToLower(path)
			if !(strings.Contains(p, "/scan/") || strings.HasPrefix(p, "scan/")) {
				continue
			}

			if _, exists := seen[path]; exists {
				continue
			}
			seen[path] = struct{}{}

			if label == "" {
				label = path
			}

			results = append(results, ScanPathResult{
				Label: label,
				Path:  path,
			})
		}
	}

	if len(results) == 0 {
		return nil, fmt.Errorf("no scan paths found")
	}

	log.Info("Found %d scan option(s)\n", len(results))
	return results, nil
}

// Displays scan options and asks user to choose
func PromptUserToSelectScanPath(scanPaths []ScanPathResult, log *logger.Logger) (string, error) {
	if len(scanPaths) == 0 {
		return "", fmt.Errorf("no scan paths to choose from")
	}

	if len(scanPaths) == 1 {
		log.Debug("Only one scan option: %s\n", scanPaths[0].Label)
		log.Debug("Auto-selecting: %s\n", scanPaths[0].Path)
		return scanPaths[0].Path, nil
	}

	// Display all options
	log.Info("\nMultiple scan versions found:\n")
	for i, sp := range scanPaths {
		log.Info(" [%d] - %s\n", i+1, sp.Label)
	}
	log.Info("")

	// Prompt user
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Info("Select a scan version (1-%d) or 0 to cancel: ", len(scanPaths))

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			log.Warn("Invalid input. Please enter a number")
			continue
		}

		if choice == 0 {
			log.Warn("Selection cancelled by user")
			return "", nil
		}

		if choice < 1 || choice > len(scanPaths) {
			log.Warn("Invalid choice. Please enter a number between 1 and %d.\n", len(scanPaths))
			continue
		}

		selected := scanPaths[choice-1]
		log.Debug("Selected: %s\n", selected.Label)
		log.Debug("Path: %s\n", selected.Path)
		return selected.Path, nil
	}
}

// Removes extra spaces and normalizes slashes
func CleanURL(baseURL, path string) string {
	baseURL = strings.TrimSpace(baseURL)
	baseURL = strings.TrimSuffix(baseURL, "/")
	path = strings.TrimPrefix(path, "/")
	path = strings.TrimSuffix(path, "/")
	return baseURL + "/" + path + "/"
}
