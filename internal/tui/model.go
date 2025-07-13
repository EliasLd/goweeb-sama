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


