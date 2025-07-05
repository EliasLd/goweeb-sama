package fetch

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"
)

// Deduces the number of chapters
// using the access url pattern
func GetChapters(slug string) ([]string, error) {
	baseName := strings.ReplaceAll(slug, "-", " ")
	baseURL := fmt.Sprintf("https://anime-sama.fr/s2/scans/%s", strings.Title(baseName))

	var chapters []string

	const maxErrorsInARow = 10

	errorsCounter := 0
	i := 0

	for {
		url := fmt.Sprintf("%s/%d/1.jpg", baseURL, i)
		fmt.Println("Checking:", url)

		resp, err := http.Get(url)
		if err != nil {
			return nil, fmt.Errorf("http.Get failed: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			chapters = append(chapters, strconv.Itoa(i))
			errorsCounter = 0
		} else {
			errorsCounter++
			fmt.Println("Did not find url: ", url)
			if errorsCounter >= maxErrorsInARow {
				// To much consecutive failures
				// We suppose the manga ends here
				fmt.Println("Reached max request failures, resuming...")
				break
			}
		}
		i++
	}
	return chapters, nil
}
