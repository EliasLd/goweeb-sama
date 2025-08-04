package fetch

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func generateNameVariants(slug string) []string {
	words := strings.Split(strings.ReplaceAll(slug, "-", " "), " ")
	joined := strings.Join(words, " ")

	// simple variants
	variants := map[string]struct{}{
		strings.Title(joined):                  {},
		strings.ToLower(joined):                {},
		joined:                                 {},
		capitalizeWordsExceptSmallWords(words): {},
	}

	// Try all possible casing variants
	n := len(words)
	maxComb := 1 << n
	for i := 0; i < maxComb; i++ {
		var combo []string
		for j := 0; j < n; j++ {
			word := words[j]
			if (i>>j)&1 == 1 {
				combo = append(combo, strings.Title(strings.ToLower(word)))
			} else {
				combo = append(combo, strings.ToLower(word))
			}
		}
		variants[strings.Join(combo, " ")] = struct{}{}
	}

	var result []string
	for v := range variants {
		result = append(result, v)
	}
	return result
}

func capitalizeWordsExceptSmallWords(words []string) string {
	smallWords := map[string]bool{
		"in": true, "on": true, "the": true, "of": true, "and": true, "or": true,
	}
	for i, word := range words {
		if i == 0 || !smallWords[strings.ToLower(word)] {
			words[i] = strings.Title(strings.ToLower(word))
		} else {
			words[i] = strings.ToLower(word)
		}
	}
	return strings.Join(words, " ")
}

func DetectCorrectMangaName(slug string, chapter int, writer io.Writer) (string, error) {
	client := &http.Client{Timeout: 10 * time.Second}
	variants := generateNameVariants(slug)

	for _, name := range variants {
		testURL := fmt.Sprintf("https://anime-sama.fr/s2/scans/%s/%d/1.jpg", name, chapter)
		fmt.Fprintln(writer, "[L] Trying variant:", testURL)

		req, _ := http.NewRequest("GET", testURL, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		resp, err := client.Do(req)
		if err == nil && resp.StatusCode == http.StatusOK {
			io.Copy(io.Discard, resp.Body)
			resp.Body.Close()
			return name, nil
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	return "", fmt.Errorf("[E] could not detect correct casing for slug: %s", slug)
}

func atoiSafe(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
