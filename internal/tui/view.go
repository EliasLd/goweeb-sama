package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

func View(m Model) string {
	var b strings.Builder

	titleStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("226")).Align(lipgloss.Center)
	labelStyle := lipgloss.NewStyle().Bold(true)

	b.WriteString(titleStyle.Render(m.Title))
	b.WriteString("\n\n")
	
	b.WriteString(labelStyle.Render("Nom du manga (espace = -)"))
	b.WriteString("\n\n")
	if m.Cursor == 0 {
		b.WriteString(m.MangaInput.View())
	} else {
		//Faint field when unfocused
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(m.MangaInput.View()))
	}
	b.WriteString("\n\n")

	
	b.WriteString(m.AllCheckbox.View(m.Cursor == 1))
	b.WriteString("\n\n")
	
	b.WriteString(labelStyle.Render("Plage de chapitres à télécharger"))
	b.WriteString("\n\n")
	// Not active if AllCheckbox is checked
	if m.AllCheckbox.Checked {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(m.RangeInput.View()))
		b.WriteString("\n")
		b.WriteString(lipgloss.NewStyle().Faint(true).Italic(true).Render("Désactivé car vous avez choisi de tout télécharger."))
	} else if m.Cursor == 2 {
		b.WriteString(m.RangeInput.View())
	} else {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(m.RangeInput.View()))
	}
	b.WriteString("\n\n")

	b.WriteString(labelStyle.Render("Dossier de destination (si le chemin n'existe pas, il sera créé)"))
	b.WriteString("\n\n")
	if m.Cursor == 3 {
		b.WriteString(m.ScanDirInput.View())
	} else {
		b.WriteString(lipgloss.NewStyle().Faint(true).Render(m.ScanDirInput.View()))
	}
	b.WriteString("\n\n")

	b.WriteString(m.KeepCheckbox.View(m.Cursor == 4))
	b.WriteString("\n\n")

	// Footer
	footerStyle := lipgloss.NewStyle().Faint(true)
	b.WriteString(footerStyle.Render("↑/↓ pour naviguer, espace ou entrée pour cocher, Ctrl+C pour quitter"))

	ui := b.String()

	boxWidth := 80  
	boxHeight := 25

	horizontalMargin := max(0, (m.Width-boxWidth)/2)
	verticalMargin := max(0, (m.Height-boxHeight)/2)

	boxStyle := lipgloss.NewStyle().
		MarginTop(verticalMargin).
		MarginLeft(horizontalMargin)

	return boxStyle.Render(ui)

}

