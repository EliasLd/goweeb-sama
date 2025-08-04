package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var lowercaseWords = map[string]bool{
	"in": true, "on": true, "the": true, "of": true, "and": true, "or": true,
}

func generateNameVariants(slug string) []string {
	words := strings.Split(strings.ReplaceAll(slug, "-", " "), " ")

	titleCase := strings.Title(strings.Join(words, " "))
	lower := strings.ToLower(strings.Join(words, " "))
	normal := strings.Join(words, " ")
	capitalized := capitalizeExceptSmallWords(words)

	return []string{titleCase, capitalized, lower, normal}
}

func capitalizeExceptSmallWords(words []string) string {
	for i, word := range words {
		if i == 0 || !lowercaseWords[strings.ToLower(word)] {
			words[i] = strings.Title(strings.ToLower(word))
		} else {
			words[i] = strings.ToLower(word)
		}
	}
	return strings.Join(words, " ")
}

// Deduces the number of chapters
// using the access url pattern
func GetChapters(slug string, writer io.Writer) ([]string, error) {
	baseName := strings.ReplaceAll(slug, "-", " ")
	baseURL := fmt.Sprintf("https://anime-sama.fr/s2/scans/%s", strings.Title(baseName))

	var chapters []string

	const maxErrorsInARow = 10
	const maxRetries = 3

	errorsCounter := 0
	i := 0

	// Reusable http client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	for {
		url := fmt.Sprintf("%s/%d/1.jpg", baseURL, i)
		fmt.Fprintln(writer, "Checking:", url)

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
			fmt.Fprintf(writer, "[E] HTTP GET failed after %d attempts: %v", maxRetries, err)
			return nil, err
		}
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			chapters = append(chapters, strconv.Itoa(i))
			errorsCounter = 0
		} else {
			fmt.Fprintln(writer, "Did not find url: ", url)
			errorsCounter++
			if errorsCounter >= maxErrorsInARow {
				// We're assuming there's no more chapters to fetch
				fmt.Fprintln(writer, "Reached max request failures, resuming...")
				break
			}
		}
		i++
		// Short pause between requests
		time.Sleep(10 * time.Millisecond)
	}
	return chapters, nil
}
