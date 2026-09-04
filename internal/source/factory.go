package source

import (
	"fmt"
	"strings"

	"github.com/EliasLd/scan-scraper/internal/source/animesama"
	sourcetypes "github.com/EliasLd/scan-scraper/internal/source/types"
)

func New(name, customDomain string) (sourcetypes.Provider, error) {
	switch strings.ToLower(strings.TrimSpace(name)) {
	case "animesama":
		return animesama.New(customDomain), nil
	case "sushiscan":
		return nil, fmt.Errorf("provider sushiscan not implemented yet")
	default:
		return nil, fmt.Errorf("unknown provider: %s", name)
	}
}
