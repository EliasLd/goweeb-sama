package tui

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/EliasLd/scan-scraper/internal/app"
	"github.com/EliasLd/scan-scraper/internal/logger"
	"github.com/EliasLd/scan-scraper/internal/source"
	"github.com/EliasLd/scan-scraper/internal/source/common"
	sourcetypes "github.com/EliasLd/scan-scraper/internal/source/types"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

type logMsg string

type setupLogPipeMsg struct {
	reader *io.PipeReader
}

type catalogSearchResultMsg struct {
	results []sourcetypes.SearchResult
	err     error
}

type scanPathResultMsg struct {
	paths []common.SelectableItem
	err   error
}

func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	if m.State == StateMangaSelection || m.State == StateScanSelection {
		return handleSelectionUpdate(msg, m)
	}

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

		if len(msg.results) == 1 {
			m.SelectedMangaURL = msg.results[0].URL
			m.Logs = append(m.Logs, "Manga found. Fetching scan versions...")
			return m, fetchScanPaths(m.SelectedMangaURL, strings.TrimSpace(m.DomainInput.Value()))
		}

		items := make([]SelectionItem, len(msg.results))
		for i, r := range msg.results {
			items[i] = SelectionItem{Label: r.Title, Value: r.URL}
		}

		m.SelectionModel = NewSelectionModel("Select a manga", items)
		m.SelectionModel.Width = m.Width
		m.SelectionModel.Height = m.Height
		m.State = StateMangaSelection
		return m, nil

	case scanPathResultMsg:
		if msg.err != nil {
			m.Logs = append(m.Logs, errorStyle.Render(fmt.Sprintf("[E] Failed to get scan paths: %v", msg.err)))
			return m, nil
		}
		if len(msg.paths) == 0 {
			m.Logs = append(m.Logs, errorStyle.Render("[E] No scan versions found"))
			return m, nil
		}

		if len(msg.paths) == 1 {
			m.SelectedScanPath = msg.paths[0].Value
			m.IsDownloading = true
			m.State = StateDownloading
			m.Logs = append(m.Logs, "Scan version selected. Starting download...")
			return m, startDownload(m)
		}

		items := make([]SelectionItem, len(msg.paths))
		for i, p := range msg.paths {
			items[i] = SelectionItem{Label: p.Label, Value: p.Value}
		}

		m.SelectionModel = NewSelectionModel("Select a version", items)
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
			if m.Cursor < 7 {
				m.Cursor++
			}
			m = updateFocus(m)
			return m, nil

		case " ":
			if m.Cursor == 0 {
				var cmd tea.Cmd
				m.MangaInput, cmd = m.MangaInput.Update(msg)
				return m, cmd
			}
			switch m.Cursor {
			case 1:
				m.AllCheckbox.Toggle()
			case 5:
				m.EbookCheckbox.Toggle()
			case 6:
				m.KeepCheckbox.Toggle()
			}
			return m, nil

		case "enter":
			switch m.Cursor {
			case 1:
				m.AllCheckbox.Toggle()
			case 5:
				m.EbookCheckbox.Toggle()
			case 6:
				m.KeepCheckbox.Toggle()
			case 7:
				if m.DownloadReady {
					m.Logs = append(m.Logs, fmt.Sprintf("Searching for: %s...", m.MangaInput.Value()))
					return m, searchCatalog(m.MangaInput.Value(), m.DomainInput.Value())
				}
			}
			return m, nil
		}

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

	case setupLogPipeMsg:
		m.pipeReader = msg.reader
		m.scanner = bufio.NewScanner(m.pipeReader)
		buf := make([]byte, 64*1024)
		m.scanner.Buffer(buf, 1024*1024)
		return m, readOneLogLine(m)

	case logMsg:
		logLine := string(msg)

		switch {
		case logLine == "Finished downloading":
			styled := highlightStyle.Render("Download complete!")
			m.Logs = append(m.Logs, styled)
			m.IsDownloading = false
			m.State = StateForm
		case strings.HasPrefix(logLine, "[DEBUG]"):
		case strings.Contains(logLine, "[ERROR]"):
			m.Logs = append(m.Logs, errorStyle.Render(logLine))
		case strings.Contains(logLine, "[!]"):
			m.Logs = append(m.Logs, logLine)
		default:
			m.Logs = append(m.Logs, logLine)
		}

		if m.IsDownloading {
			return m, readOneLogLine(m)
		}
		return m, nil
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

	if m.SelectionModel.Selected != "" {
		if m.State == StateMangaSelection {
			m.SelectedMangaURL = m.SelectionModel.Selected
			m.State = StateForm
			m.Logs = append(m.Logs, "Manga selected. Fetching scan versions...")
			return m, fetchScanPaths(m.SelectedMangaURL, strings.TrimSpace(m.DomainInput.Value()))
		} else if m.State == StateScanSelection {
			m.SelectedScanPath = m.SelectionModel.Selected
			m.State = StateDownloading
			m.IsDownloading = true
			m.Logs = append(m.Logs, "Scan version selected. Starting download...")
			return m, startDownload(m)
		}
	}

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

func searchCatalog(query, customDomain string) tea.Cmd {
	return func() tea.Msg {
		log := logger.New(io.Discard, logger.LevelInfo)

		provider, err := source.New("animesama", strings.TrimSpace(customDomain))
		if err != nil {
			return catalogSearchResultMsg{results: nil, err: err}
		}

		results, err := provider.Search(query, log)
		return catalogSearchResultMsg{results: results, err: err}
	}
}

func fetchScanPaths(mangaURL string, customDomain string) tea.Cmd {
	return func() tea.Msg {
		log := logger.New(io.Discard, logger.LevelInfo)

		provider, err := source.New("animesama", strings.TrimSpace(customDomain))
		if err != nil {
			return scanPathResultMsg{paths: nil, err: err}
		}

		paths, err := provider.ListScanPaths(mangaURL, log)
		if err != nil {
			return scanPathResultMsg{paths: nil, err: err}
		}

		return scanPathResultMsg{paths: paths, err: nil}
	}
}

func startDownload(m Model) tea.Cmd {
	return func() tea.Msg {
		opts := app.Options{
			Slug:          m.MangaInput.Value(),
			All:           m.AllCheckbox.Checked,
			Source:        "animesama",
			ScanDir:       m.ScanDirInput.Value(),
			Cleanup:       !m.KeepCheckbox.Checked,
			CustomDomain:  strings.TrimSpace(m.DomainInput.Value()),
			MangaURL:      m.SelectedMangaURL,
			ScanPath:      m.SelectedScanPath,
			EbookFriendly: m.EbookCheckbox.Checked,
		}

		if !opts.All {
			r := strings.TrimSpace(m.RangeInput.Value())
			var rangeMode app.RangeMode = app.RangeNone
			var chapterRange [2]int

			switch {
			case strings.HasPrefix(r, "-") && len(r) > 1:
				var nLast int
				if n, err := fmt.Sscanf(r, "-%d", &nLast); err == nil && n == 1 && nLast > 0 {
					chapterRange[1] = nLast
					rangeMode = app.RangeLastN
				}
			case strings.HasSuffix(r, "-"):
				var start int
				trim := strings.TrimSuffix(r, "-")
				if n, err := fmt.Sscanf(trim, "%d", &start); err == nil && n == 1 {
					chapterRange[0] = start
					chapterRange[1] = 0
					rangeMode = app.RangeOpenEnded
				}
			default:
				var start, end int
				if n, err := fmt.Sscanf(r, "%d-%d", &start, &end); err == nil && n == 2 && start <= end {
					chapterRange[0] = start
					chapterRange[1] = end
					rangeMode = app.RangeNormal
				} else {
					var solo int
					if n, err := fmt.Sscanf(r, "%d", &solo); err == nil && n == 1 {
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
			_ = pw.Close()
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
			return logMsg(m.scanner.Text())
		}

		if err := m.scanner.Err(); err != nil {
			return logMsg(fmt.Sprintf("Error: failed to read log: %v", err))
		}

		return logMsg("Finished downloading")
	}
}
