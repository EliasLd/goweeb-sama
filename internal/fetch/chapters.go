package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"strconv"
	"time"
)

// Deduces the number of chapters
// using the access url pattern
func GetChapters(slug string) ([]string, error) {
	baseName := strings.ReplaceAll(slug, "-", " ")
	baseURL := fmt.Sprintf("https://anime-sama.fr/s2/scans/%s", strings.Title(baseName))

	var chapters []string

	const maxErrorsInARow = 10
	const maxRetries = 3

	errorsCounter := 0
	i := 0

	// Reusable http client with timeout
	client := &http.Client {
		Timeout: 10 * time.Second,
	}

	for {
		url := fmt.Sprintf("%s/%d/1.jpg", baseURL, i)
		fmt.Println("Checking:", url)

		var resp *http.Response
		var err error

		// Retry in case of network error
		for attempt := 1; attempt <= maxRetries; attempt++ {
			req, _ := http.NewRequest("GET", url, nil)
			req.Header.Set("User-Agent", "Mozilla/5.0")

			resp, err = client.Do(req)
			if err == nil {
				break
			}

			if attempt < maxRetries {
				time.Sleep(1 * time.Second)
			}
		}
		
		if err != nil {
			return nil, fmt.Errorf("HTTP GET failed after %d attemps: %v", maxRetries, err)
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			chapters = append(chapters, strconv.Itoa(i))
			errorsCounter = 0
		} else {
			fmt.Println("Did not find url: ", url)
			errorsCounter++
			if errorsCounter >= maxErrorsInARow {
				// We're assuming there's no more chapters to fetch
				fmt.Println("Reached max request failures, resuming...")
				break
			}
		}
		i++
		// Short pause between requests
		time.Sleep(10 * time.Millisecond)
	}
	return chapters, nil
}
