package tui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var asciiArt string = `
   _____                                                         
  / ___/_________ _____     __________________  _____  ___  _____
  \__ \/ ___/ __ ` + "`" + `/ __ \   / ___/ ___/ ___/ __ ` + "`" + `/ __ \/ _ \/ ___/
 ___/ / /__/ /_/ / / / /  (__  ) /__/ /  / /_/ / /_/ /  __/ /    
/____/\___/\__,_/_/ /_/  /____/\___/_/   \__,_/ .___/\___/_/     
                                             /_/                 
`

type Model struct {
	Title		string
	MangaInput	textinput.Model
	AllCheckbox	Checkbox
	RangeInput	textinput.Model
	ScanDirInput	textinput.Model
	KeepCheckbox	Checkbox
	Cursor		int
	Width		int
	Height		int
}

func InitialModel() Model {
	manga := textinput.New()
	manga.Placeholder = "ex: jujutsu-kaisen"
	manga.Focus()
	manga.Prompt = "> "
	manga.CharLimit = 50
	manga.Width = 60

	rangeInput := textinput.New()
	rangeInput.Placeholder = "ex: 1-28"
	rangeInput.Prompt = "> "
	rangeInput.Width = 60

	scanDir := textinput.New()
	scanDir.Placeholder = "ex: C:\\Users\\<username>\\Documents\\scans\\jjk"
	scanDir.Prompt = "> "
	scanDir.Width = 70

	return Model {
		Title:		asciiArt,
		MangaInput:	manga,
		AllCheckbox:	Checkbox{Label: "Télécharger tous les chapitres.", Checked: false},
		RangeInput:	rangeInput,
		ScanDirInput:	scanDir,
		KeepCheckbox:	Checkbox{Label: "Garder les images après conversion (déconseillé).", Checked: false},
		Cursor:		0,
		Width:		0,
		Height:		0,
	}
}

func (m Model) Init() tea.Cmd {
	return nil 
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return Update(msg, m)
}

func (m Model) View() string {
    	return View(m)
}

