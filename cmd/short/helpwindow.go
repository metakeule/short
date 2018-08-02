package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"sort"
)

type HelpWindow struct {
	*ModalWindow
}

func NewHelpWindow(s *Screen) *HelpWindow {
	return &HelpWindow{ModalWindow: NewModalWindow(s)}
}

var help = map[string]string{
	"F1": "help",
	// "F2":         "rename selected shortcut",
	// "F3":         "edit command of selected shortcut",
	// "....":         "edit defaults of selected shortcut",
	// "F5":         "discard current parameters",
	"F4": "set parameters as defauls",
	"F6": "add shortcut",
	// "F7":         "add shortcut based on shell history",
	// "F8":         "add shortcut based on selected shortcut and current parameters",
	// "F9":         "save copy of selected shortcut as...",
	// "F10":        "run non-interactive shell command",
	"DEL":        "remove selected shortcut",
	"UP/DOWN":    "select shortcut",
	"LEFT/RIGHT": "previous/next page",
	"ENTER":      "execute command",
	"ESC":        "clear search bar",
	"TAB":        "edit parameters",
	"CTRL-SPACE": "clear parameter buffer",
	"CTRL-C":     "quit",
}

func helpSorted() (h []string) {

	for k, v := range help {
		h = append(h, fmt.Sprintf("%-20v%s", k, v))
	}

	sort.Strings(h)
	return
}

func (h *HelpWindow) Print() {
	h.s.Clear()
	h.s.HideCursor()
	h.s.puts(tcell.StyleDefault, 10, 2, "HELP")
	hp := helpSorted()

	for n, line := range hp {
		h.s.puts(tcell.StyleDefault, 10, n+5, line)
	}

	h.s.puts(tcell.StyleDefault, 1, h.s.height-1, "press ENTER to leave screen")
	h.s.Show()
}

func (h *HelpWindow) KeyEnter(ev *tcell.EventKey) (quit bool) {
	h.s.switchWindow(NewMainWindow(h.s))
	return false
}
