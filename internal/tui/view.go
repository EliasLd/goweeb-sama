package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")).Align(lipgloss.Center)
	labelStyle  = lipgloss.NewStyle().Bold(true)
	logBoxStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1).
			Width(70).
			Height(20)
)

func View(m Model) string {
	// Show selection screen if active
	if m.State == StateMangaSelection || m.State == StateScanSelection {
		return m.SelectionModel.View()
	}

	// Show form otherwise
	return viewForm(m)
}

func viewForm(m Model) string {
	var form strings.Builder

	form.WriteString(titleStyle.Render(m.Title))
	form.WriteString("\n\n")

	form.WriteString(labelStyle.Render("Nom du manga"))
	form.WriteString("\n\n")
	if m.Cursor == 0 {
		form.WriteString(m.MangaInput.View())
	} else {
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.MangaInput.View()))
	}
	form.WriteString("\n\n")

	form.WriteString(m.AllCheckbox.View(m.Cursor == 1))
	form.WriteString("\n\n")

	form.WriteString(labelStyle.Render("Plage de chapitres à télécharger"))
	form.WriteString("\n\n")
	if m.AllCheckbox.Checked {
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.RangeInput.View()))
		form.WriteString("\n")
		form.WriteString(lipgloss.NewStyle().Faint(true).Italic(true).Render("Désactivé car vous avez choisi de tout télécharger."))
	} else if m.Cursor == 2 {
		form.WriteString(m.RangeInput.View())
	} else {
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.RangeInput.View()))
	}
	form.WriteString("\n\n")

	form.WriteString(labelStyle.Render("Dossier de destination (si le chemin n'existe pas, il sera créé)"))
	form.WriteString("\n\n")
	if m.Cursor == 3 {
		form.WriteString(m.ScanDirInput.View())
	} else {
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.ScanDirInput.View()))
	}
	form.WriteString("\n\n")

	form.WriteString(labelStyle.Render("Domaine anime-sama (optionnel)"))
	form.WriteString("\n\n")
	if m.Cursor == 4 {
		form.WriteString(m.DomainInput.View())
	} else {
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.DomainInput.View()))
	}
	form.WriteString("\n\n")

	form.WriteString(m.EbookCheckbox.View(m.Cursor == 5))
	form.WriteString("\n\n")

	form.WriteString(m.KeepCheckbox.View(m.Cursor == 6))
	form.WriteString("\n\n")

	if m.DownloadReady {
		button := "[ Télécharger ]"
		if m.Cursor == 7 {
			button = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")).Render(button)
		} else {
			button = lipgloss.NewStyle().Faint(true).Render(button)
		}
		form.WriteString(button)
		form.WriteString("\n\n")
	}

	footerStyle := lipgloss.NewStyle().Faint(true)
	form.WriteString(footerStyle.Render("↑/↓ pour naviguer, espace ou entrée pour cocher, Ctrl+C ou Esc pour quitter"))

	var logs strings.Builder
	const maxLogs = 18
	start := 0
	if len(m.Logs) > maxLogs {
		start = len(m.Logs) - maxLogs
	}

	visibleLogs := m.Logs[start:]
	for _, line := range visibleLogs {
		logs.WriteString(line + "\n")
	}

	if len(m.Logs) == 0 {
		logs.WriteString("Aucun log pour le moment...")
	}

	for i := len(visibleLogs); i < maxLogs; i++ {
		logs.WriteString("\n")
	}

	logsView := logBoxStyle.Render(logs.String())

	content := lipgloss.JoinHorizontal(lipgloss.Top, form.String(), logsView)

	boxWidth := lipgloss.Width(content)
	boxHeight := lipgloss.Height(content)
	horizontalMargin := max(0, (m.Width-boxWidth)/2)
	verticalMargin := max(0, (m.Height-boxHeight)/2)

	boxStyle := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin)

	return boxStyle.Render(content)
}
