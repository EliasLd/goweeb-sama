package convert

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/signintech/gopdf"
)

func ImagesToPDF(imagesDir, outputPDFPath string) error {
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
	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	for _, imageName := range imageFiles {
		imagePath := filepath.Join(imagesDir, imageName)

		pdf.AddPage()

		err := pdf.Image(imagePath, 0, 0, nil)
		if err != nil {
			return fmt.Errorf("Error adding image %s: %w", imageName, err)
		}
	}

	// Save newly created PDF file
	if err := pdf.WritePdf(outputPDFPath); err != nil {
		return fmt.Errorf("Failed to write PDF: %w", err)
	}

	fmt.Printf("PDF generated at: %s\n",outputPDFPath)
	return nil
}
