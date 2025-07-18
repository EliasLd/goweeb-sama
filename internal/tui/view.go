package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")).Align(lipgloss.Center)
	labelStyle = lipgloss.NewStyle().Bold(true)
	logBoxStyle = lipgloss.NewStyle().
		Border(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1).
		Width(50).
		Height(20)
)


func View(m Model) string {
	var form strings.Builder

	form.WriteString(titleStyle.Render(m.Title))
	form.WriteString("\n\n")
	
	form.WriteString(labelStyle.Render("Nom du manga (espace = -)"))
	form.WriteString("\n\n")
	if m.Cursor == 0 {
		form.WriteString(m.MangaInput.View())
	} else {
		//Faint field when unfocused
		form.WriteString(lipgloss.NewStyle().Faint(true).Render(m.MangaInput.View()))
	}
	form.WriteString("\n\n")

	
	form.WriteString(m.AllCheckbox.View(m.Cursor == 1))
	form.WriteString("\n\n")
	
	form.WriteString(labelStyle.Render("Plage de chapitres à télécharger"))
	form.WriteString("\n\n")
	// Not active if AllCheckbox is checked
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

	form.WriteString(m.KeepCheckbox.View(m.Cursor == 4))
	form.WriteString("\n\n")

	if m.DownloadReady {
		button := "[ Télécharger ]"
		if m.Cursor == 5 {
			button = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")).Render(button)
		} else {
			button = lipgloss.NewStyle().Faint(true).Render(button)
		}
		form.WriteString(button)
		form.WriteString("\n\n")
	}

	// Footer
	footerStyle := lipgloss.NewStyle().Faint(true)
	form.WriteString(footerStyle.Render("↑/↓ pour naviguer, espace ou entrée pour cocher, Ctrl+C pour quitter"))

	var logs strings.Builder

	const maxLogs = 30
	start := 0
	if len(m.Logs) > maxLogs {
		start = len(m.Logs) - maxLogs
	}

	if len(m.Logs) == 0 {
		logs.WriteString("Aucun log pour le moment...")
	} else {
		for _, line := range m.Logs[start:] {
			logs.WriteString(line + "\n")
		}
	}

	// Adapt logBox height
	height := m.Height - 5
	if height < 10 {
		height = 10
	}

	logBoxStyleDynamic := logBoxStyle.Copy().Height(height)

	logsView := logBoxStyleDynamic.Render(logs.String())

	content := lipgloss.JoinHorizontal(lipgloss.Top, form.String(), logsView)

	// Center content
	boxWidth := lipgloss.Width(content)
	boxHeight := lipgloss.Height(content)
	horizontalMargin := max(0, (m.Width - boxWidth)/2)
	verticalMargin := max(0, (m.Height - boxHeight)/2)

	boxStyle := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin)
	
	return boxStyle.Render(content)
}

