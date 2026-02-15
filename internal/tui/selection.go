package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Represents a generic selectable item
type SelectionItem struct {
	Label string
	Value string
}

// Represents a generic selection screen
type SelectionModel struct {
	Title     string
	Items     []SelectionItem
	Cursor    int
	Selected  string
	Cancelled bool
	Width     int
	Height    int
}

func NewSelectionModel(title string, items []SelectionItem) SelectionModel {
	return SelectionModel{
		Title:     title,
		Items:     items,
		Cursor:    0,
		Selected:  "",
		Cancelled: false,
	}
}

func (m SelectionModel) Init() tea.Cmd {
	return nil
}

func (m SelectionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Cancelled = true
			return m, tea.Quit

		case "esc":
			m.Cancelled = true
			return m, nil

		case "up", "k":
			if m.Cursor > 0 {
				m.Cursor--
			}
			return m, nil

		case "down", "j":
			if m.Cursor < len(m.Items)-1 {
				m.Cursor++
			}
			return m, nil

		case "enter", " ":
			if len(m.Items) > 0 {
				m.Selected = m.Items[m.Cursor].Value
			}
			return m, nil
		}
	}

	return m, nil
}

func (m SelectionModel) View() string {
	var s string

	titleStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("226")).
		Padding(1, 0)

	s += titleStyle.Render(m.Title) + "\n\n"

	// Render items
	for i, item := range m.Items {
		cursor := "  "
		if m.Cursor == i {
			cursor = "> "
		}

		itemStyle := lipgloss.NewStyle()
		if m.Cursor == i {
			itemStyle = itemStyle.Bold(true).Foreground(lipgloss.Color("226"))
		} else {
			itemStyle = itemStyle.Faint(true)
		}

		s += fmt.Sprintf("%s%s\n", cursor, itemStyle.Render(item.Label))
	}

	s += "\n"
	footerStyle := lipgloss.NewStyle().Faint(true)
	s += footerStyle.Render("↑/↓ pour naviguer • Entrée pour sélectionner • Esc pour annuler")

	// Center content
	boxWidth := lipgloss.Width(s)
	boxHeight := lipgloss.Height(s)
	horizontalMargin := max(0, (m.Width-boxWidth)/2)
	verticalMargin := max(0, (m.Height-boxHeight)/2)

	boxStyle := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin)

	return boxStyle.Render(s)
}
