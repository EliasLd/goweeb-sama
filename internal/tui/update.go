package tui

import (
	tea "github.com/charmbracelet/bubbletea"
)

func Update(msg tea.Msg, m Model) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil
	case tea.KeyMsg:
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
			if m.Cursor < 5 {
				m.Cursor++
			}
			m = updateFocus(m)
			return m, nil

		case "enter", " ":
			switch m.Cursor {
			case 1: // AllCheckbox
				m.AllCheckbox.Toggle()
			case 4: // KeepCheckbox
				m.KeepCheckbox.Toggle()
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
		}
	}

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
	}
	return m
}
