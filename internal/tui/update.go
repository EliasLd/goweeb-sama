package tui

import (
	"bufio"
	"fmt"
	app "github.com/EliasLd/scan-scraper/internal/app"
	"io"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var highlightStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("226"))
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("9"))

type logMsg string

type setupLogPipeMsg struct {
	reader *io.PipeReader
}

func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:

		if m.IsDownloading && msg.String() != "ctrl+c" {
			return m, nil
		}

		switch msg.String() {

		case "ctrl+c":
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

		case "enter", " ":
			switch m.Cursor {
			case 1: // AllCheckbox
				m.AllCheckbox.Toggle()
			case 5: // KeepCheckbox
				m.KeepCheckbox.Toggle()
			case 6: // Download button
				if m.DownloadReady {
					opts := app.Options{
						Slug:         m.MangaInput.Value(),
						All:          m.AllCheckbox.Checked,
						ScanDir:      m.ScanDirInput.Value(),
						Cleanup:      !m.KeepCheckbox.Checked,
						CustomDomain: strings.TrimSpace(m.DomainInput.Value()),
					}

					if !opts.All {
						r := strings.TrimSpace(m.RangeInput.Value())

						// Support:
						// - "10-20"
						// - "10-" (open ended)
						var from, to int
						if strings.HasSuffix(r, "-") {
							trim := strings.TrimSuffix(r, "-")
							_, err := fmt.Sscanf(trim, "%d", &from)
							if err == nil {
								opts.Range = [2]int{from, 0}
							}
						} else {
							_, err := fmt.Sscanf(r, "%d-%d", &from, &to)
							if err == nil {
								opts.Range = [2]int{from, to}
							}
						}
					}

					m.IsDownloading = true
					m.Logs = nil
					return m, runDownloadInBackground(opts)
				}
			}
			return m, nil
		}

		// input handling
		switch m.Cursor {
		case 0:
			var cmd tea.Cmd
			m.MangaInput, cmd = m.MangaInput.Update(msg)
			return m, cmd
		case 2:
			// Range only active if not "all" checked
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
		case strings.HasPrefix(logLine, "[L]"):
			m.Logs = append(m.Logs, logLine)
		case strings.HasPrefix(logLine, "[E]"):
			styled := errorStyle.Render(logLine)
			m.Logs = append(m.Logs, styled)
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

// Updates focus/blur of textinput field depending on the cursor
func updateFocus(m Model) Model {
	m.MangaInput.Blur()
	m.RangeInput.Blur()
	m.ScanDirInput.Blur()

	// Apply focus on the selected field
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

func runDownloadInBackground(opts app.Options) tea.Cmd {
	return func() tea.Msg {
		pr, pw := io.Pipe()

		go func() {
			app.Run(opts, pw)
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
