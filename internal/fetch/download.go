package fetch

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// Downloads all pages of a chapter
// ans saves them into destDir.
func DownloadChapter(slug , chapter , destDir string, writer io.Writer) error {
	baseName := strings.ReplaceAll(slug, "-", " ")
	baseURL := fmt.Sprintf("https://anime-sama.fr/s2/scans/%s/%s", strings.Title(baseName), chapter)

	const defaultDirPerm = 0755
	// Creates dest directory if needed
	if err := os.MkdirAll(destDir, defaultDirPerm); err != nil {
		fmt.Fprintf(writer, "[E] Failed to create destination folder: %v", err)
		return fmt.Errorf("Failed to create destDir: %v", err)
	}
	
	fmt.Fprintln(writer, "Using directory: ", destDir)
	for page := 1; ; page++ {
		imgURL := fmt.Sprintf("%s/%d.jpg", baseURL, page)
		fmt.Fprintln(writer, "Downloading: ", imgURL)

		resp, err := http.Get(imgURL)
		if err != nil {
			fmt.Fprintf(writer, "HTTP.GET failed: %v", err)
			return fmt.Errorf("HTTP.GET failed: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			fmt.Fprintln(writer, "[L] No more pages or invalid response: ", resp.StatusCode)
			// Stop if an image is missing.
			// It probably means that we reached the
			// end of the current chapter.
			break
		}
		
		// Save jpg file with 3
		imgPath := filepath.Join(destDir, fmt.Sprintf("%03d.jpg", page))
		outFile, err := os.Create(imgPath)
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to create image file: %s", imgPath)
			return fmt.Errorf("Failed to create image file: %v", err)
		}

		_, err = io.Copy(outFile, resp.Body)
		outFile.Close()
		if err != nil {
			fmt.Fprintf(writer, "[E] Failed to save image: %s", outFile)
			return fmt.Errorf("Failed to save image: %v", err)
		}
	}
	return nil
}
