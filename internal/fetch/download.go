package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// DownloadChapterFromBaseURL downloads chapter images using an exact base URL
func DownloadChapterFromBaseURL(baseURL, chapter, destDir string, writer io.Writer) error {
	chapterURL := fmt.Sprintf("%s/%s", baseURL, chapter)

	const defaultDirPerm = 0755
	if err := os.MkdirAll(destDir, defaultDirPerm); err != nil {
		return fmt.Errorf("failed to create destDir: %w", err)
	}

	fmt.Fprintf(writer, "[L] Downloading from: %s\n", chapterURL)

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	for page := 1; ; page++ {
		imgURL := fmt.Sprintf("%s/%d.jpg", chapterURL, page)

		req, err := http.NewRequest("GET", imgURL, nil)
		if err != nil {
			return fmt.Errorf("failed to create request: %w", err)
		}

		req.Header.Set("User-Agent", "Mozilla/5.0")

		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("HTTP GET failed: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			fmt.Fprintf(writer, "[L] No more pages (status %d at page %d)\n", resp.StatusCode, page)
			break
		}

		imgPath := filepath.Join(destDir, fmt.Sprintf("%03d.jpg", page))
		outFile, err := os.Create(imgPath)
		if err != nil {
			resp.Body.Close()
			return fmt.Errorf("failed to create image file: %w", err)
		}

		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()
		resp.Body.Close()

		if err != nil {
			return fmt.Errorf("failed to save image: %w", err)
		}

		fmt.Fprintf(writer, "[L] Downloaded page %d\n", page)
	}

	return nil
}
