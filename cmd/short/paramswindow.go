package main

import (
	"fmt"
	"github.com/metakeule/pager"

	// "fmt"
	"github.com/gdamore/tcell"
	"github.com/metakeule/short"
)

type ParamsWindow struct {
	*ModalWindow
	selectedParam int
	finalCmd      string
	params        params
	finalDefaults map[string]string
	defaults      map[string]string
	cmd           string
	cursorX       int
	shortcut      string
	origCMD       string
	origShortcut  string
	selected      int
	mainWindow    *MainWindow
}

func (p *ParamsWindow) refresh() {
	// cs := p.s.pagedCuts()

	c := p.s.filteredCuts[p.selected]

	p.finalCmd, p.params, p.finalDefaults = p.s.paramLines(c.Name)
	/*
		p.finalCmd, p.params, p.finalDefaults = p.s.paramLines(cs)
			p.origCMD = cs[p.s.Selected].Command
			p.origShortcut = cs[p.s.Selected].Name
			p.defaults = cs[p.s.Selected].Defaults
	*/
	p.origCMD = c.Command
	p.origShortcut = c.Name
	p.defaults = c.Defaults
	p.cmd = p.origCMD
	p.shortcut = p.origShortcut
	p.cursorX = -1
}

func NewParamsWindow(m *MainWindow, selected int) *ParamsWindow {
	p := &ParamsWindow{ModalWindow: NewModalWindow(m.s)}
	p.selected = selected
	p.mainWindow = m
	p.refresh()
	return p
}

func (p *ParamsWindow) insertRuneToShortcut(r rune) {
	var f Field
	f.Value = p.shortcut
	f.Cursor = p.cursorX
	f.Insert(r)
	p.shortcut = f.Value
	p.cursorX = f.Cursor
	p.Print()
}

func (p *ParamsWindow) insertRuneToCMD(r rune) {
	var f Field
	f.Value = p.cmd
	f.Cursor = p.cursorX
	f.Insert(r)
	p.cmd = f.Value
	p.cursorX = f.Cursor
	p.Print()
}

func (p *ParamsWindow) KeyOther(ev *tcell.EventKey) (quit bool) {

	switch p.selectedParam {
	case -2:
		p.insertRuneToShortcut(ev.Rune())
	case -1:
		p.insertRuneToCMD(ev.Rune())
	default:
		p.insertRuneToParam(ev.Rune())
	}

	return
}

func (p *ParamsWindow) KeyCtrlS(ev *tcell.EventKey) (quit bool) {
	name := p.shortcut
	err := p.save()

	if err != nil {
		p.s.Clear()
		p.s.puts(tcell.StyleDefault, 10, 10, fmt.Sprintf("ERROR: can't save: %s", err.Error()))
		p.s.Show()
		return
	}

	err = p.s.loadShortCuts()
	if err != nil {
		p.s.Clear()
		p.s.puts(tcell.StyleDefault, 10, 10, fmt.Sprintf("ERROR: can't save: %s", err.Error()))
		p.s.Show()
		return
	}

	var selected int = -1

	//for idx, n := range p.s.pagedCuts() {
	for idx, n := range p.s.filteredCuts {
		if n.Name == name {
			selected = idx
		}
	}

	if selected == -1 {
		panic("must not happen: didn't find shortcut " + name)
	}

	p.mainWindow.pager = pager.New(p.s.height-3, len(p.s.filteredCuts), pager.PreSelect(uint(selected)))
	p.selected = selected
	p.refresh()
	p.Print()
	return
}

func (p *ParamsWindow) setParamsAsDefaults() error {
	if p.origShortcut != p.shortcut || p.cmd != p.origCMD {
		return fmt.Errorf("you have unsaved changed, save first")
	}

	if _, has := p.s.allCuts[p.shortcut]; has {
		c := p.s.allCuts[p.shortcut]
		c.Defaults = map[string]string{}
		for _, l := range p.params {
			vl, has := p.s.currentParameters[l[0]]
			if has {
				if len(vl) > 0 && vl[0] == '#' {
					vl = vl[1:]
				}
				c.Defaults[l[0]] = vl
			}

		}

		p.s.allCuts[p.shortcut] = c
		_, err := short.Command(p.shortcut, p.s.allCuts, nil)
		if err != nil {
			return err
		}
		return save(p.s.allCuts)

	} else {
		return fmt.Errorf("you have unsaved changed, save first")
	}

}

func (p *ParamsWindow) KeyCtrlF(ev *tcell.EventKey) (quit bool) {
	switch p.selectedParam {
	case -2:
		p.s.bark()
	case -1:
		p.s.bark()
	default:
		//	p.insertRuneToParam(ev.Rune())

		// pName := p.params[p.selectedParam][0]
		// p.s.currentParameters[pName] = value
		val := p.params[p.selectedParam][1]
		if len(val) > 0 && val[0] == '#' {
			val = val[1:]
		}
		p.s.switchWindow(NewFileWindow(p, p.params[p.selectedParam][0], val, p.selected))
	}
	return
}

func (p *ParamsWindow) KeyF4(ev *tcell.EventKey) (quit bool) {
	name := p.shortcut

	err := p.setParamsAsDefaults()

	if err != nil {
		p.s.Clear()
		p.s.puts(tcell.StyleDefault, 10, 10, fmt.Sprintf("ERROR: can't save: %s", err.Error()))
		p.s.Show()
		return
	}

	err = p.s.loadShortCuts()
	if err != nil {
		p.s.Clear()
		p.s.puts(tcell.StyleDefault, 10, 10, fmt.Sprintf("ERROR: can't save: %s", err.Error()))
		p.s.Show()
		return
	}

	var selected int = -1

	//for idx, n := range p.s.pagedCuts() {
	for idx, n := range p.s.filteredCuts {
		if n.Name == name {
			selected = idx
		}
	}

	if selected == -1 {
		panic("must not happen: didn't find shortcut " + name)
	}

	p.mainWindow.pager = pager.New(p.s.height-3, len(p.s.filteredCuts), pager.PreSelect(uint(selected)))
	p.selected = selected

	p.refresh()
	p.Print()
	return
}

func (p *ParamsWindow) insertRuneToParam(r rune) {
	pName := p.params[p.selectedParam][0]

	if _, has := p.s.currentParameters[pName]; !has {
		p.copyDefaultToCurrentParam()
	}

	var f Field
	f.Value = p.s.currentParameters[pName]
	f.Cursor = p.cursorX
	f.Insert(r)
	p.s.currentParameters[pName] = f.Value
	p.cursorX = f.Cursor
	p.Print()

}

func (p *ParamsWindow) backspaceParams() {
	pName := p.params[p.selectedParam][0]

	if _, has := p.s.currentParameters[pName]; !has {
		p.copyDefaultToCurrentParam()
		//p.s.currentParameters[pName] = p.params[p.selectedParam][1]
	}

	if p.s.currentParameters[pName] == "" {
		delete(p.s.currentParameters, pName)
		p.cursorX = -1
		p.Print()
		return
	}

	var f Field
	f.Value = p.s.currentParameters[pName]
	f.Cursor = p.cursorX

	f.Delete()
	p.s.currentParameters[pName] = f.Value
	p.cursorX = f.Cursor
	p.Print()
	return
}

func (p *ParamsWindow) backspaceCMD() {
	var f Field
	f.Value = p.cmd
	f.Cursor = p.cursorX
	f.Delete()
	p.cmd = f.Value
	p.cursorX = f.Cursor
	p.Print()
}

func (p *ParamsWindow) backspaceName() {
	var f Field
	f.Value = p.shortcut
	f.Cursor = p.cursorX
	f.Delete()
	p.shortcut = f.Value
	p.cursorX = f.Cursor
	p.Print()
}

func (p *ParamsWindow) KeyBackspace(ev *tcell.EventKey) (quit bool) {

	switch p.selectedParam {
	case -2:
		p.backspaceName()
	case -1:
		p.backspaceCMD()
	default:
		p.backspaceParams()
	}

	return
}

func (p *ParamsWindow) KeyF2(ev *tcell.EventKey) (quit bool) {
	p.s.CopyAllDefaultsToCurrentParams(p.s.filteredCuts[p.selected].Name)
	p.cursorX = -1
	p.Print()
	return
}

func (p *ParamsWindow) KeyCtrlP(ev *tcell.EventKey) (quit bool) {
	p.copyDefaultToCurrentParam()
	p.cursorX = -1
	p.Print()
	return
}

func (p *ParamsWindow) copyDefaultToCurrentParam() {
	pName := p.params[p.selectedParam][0]
	p.s.currentParameters[pName] = p.finalDefaults[pName]
}

func (p *ParamsWindow) KeyEscape(ev *tcell.EventKey) (quit bool) {
	p.s.switchWindow(p.mainWindow)
	return false
}

func (p *ParamsWindow) KeyLeft(ev *tcell.EventKey) (quit bool) {
	if p.cursorX > 0 {
		p.cursorX--
		p.Print()
		return
	}
	p.s.bark()
	return
}

func (p *ParamsWindow) KeyEnter(ev *tcell.EventKey) (quit bool) {
	return p.mainWindow.KeyEnter(ev)
}

func (p *ParamsWindow) KeyRight(ev *tcell.EventKey) (quit bool) {
	if p.cursorX < 0 {
		p.s.bark()
		return
	}
	// if p.cursorX > 0 {
	if p.selectedParam > -1 {
		val := p.s.currentParameters[p.params[p.selectedParam][0]]
		if p.cursorX >= len(val) {
			p.cursorX = len(val)
		} else {
			p.cursorX++
		}
	}

	if p.selectedParam == -1 {
		if p.cursorX >= len(p.cmd) {
			p.cursorX = len(p.cmd)
		} else {
			p.cursorX++
		}
	}

	if p.selectedParam == -2 {
		if p.cursorX >= len(p.shortcut) {
			p.cursorX = len(p.shortcut)
		} else {
			p.cursorX++
		}
	}
	p.Print()
	return
	// }
	// p.s.bark()
	// return
}

func (p *ParamsWindow) KeyUp(ev *tcell.EventKey) (quit bool) {
	if p.selectedParam > -2 {
		p.selectedParam--
		p.cursorX = -1
		p.Print()
		return
	}
	p.s.bark()
	return
}

func (p *ParamsWindow) KeyDown(ev *tcell.EventKey) (quit bool) {
	if p.selectedParam < len(p.params)-1 {
		p.selectedParam++
		p.cursorX = -1
		p.Print()
		return
	}
	p.s.bark()
	return
}

func (p *ParamsWindow) KeyCtrlSpace(ev *tcell.EventKey) (quit bool) {
	p.cursorX = -1
	return p.s.doDeleteParams(ev)
}

func (p *ParamsWindow) KeyTab(ev *tcell.EventKey) (quit bool) {
	return p.KeyDown(ev)
}

func (p *ParamsWindow) KeyBackTab(ev *tcell.EventKey) (quit bool) {
	return p.KeyUp(ev)
}

func (p *ParamsWindow) printCachedParams() {
	if len(p.s.currentParameters) > 0 {
		for x := 0; x < p.s.width; x++ {
			p.s.puts(p.s.style.search, x, p.s.height-3, " ")
		}
		p.s.puts(p.s.style.search, 0, p.s.height-3, " "+mapToString(p.s.currentParameters)+" ")
	}
}

func (p *ParamsWindow) printHelp() {
	p.s.puts(tcell.StyleDefault, 1, p.s.height-1, "press ENTER to execute Commad, press ESC to return")
	p.s.puts(p.s.style.name, 1, 8, "press F2 to copy all defaults to your parameter buffer")
	p.s.puts(p.s.style.name, 1, 9, "press CTRL+P to copy default value to your parameter buffer")
	p.s.puts(p.s.style.name, 1, 10, "press F4 to safe the current parameter buffer as default values")
}

func (p *ParamsWindow) printShortcut() {
	p.s.puts(p.s.style.name, 1, 1, "Name:")
	style := p.s.style.highlighted
	if p.shortcut != p.origShortcut {
		style = p.s.style.selected
	}
	p.s.puts(style, 50, 1, p.shortcut)
	if p.selectedParam == -2 {
		if p.cursorX < 0 {
			p.cursorX = len(p.shortcut)
		}
		p.s.Screen.ShowCursor(50+p.cursorX, 1)
	}

}

func (p *ParamsWindow) printCmd() {
	p.s.puts(p.s.style.name, 1, 2, "Command:")
	style := p.s.style.highlighted
	if p.cmd != p.origCMD {
		style = p.s.style.selected
	}
	p.s.puts(style, 50, 2, p.cmd)
	if p.selectedParam == -1 {
		if p.cursorX < 0 {
			p.cursorX = len(p.cmd)
		}
		p.s.Screen.ShowCursor(50+p.cursorX, 2)
	}
}

func (p *ParamsWindow) printFinalCmd() {
	p.s.puts(p.s.style.name, 1, 5, "final Command:")
	style := p.s.style.highlighted
	p.s.puts(style, 50, 5, p.finalCmd)
}

func (p *ParamsWindow) printCmdInfo() {
	p.printShortcut()
	p.printCmd()
	p.s.puts(p.s.style.name, 1, 3, "Defaults:")
	p.s.puts(p.s.style.highlighted, 50, 3, mapToString(p.defaults))

	p.printFinalCmd()
	p.s.puts(p.s.style.name, 1, 6, "final Defaults:")
	p.s.puts(p.s.style.highlighted, 50, 6, mapToString(p.finalDefaults))
}

func (p *ParamsWindow) printParams() {
	for i, l := range p.params {
		p.printParam(i, l[0], l[1], l[2])
	}
}

func (p *ParamsWindow) printParam(line int, name, val, typ string) {
	y := line + 15
	p.s.puts(p.s.style.name, 1, y, name+" ("+typ+"):")
	// vl := val

	style := p.s.style.selected

	if len(val) > 0 && val[0] == '#' {
		style = p.s.style.highlighted
		val = val[1:]
	}

	if line == p.selectedParam {
		if p.cursorX < 0 {
			p.cursorX = len(val)
		}

		p.s.Screen.ShowCursor(50+p.cursorX, y)
	}

	if val == "" {
		val = "            "
	}

	p.s.puts(style, 50, y, val)
}

func (p *ParamsWindow) Print() {
	/*
		cs := p.s.pagedCuts()
		if len(cs) == 0 {
			p.s.bark()
			return
		}
	*/

	//p.finalCmd, p.params, p.finalDefaults = p.s.paramLines(cs)
	p.finalCmd, p.params, p.finalDefaults = p.s.paramLines(p.s.filteredCuts[p.selected].Name)

	p.s.Screen.Clear()
	p.s.Screen.HideCursor()

	p.printCmdInfo()
	p.printParams()
	//p.printCachedParams()

	p.mainWindow.printCMD(p.s.filteredCuts[p.selected].Name)

	p.printHelp()

	p.s.Screen.Show()
}

func (p *ParamsWindow) save() error {
	if _, has := p.s.allCuts[p.shortcut]; !has {
		var c short.Cut
		c.Name = p.shortcut
		c.Command = p.cmd
		c.Defaults = map[string]string{}

		p.s.allCuts[c.Name] = c

		err := short.Add(c, p.s.allCuts)

		if err != nil {
			return err
		}

	} else {
		c := p.s.allCuts[p.shortcut]
		c.Command = p.cmd
		p.s.allCuts[p.shortcut] = c
		_, err := short.Command(p.shortcut, p.s.allCuts, nil)
		// p.s.allCuts[p.shortcut].
		if err != nil {
			return err
		}
	}
	return save(p.s.allCuts)

}
