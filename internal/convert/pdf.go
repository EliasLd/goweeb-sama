package convert

import (
	"fmt"
	"os"
	"image"
	_ "image/jpeg"
	"path/filepath"
	"sort"
	"strings"

	"github.com/signintech/gopdf"
)

func ImagesToPDF(imagesDir, outputPDFPath string, cleanup bool) error {
	files, err := os.ReadDir(imagesDir)
	if err != nil {
		return fmt.Errorf("Failed to read image dir: %v", err)
	}
	
	var imageFiles []string
	for _, f := range files {	
		if !f.IsDir() && strings.HasSuffix(strings.ToLower(f.Name()), ".jpg") {
			imageFiles = append(imageFiles, f.Name())
		}
	}

	if len(imageFiles) == 0 {
		return fmt.Errorf("No JPG files found in %s", imagesDir)
	}

	sort.Strings(imageFiles)

	pdf := gopdf.GoPdf{}
	pdf.Start(gopdf.Config{})


	for _, imageName := range imageFiles {
		imagePath := filepath.Join(imagesDir, imageName)

		imageFile, err := os.Open(imagePath)
		if err != nil {
			return fmt.Errorf("Failed to open image %s: %w", imageName, err)
		}

		imageConfig, _, err := image.DecodeConfig(imageFile)
		imageFile.Close()
		if err != nil {
			return fmt.Errorf("Failed to decode image %s: %w", imageName, err)
		}

		width := float64(imageConfig.Width)
		height := float64(imageConfig.Height)

		// Start a new page with exact dimensions
		pdf.AddPageWithOption(gopdf.PageOption{
			PageSize: &gopdf.Rect{W: width, H: height},
		})

		err = pdf.Image(imagePath, 0, 0, &gopdf.Rect{W: width, H: height})
		if err != nil {
			return fmt.Errorf("Error adding image %s: %w", imageName, err)
		}
	}

	// Save newly created PDF file
	if err := pdf.WritePdf(outputPDFPath); err != nil {
		return fmt.Errorf("Failed to write PDF: %w", err)
	}

	fmt.Printf("PDF generated at: %s\n",outputPDFPath)

	if cleanup {
		fmt.Printf("Cleaning up directory: %s\n", imagesDir)
		err := os.RemoveAll(imagesDir)
		if err != nil {
			return fmt.Errorf("Failed to remove temp dir: %w", err)
		}
	}

	return nil
}
