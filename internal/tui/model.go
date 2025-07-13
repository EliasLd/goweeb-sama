package tui

import (
	"github.com/charmbracelet/bubbles/checkbox"
	"github.com/charmbracelet/bubbles/textinput"
)

type Model struct {
	Title		string
	MangaInput	textinput.Model
	AllCheckbox	checkbox.Model
	RangeInput	textinput.Model
	ScanDirInput	textinput.Model
	KeepCheckbox	checkbox.Model
	Cursor		int
}

func InitialModel() Model {
	manga := textinput.New()
	manga.Placeholder = "Entrer le nom du manga (espace = -)"
	manga.Focus()
	manga.Prompt = "> "
	manga.CharLimit = 50

	rangeInput := textinput.New()
	rangeInput.Placeholder = "exemple: 10-77"
	rangeInput.Prompt("> ")

	scanDir := textinput.New()
	scanDir.Placeholder = "Dossier de destination"
	scanDir.Prompt = "> "
	scanDir.SetValue("pdf")

	return Model {
		Title:		"Scan scraper",
		MangaInput:	manga,
		AllCheckBox:	checkbox.New("Télécharger tous les chapitres", false),
		RangeInput:	rangeInput,
		ScanDirInput:	scanDir,
		KeepCheckbox:	checkbox.New("Garder les images après conversion", false),
		Cursor:		0,
	}
}
