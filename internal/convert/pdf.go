package convert

import (
	"fmt"
	"os"
	"image"
	"image/jpeg"
	"path/filepath"
	"sort"
	"strings"

	"github.com/signintech/gopdf"
	"golang.org/x/image/webp"
)

// Opens the file at path and checks
// if it is encoded in WEBP format
func IsWebP(path string) (bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	header := make([]byte, 12)
	_, err = f.Read(header)
	if err != nil {
		return false, err
	}

	if string(header[0:4]) == "RIFF" && string(header[8:12]) == "WEBP" {
		return true, nil
	}
	return false, nil
}

// Converts a WebP image in the jpg format
func WebPToJPG(srcPath, destPath string) error {
	f, err := os.Open(srcPath)
	if err != nil {
		return fmt.Errorf("Failed to open WebP file: %w", err)
	}
	defer f.Close()

	img, err := webp.Decode(f)
	if err != nil {
		return fmt.Errorf("Failed to decode WebP image: %w", err)
	}

	out, err := os.Create(destPath)
	if err != nil {
		return fmt.Errorf("Failed to create JPEG file: %w", err)
	}
	defer out.Close()

	// Encode image into JPEG format
	opt := jpeg.Options{Quality: 90}
	if err := jpeg.Encode(out, img, &opt); err != nil {
		return fmt.Errorf("Failed to encode JPEG: %w", err)
	}

	return nil
}

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
