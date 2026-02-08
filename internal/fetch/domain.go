package fetch

import (
	"fmt"
	"io"
	"strings"
)

const (
	// Default domain to use
	DefaultDomain = "https://anime-sama.tv"
)

// Returns the custom domain if provided, otherwise returns the default domain
func GetActiveDomain(customDomain string, writer io.Writer) string {
	// If user specified a custom domain, use it directly
	if customDomain != "" {
		fmt.Fprintf(writer, "[L] Using custom domain: %s\n", customDomain)
		return customDomain
	}

	// Otherwise use default domain
	fmt.Fprintf(writer, "[L] Using default domain: %s\n", DefaultDomain)
	return DefaultDomain
}

// BuildScanBaseURL constructs the base URL for scans
func BuildScanBaseURL(domain, correctName string) string {
	// Remove trailing slash if present
	domain = strings.TrimSuffix(domain, "/")
	return fmt.Sprintf("%s/s2/scans/%s", domain, correctName)
}
