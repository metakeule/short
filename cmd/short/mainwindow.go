package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/metakeule/pager"
	"github.com/metakeule/short"
	"os"
)

type MainWindow struct {
	*ModalWindow
	pager pager.Pager
}

func NewMainWindow(s *Screen) *MainWindow {
	m := &MainWindow{ModalWindow: NewModalWindow(s)}
	m.resetPager()
	return m
}

func (m *MainWindow) KeyF1(ev *tcell.EventKey) (quit bool) {
	m.s.switchWindow(NewHelpWindow(m.s))
	return
}

func (m *MainWindow) KeyCtrlP(ev *tcell.EventKey) (quit bool) {
	return false
}

func (m *MainWindow) KeyCtrlSpace(ev *tcell.EventKey) (quit bool) {
	return m.s.doDeleteParams(ev)
}

func (m *MainWindow) del() error {
	from, _, selected := m.pager.Indexes()
	//if len(m.s.All) > m.s.Selected {
	if selected > -1 {
		c := m.s.filteredCuts[selected+from]
		delete(m.s.allCuts, c.Name)

		err := save(m.s.allCuts)

		if err != nil {
			return err
		}

		return m.s.loadShortCuts()
	}

	return fmt.Errorf("not found")

}

func (m *MainWindow) KeyDEL(ev *tcell.EventKey) (quit bool) {
	err := m.del()

	m.resetPager()

	if err != nil {
		m.s.Clear()
		m.s.puts(tcell.StyleDefault, 10, 10, fmt.Sprintf("ERROR: can't delete: %s", err.Error()))
		m.s.Show()
		return
	}

	m.Print()
	return
}

func (m *MainWindow) KeyLeft(ev *tcell.EventKey) (quit bool) {
	if !m.pager.PageUp() {
		m.s.bark()
		return
	}
	m.Print()
	// m.Up()
	return
}

func (m *MainWindow) KeyRight(ev *tcell.EventKey) (quit bool) {
	if !m.pager.PageDown() {
		m.s.bark()
		return
	}
	m.Print()
	// m.Up()
	return
}

func (m *MainWindow) KeyUp(ev *tcell.EventKey) (quit bool) {
	if !m.pager.Prev() {
		m.s.bark()
		return
	}
	m.Print()
	// m.Up()
	return
}

func (m *MainWindow) KeyDown(ev *tcell.EventKey) (quit bool) {
	// m.Down()
	if !m.pager.Next() {
		m.s.bark()
		return
	}
	m.Print()
	return
}

func (m *MainWindow) KeyEnter(ev *tcell.EventKey) (quit bool) {
	from, _, selected := m.pager.Indexes()

	if selected == -1 {
		m.s.bark()
		return
	}

	m.s.switchWindow(NewParamsWindow(m, from+selected))
	return
}

func (m *MainWindow) KeyEscape(ev *tcell.EventKey) (quit bool) {
	m.s.resetFilter()
	m.resetPager()
	m.Print()
	return
}

func (m *MainWindow) KeyCtrlE(ev *tcell.EventKey) (quit bool) {
	_, _, selected := m.pager.Indexes()

	//if len(m.s.pagedCuts()) == 0 {
	if selected == -1 {
		m.s.bark()
	} else {

		// TODO: query params
		name := m.s.filteredCuts[selected].Name

		//cmd, err := short.Command(m.s.SelectedName(), m.s.allCuts, m.s.currentParameters)
		cmd, err := short.Command(name, m.s.allCuts, m.s.currentParameters)
		if err != nil {
			m.s.Clear()
			fmt.Fprintf(os.Stderr, "Error: %s", err.Error())
		} else {
			m.s.Fini()
			m.s.finished = true
			fmt.Fprintf(os.Stdout, "\n###################### running ######################\n%s\n#####################################################\n\n\n", cmd)
			//short.Exec(m.s.SelectedName(), m.s.allCuts, m.s.currentParameters)
			short.Exec(name, m.s.allCuts, m.s.currentParameters)
			return true
		}
	}

	return false
}

func (m *MainWindow) resetPager() {
	m.pager = pager.New(m.s.height-3, len(m.s.filteredCuts), 0)
}

func (m *MainWindow) KeyBackspace(ev *tcell.EventKey) (quit bool) {

	rs := []rune(m.s.Search)
	if len(rs) > 0 {
		rs = rs[0 : len(rs)-1]
		m.s.Search = string(rs)
	}

	// m.s.Selected = 0
	if m.s.Search != "" {
		m.s.fuzzyFind()
	} else {
		m.s.resetFilter()
	}
	m.resetPager()

	m.Print()
	return
}

func (m *MainWindow) Print() {
	//cs := m.s.pagedCuts()

	m.s.Screen.Clear()

	from, to, selected := m.pager.Indexes()
	var cs short.Cuts

	if from > -1 {
		cs = m.s.filteredCuts[from:to]
	}

	for i, c := range cs {
		styleName := m.s.style.name
		styleCode := m.s.style.code
		//if i == m.s.Selected {
		if i == selected {
			styleName = m.s.style.highlighted.Bold(true)
			styleCode = m.s.style.highlighted

			for j := 0; j < m.s.width; j++ {
				m.s.puts(styleCode, j, i, " ")
			}
		}
		m.s.puts(styleName, 1, i, c.Name)
		m.s.puts(styleCode, 50, i, c.Command)

		defaults := mapToString(c.Defaults)

		if defaults != "" {
			m.s.puts(styleCode, 120, i, defaults)
		}
	}

	for j := 0; j < m.s.width; j++ {
		m.s.puts(m.s.style.selected, j, m.s.height-2, " ")
	}

	//if m.s.Selected > len(cs)-1 {
	// m.s.puts(m.s.style.selected, 1, m.s.height-2, " ")
	// } else {
	if selected > -1 {
		//cmd, err := short.Command(cs[m.s.Selected].Name, m.s.allCuts, m.s.currentParameters)

		cmd, err := short.Command(cs[selected].Name, m.s.allCuts, m.s.currentParameters)
		if err != nil {
			cmd = "ERROR: " + err.Error()
		}
		//m.s.puts(m.s.style.highlighted, 1, m.s.height-3, "custom parameters: ")
		if len(m.s.currentParameters) > 0 {
			for x := 0; x < m.s.width; x++ {
				m.s.puts(m.s.style.search, x, m.s.height-3, " ")
			}
			m.s.puts(m.s.style.search, 0, m.s.height-3, " "+mapToString(m.s.currentParameters)+" ")
		}
		m.s.puts(m.s.style.selected, 1, m.s.height-2, cmd)
	}

	pretext := " search (F1 for help, CTRL+e to execute) "
	lenpre := len(pretext)

	m.s.puts(m.s.style.highlighted, 0, m.s.height-1, pretext)

	for j := lenpre; j < m.s.width; j++ {
		m.s.puts(m.s.style.search, j, m.s.height-1, " ")
	}

	m.s.puts(m.s.style.search, lenpre, m.s.height-1, m.s.Search)
	m.s.Screen.ShowCursor(lenpre+len(m.s.Search), m.s.height-1)
	m.s.Screen.Show()
}

/*
func (m *MainWindow) Up() {
	// TODO: page up if we are not on the first page
	if m.s.Selected == 0 {
		m.s.bark()
		return
	}
	m.s.Selected--
	m.Print()
}
*/

/*
func (m *MainWindow) Down() {
	// TODO: page down if we are not on the last page
	cs := m.s.pagedCuts()

	if m.s.Selected < len(cs)-1 {
		m.s.Selected++
		m.Print()
		return
	}
	m.s.bark()

	//s.Print()
}
*/

func (m *MainWindow) KeyOther(ev *tcell.EventKey) (quit bool) {
	m.s.Search += string(ev.Rune())
	m.s.fuzzyFind()
	m.resetPager()
	m.Print()
	return
}
