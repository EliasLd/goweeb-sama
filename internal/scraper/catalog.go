package scraper

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
	"time"
)

// Searches anime-sama catalog and returns the manga page URL
// Returns empty string if not found
func SearchCatalog(domain, query string, writer io.Writer) (string, error) {
	// Build search URL
	searchURL := fmt.Sprintf("%s/catalogue/?type=Scans&search=%s",
		strings.TrimSuffix(domain, "/"),
		url.QueryEscape(query),
	)

	fmt.Fprintf(writer, "[L] Searching catalog: %s\n", searchURL)

	client := &http.Client{
		Timeout: 15 * time.Second,
	}

	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		return "", fmt.Errorf("Failed to create request: %w", err)
	}

	req.Header.Set("User-Agent", "Mozilla/5.0")

	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Failed to fetch catalog; %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Catalog returned status: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read catalog response: %w", err)
	}

	html := string(body)

	// Extract the first manga card link
	// Pattern: <a href="https://anime-sama.tv/catalogue/one-piece/">
	// We look for: <div class="shrink-0 catalog-card card-base"> followed by <a href="...">
	cardRegex := regexp.MustCompile(`(?s)<div class="shrink-0 catalog-card card-base">.*?<a href="([^"]+)"`)
	matches := cardRegex.FindStringSubmatch(html)

	if len(matches) < 2 {
		fmt.Fprintln(writer, "[W] No manga found in catalog")
		return "", nil
	}

	mangaURL := matches[1]
	fmt.Fprintf(writer, "[L] Found manga page: %s\n", mangaURL)

	return mangaURL, nil
}
