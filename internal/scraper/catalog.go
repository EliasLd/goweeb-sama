package scraper

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/EliasLd/scan-scraper/internal/logger"
	"github.com/PuerkitoBio/goquery"
)

type MangaResult struct {
	Title string
	URL   string
}

// Searches anime-sama catalog and returns the manga page URL
// Returns empty string if not found
func SearchCatalog(domain, query string, log *logger.Logger) ([]MangaResult, error) {
	// Build search URL
	searchURL := fmt.Sprintf("%s/catalogue/?type[]=Scans&search=%s",
		strings.TrimSuffix(domain, "/"),
		url.QueryEscape(query),
	)

	log.Debug("Searching catalog: %s\n", searchURL)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch catalog; %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Catalog returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read catalog response: %w", err)
	}

	html := string(body)

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return nil, fmt.Errorf("failed to parse catalog page: %w", err)
	}

	var results []MangaResult

	doc.Find("div.catalog-card").Each(func(i int, card *goquery.Selection) {
		link := card.Find("a").First()

		href, exists := link.Attr("href")
		if !exists {
			return
		}

		title := strings.TrimSpace(
			card.Find(".card-title").Text(),
		)

		if strings.HasPrefix(href, "/") {
			href = strings.TrimSuffix(domain, "/") + href
		}

		results = append(results, MangaResult{
			Title: title,
			URL:   href,
		})
	})

	if len(results) == 0 {
		log.Warn("No manga found in catalog")
		return nil, nil
	}

	log.Info("Found %d result(s)\n", len(results))

	return results, nil
}

func PromptUserToSelectManga(results []MangaResult, log *logger.Logger) (string, error) {
	if len(results) == 0 {
		return "", fmt.Errorf("no results to choose from")
	}

	if len(results) == 1 {
		log.Debug("Only one result found: %s\n", results[0].Title)
		log.Debug("Auto-selecting: %s\n", results[0].URL)
		return results[0].URL, nil
	}

	// Display all results
	log.Info("\nMultiple results found:\n")
	for i, result := range results {
		log.Info("[%d] - %s\n", i+1, result.Title)
	}
	log.Info("")

	// Prompt user
	reader := bufio.NewReader(os.Stdin)
	for {
		log.Info("Select a manga (1-%d) or 0 to cancel: ", len(results))

		input, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to read input: %w", err)
		}

		input = strings.TrimSpace(input)

		choice, err := strconv.Atoi(input)
		if err != nil {
			log.Warn("Invalid input. Please enter a number.")
			continue
		}

		if choice == 0 {
			log.Debug("Selection cancelled by user")
			return "", nil
		}

		if choice < 1 || choice > len(results) {
			log.Warn("Invalid choice. Please enter a number between 1 and %d.\n", len(results))
			continue
		}

		selected := results[choice-1]
		log.Info("Selected: %s\n", selected.Title)
		log.Debug("URL: %s\n", selected.URL)
		return selected.URL, nil
	}
}
