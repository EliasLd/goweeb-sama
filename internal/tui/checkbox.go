package tui

import (
	"fmt"
)

type Checkbox struct {
	Label   string
	Checked bool
}

// Toggle checkbox state
func (c *Checkbox) Toggle() {
	c.Checked = !c.Checked
}

// Returns a string corresponding to the checkbox
func (c Checkbox) View(focused bool) string {
	cursor := "  "
	if focused {
		cursor = "> "
	}

	checked := "[ ]"
	if c.Checked {
		checked = "[x]"
	}

	return fmt.Sprintf("%s%s %s", cursor, checked, c.Label)
}

