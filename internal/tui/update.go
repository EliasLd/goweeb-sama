package tui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/EliasLd/scan-scraper/internal/app"
	"github.com/EliasLd/scan-scraper/internal/fetch"
	"github.com/EliasLd/scan-scraper/internal/logger"
	"github.com/EliasLd/scan-scraper/internal/scraper"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

type logMsg string

type setupLogPipeMsg struct {
	reader *io.PipeReader
}

// Messages for async operations
type catalogSearchResultMsg struct {
	results []scraper.MangaResult
	err     error
}

type scanPathResultMsg struct {
	paths []scraper.ScanPathResult
	err   error
}

func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	// Handle selection screen separately
	if m.State == StateMangaSelection || m.State == StateScanSelection {
		return handleSelectionUpdate(msg, m)
	}

	// Handle form and downloading states
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case catalogSearchResultMsg:
		if msg.err != nil {
			m.Logs = append(m.Logs, errorStyle.Render(fmt.Sprintf("[E] Catalog search failed: %v", msg.err)))
			return m, nil
		}

		if len(msg.results) == 0 {
			m.Logs = append(m.Logs, errorStyle.Render("[E] No manga found"))
			return m, nil
		}

		// If only one result, auto-select
		if len(msg.results) == 1 {
			m.SelectedMangaURL = msg.results[0].URL
			return m, fetchScanPaths(m.SelectedMangaURL)
		}

		// Multiple results, show selection screen
		items := make([]SelectionItem, len(msg.results))
		for i, r := range msg.results {
			items[i] = SelectionItem{Label: r.Title, Value: r.URL}
		}

		m.SelectionModel = NewSelectionModel("Sélectionnez un manga", items)
		m.SelectionModel.Width = m.Width
		m.SelectionModel.Height = m.Height
		m.State = StateMangaSelection
		return m, nil

	case scanPathResultMsg:
		if msg.err != nil {
			m.Logs = append(m.Logs, errorStyle.Render(fmt.Sprintf("[E] Failed to get scan paths: %v", msg.err)))
			return m, nil
		}

		// If only one scan path, auto-select
		if len(msg.paths) == 1 {
			m.SelectedScanPath = msg.paths[0].Path
			return m, startDownload(m)
		}

		// Multiple scan paths, show selection screen
		items := make([]SelectionItem, len(msg.paths))
		for i, p := range msg.paths {
			items[i] = SelectionItem{Label: p.Label, Value: p.Path}
		}

		m.SelectionModel = NewSelectionModel("Sélectionnez une version", items)
		m.SelectionModel.Width = m.Width
		m.SelectionModel.Height = m.Height
		m.State = StateScanSelection
		return m, nil

	case tea.KeyMsg:
		if m.IsDownloading && msg.String() != "ctrl+c" && msg.String() != "esc" {
			return m, nil
		}

		switch msg.String() {

		case "ctrl+c", "esc":
			return m, tea.Quit

		case "up":
			if m.Cursor > 0 {
				m.Cursor--
			}
			m = updateFocus(m)
			return m, nil

		case "down", "tab":
			if m.Cursor < 6 {
				m.Cursor++
			}
			m = updateFocus(m)
			return m, nil

		case " ":
			// Allow space only in manga Input
			if m.Cursor == 0 {
				var cmd tea.Cmd
				m.MangaInput, cmd = m.MangaInput.Update(msg)
				return m, cmd
			}

			switch m.Cursor {
			case 1:
				m.AllCheckbox.Toggle()
			case 5:
				m.KeepCheckbox.Toggle()
			}
			return m, nil

		case "enter":
			switch m.Cursor {
			case 1: // AllCheckbox
				m.AllCheckbox.Toggle()
			case 5: // KeepCheckbox
				m.KeepCheckbox.Toggle()
			case 6: // Download button
				if m.DownloadReady {
					// Start catalog search
					m.Logs = append(m.Logs, fmt.Sprintf("Searching for: %s", m.MangaInput.Value()))
					return m, searchCatalog(m.MangaInput.Value(), m.DomainInput.Value())
				}
			}
			return m, nil
		}

		// Input handling
		switch m.Cursor {
		case 0:
			var cmd tea.Cmd
			m.MangaInput, cmd = m.MangaInput.Update(msg)
			return m, cmd
		case 2:
			if !m.AllCheckbox.Checked {
				var cmd tea.Cmd
				m.RangeInput, cmd = m.RangeInput.Update(msg)
				return m, cmd
			}
		case 3:
			var cmd tea.Cmd
			m.ScanDirInput, cmd = m.ScanDirInput.Update(msg)
			return m, cmd
		case 4:
			var cmd tea.Cmd
			m.DomainInput, cmd = m.DomainInput.Update(msg)
			return m, cmd
		}

	case logMsg:
		logLine := string(msg)

		switch {
		case logLine == "Finished downloading":
			styled := highlightStyle.Render(string(msg))
			m.Logs = append(m.Logs, styled)
			m.IsDownloading = false
		case strings.HasPrefix(logLine, "[DEBUG]"):
			// Skip debug logs in TUI
		case strings.Contains(logLine, "[ERROR]"):
			styled := errorStyle.Render(logLine)
			m.Logs = append(m.Logs, styled)
		case strings.Contains(logLine, "[!]"):
			// Warning logs
			m.Logs = append(m.Logs, logLine)
		default:
			// Info level logs
			m.Logs = append(m.Logs, logLine)
		}
		return m, readOneLogLine(m)

	case setupLogPipeMsg:
		m.pipeReader = msg.reader
		m.scanner = bufio.NewScanner(m.pipeReader)
		return m, readOneLogLine(m)
	}

	m.DownloadReady = m.MangaInput.Value() != "" &&
		(m.AllCheckbox.Checked || m.RangeInput.Value() != "") &&
		m.ScanDirInput.Value() != ""

	return m, nil
}

func handleSelectionUpdate(msg tea.Msg, m Model) (Model, tea.Cmd) {
	var cmd tea.Cmd
	selectionModel, cmd := m.SelectionModel.Update(msg)
	m.SelectionModel = selectionModel.(SelectionModel)

	// Check if selection was made
	if m.SelectionModel.Selected != "" {
		if m.State == StateMangaSelection {
			m.SelectedMangaURL = m.SelectionModel.Selected
			m.State = StateForm
			m.Logs = append(m.Logs, "Manga selected, fetching scan versions...")
			return m, fetchScanPaths(m.SelectedMangaURL)
		} else if m.State == StateScanSelection {
			m.SelectedScanPath = m.SelectionModel.Selected
			m.State = StateForm
			m.Logs = append(m.Logs, "Scan version selected, starting download...")
			return m, startDownload(m)
		}
	}

	// Check if cancelled
	if m.SelectionModel.Cancelled {
		m.State = StateForm
		m.Logs = append(m.Logs, "Selection cancelled")
		return m, nil
	}

	return m, cmd
}

func updateFocus(m Model) Model {
	m.MangaInput.Blur()
	m.RangeInput.Blur()
	m.ScanDirInput.Blur()
	m.DomainInput.Blur()

	switch m.Cursor {
	case 0:
		m.MangaInput.Focus()
	case 2:
		if !m.AllCheckbox.Checked {
			m.RangeInput.Focus()
		}
	case 3:
		m.ScanDirInput.Focus()
	case 4:
		m.DomainInput.Focus()
	}
	return m
}

// Async commands

func searchCatalog(query, customDomain string) tea.Cmd {
	return func() tea.Msg {
		log := logger.New(io.Discard, logger.LevelInfo)
		domain := fetch.GetActiveDomain(customDomain, log)
		results, err := scraper.SearchCatalog(domain, query, log)
		return catalogSearchResultMsg{results: results, err: err}
	}
}

func fetchScanPaths(mangaURL string) tea.Cmd {
	return func() tea.Msg {
		log := logger.New(io.Discard, logger.LevelInfo)
		paths, err := scraper.GetAllScanPaths(mangaURL, log)
		return scanPathResultMsg{paths: paths, err: err}
	}
}

func startDownload(m Model) tea.Cmd {
	return func() tea.Msg {
		// Build download options
		opts := app.Options{
			Slug:         m.MangaInput.Value(),
			All:          m.AllCheckbox.Checked,
			ScanDir:      m.ScanDirInput.Value(),
			Cleanup:      !m.KeepCheckbox.Checked,
			CustomDomain: strings.TrimSpace(m.DomainInput.Value()),
			MangaURL:     m.SelectedMangaURL,
			ScanPath:     m.SelectedScanPath,
		}

		if !opts.All {
			// Range parsing
			r := strings.TrimSpace(m.RangeInput.Value())
			var rangeMode app.RangeMode = app.RangeNone
			var chapterRange [2]int

			switch {
			case strings.HasPrefix(r, "-") && len(r) > 1:
				// N last chapters
				var nLast int
				n, err := fmt.Sscanf(r, "-%d", &nLast)
				if err == nil && n == 1 && nLast > 0 {
					chapterRange[1] = nLast
					rangeMode = app.RangeLastN
				}
			case strings.HasSuffix(r, "-"):
				// From chapter N to last
				var start int
				trim := strings.TrimSuffix(r, "-")
				n, err := fmt.Sscanf(trim, "%d", &start)
				if err == nil && n == 1 {
					chapterRange[0] = start
					chapterRange[1] = 0
					rangeMode = app.RangeOpenEnded
				}
			default:
				// N-M, from chapter N to chapter M, normal case
				var start, end int
				n, err := fmt.Sscanf(r, "%d-%d", &start, &end)
				if err == nil && n == 2 && start <= end {
					chapterRange[0] = start
					chapterRange[1] = end
					rangeMode = app.RangeNormal
				} else {
					// chapter N alone
					var solo int
					n, err := fmt.Sscanf(r, "%d", &solo)
					if err == nil && n == 1 {
						chapterRange[0] = solo
						chapterRange[1] = solo
						rangeMode = app.RangeNormal
					}
				}
			}

			opts.Range = chapterRange
			opts.RangeMode = rangeMode
		}

		pr, pw := io.Pipe()

		go func() {
			log := logger.New(pw, logger.LevelInfo)
			app.RunWithWorkflow(opts, log)
			pw.Close()
		}()

		return setupLogPipeMsg{reader: pr}
	}
}

func readOneLogLine(m Model) tea.Cmd {
	return func() tea.Msg {
		if m.scanner == nil {
			return nil
		}

		if m.scanner.Scan() {
			line := m.scanner.Text()
			return logMsg(line)
		}

		if err := m.scanner.Err(); err != nil {
			return logMsg(fmt.Sprintf("Error: failed to read log: %v", err))
		}

		if m.IsDownloading {
			return logMsg("Finished downloading")
		}

		return nil
	}
}
