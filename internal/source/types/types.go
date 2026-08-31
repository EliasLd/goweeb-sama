package types

import "github.com/EliasLd/scan-scraper/internal/logger"

type ItemKind int

const (
	ItemChapter ItemKind = iota
	ItemVolume
)

type SearchResult struct {
	Title string
	URL   string
}

type Entry struct {
	Number int
	Label  string
	URL    string
}

type Work struct {
	Title string
	Kind  ItemKind
}

type Provider interface {
	Name() string
	Search(query string, log *logger.Logger) ([]SearchResult, error)
	ListEntries(workURL string, log *logger.Logger) (Work, []Entry, error)
	GetPageImageURLs(entryURL string, log *logger.Logger) ([]string, error)
}
