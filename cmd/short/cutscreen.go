package main

import (
	"fmt"
	"github.com/gdamore/tcell"
	"github.com/gdamore/tcell/encoding"
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/mattn/go-runewidth"
	"github.com/metakeule/short"
	"sort"
	"time"
)

type Window interface {
	KeyDEL(ev *tcell.EventKey) (quit bool)
	KeyCtrlP(ev *tcell.EventKey) (quit bool)
	KeyCtrlSpace(ev *tcell.EventKey) (quit bool)
	KeyF1(ev *tcell.EventKey) (quit bool)
	KeyUp(ev *tcell.EventKey) (quit bool)
	KeyLeft(ev *tcell.EventKey) (quit bool)
	KeyDown(ev *tcell.EventKey) (quit bool)
	KeyRight(ev *tcell.EventKey) (quit bool)
	KeyTab(ev *tcell.EventKey) (quit bool)
	KeyBackTab(ev *tcell.EventKey) (quit bool)
	KeyEnter(ev *tcell.EventKey) (quit bool)
	KeyEscape(ev *tcell.EventKey) (quit bool)
	KeyCtrlS(ev *tcell.EventKey) (quit bool)
	KeyCtrlC(ev *tcell.EventKey) (quit bool)
	KeyCtrlE(ev *tcell.EventKey) (quit bool)
	KeyCtrlL(ev *tcell.EventKey) (quit bool)
	KeyCtrlR(ev *tcell.EventKey) (quit bool)
	KeyBackspace(ev *tcell.EventKey) (quit bool)
	KeyF4(ev *tcell.EventKey) (quit bool)

	KeyOther(ev *tcell.EventKey) (quit bool)
	Print()
}

type Screen struct {
	tcell.Screen
	All          short.Cuts
	allCuts      map[string]short.Cut
	Lines        int
	First        int
	Selected     int
	filteredCuts short.Cuts
	Search       string
	style        struct {
		name        tcell.Style
		code        tcell.Style
		highlighted tcell.Style
		selected    tcell.Style
		search      tcell.Style
	}

	currentWindow Window

	KeyMap map[tcell.Key]func(*tcell.EventKey) (quit bool)

	currentParameters map[string]string
	// paramsMode        bool

	// modalWin bool
	finished bool
	width    int
	height   int
}

func (s *Screen) loadShortCuts() error {
	allCuts, err := load()
	if err != nil {
		return err
	}

	all := short.Sort(allCuts)

	s.All = all
	s.allCuts = allCuts
	s.First = 0
	s.Selected = 0
	s.filteredCuts = all

	return nil

}

func NewScreen() (*Screen, error) {

	/*
		allCuts, err := load()
		if err != nil {
			return nil, err
		}

		all := short.Sort(allCuts)
	*/
	sc, e := tcell.NewScreen()
	if e != nil {
		return nil, e
	}

	encoding.Register()

	if e = sc.Init(); e != nil {
		return nil, e
	}

	sc.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	sc.Clear()

	s := &Screen{
		//	All:          all,
		//allCuts:      allCuts,
		//		First:        0,
		//	Selected:     0,
		Screen: sc,
		//		filteredCuts: all,

		currentParameters: map[string]string{},
	}

	err := s.loadShortCuts()
	if err != nil {
		sc.Fini()
		return nil, err
	}

	s.style.name = tcell.StyleDefault
	s.style.code = tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorGray)
	s.style.highlighted = tcell.StyleDefault.Background(tcell.ColorBlue).Foreground(tcell.ColorWhite)
	s.style.selected = tcell.StyleDefault.Background(tcell.ColorRed).Foreground(tcell.ColorWhite).Bold(true)
	s.style.search = tcell.StyleDefault.Background(tcell.ColorWhite).Foreground(tcell.ColorBlack).Bold(true)

	s.width, s.height = s.Screen.Size()
	s.KeyMap = map[tcell.Key]func(*tcell.EventKey) (quit bool){}
	s.switchWindow(NewMainWindow(s))
	return s, nil
}

func (s *Screen) Print() {
	s.currentWindow.Print()
}

func (s *Screen) switchWindow(w Window) {
	s.currentWindow = w
	s.setKeyMaps()
	s.Print()
}

func (s *Screen) setKeyMaps() {
	s.KeyMap[271] = s.currentWindow.KeyDEL /* DEL or ENTF key */
	s.KeyMap[tcell.KeyCtrlP] = s.currentWindow.KeyCtrlP
	s.KeyMap[tcell.KeyCtrlSpace] = s.currentWindow.KeyCtrlSpace
	s.KeyMap[tcell.KeyF1] = s.currentWindow.KeyF1
	s.KeyMap[tcell.KeyF2] = nil
	s.KeyMap[tcell.KeyF3] = nil
	s.KeyMap[tcell.KeyF4] = nil
	s.KeyMap[tcell.KeyF5] = nil
	s.KeyMap[tcell.KeyF6] = nil
	s.KeyMap[tcell.KeyF7] = nil
	s.KeyMap[tcell.KeyF8] = nil
	s.KeyMap[tcell.KeyF9] = nil
	s.KeyMap[tcell.KeyF10] = nil
	s.KeyMap[tcell.KeyF11] = s.doBark
	s.KeyMap[tcell.KeyF12] = s.doBark /* nothing to do or to see here */
	s.KeyMap[tcell.KeyUp] = s.currentWindow.KeyUp
	s.KeyMap[tcell.KeyDown] = s.currentWindow.KeyDown
	s.KeyMap[tcell.KeyLeft] = s.currentWindow.KeyLeft
	s.KeyMap[tcell.KeyRight] = s.currentWindow.KeyRight
	s.KeyMap[tcell.KeyTab] = s.currentWindow.KeyTab
	s.KeyMap[tcell.KeyEnter] = s.currentWindow.KeyEnter
	s.KeyMap[tcell.KeyBacktab] = s.currentWindow.KeyBackTab
	s.KeyMap[tcell.KeyEscape] = s.currentWindow.KeyEscape
	s.KeyMap[tcell.KeyCtrlC] = s.currentWindow.KeyCtrlC
	s.KeyMap[tcell.KeyCtrlL] = s.currentWindow.KeyCtrlL
	s.KeyMap[tcell.KeyCtrlR] = s.currentWindow.KeyCtrlR
	s.KeyMap[tcell.KeyBackspace] = s.currentWindow.KeyBackspace
	s.KeyMap[tcell.KeyBackspace2] = s.currentWindow.KeyBackspace
	s.KeyMap[tcell.KeyCtrlE] = s.currentWindow.KeyCtrlE
	s.KeyMap[tcell.KeyCtrlS] = s.currentWindow.KeyCtrlS
	s.KeyMap[tcell.KeyF4] = s.currentWindow.KeyF4
}

func (s *Screen) pagedCuts() short.Cuts {
	cs := s.filteredCuts[s.First:]

	if len(cs) > s.height-3 {
		cs = cs[0 : s.height-3]
	}

	return cs
}

func (s *Screen) paramLines(cs short.Cuts) (cmd string, lines params, defaults map[string]string) {

	paramDefs, err := short.Params(cs[s.Selected].Name, s.allCuts)

	if err != nil {
		s.puts(s.style.code, 10, 10, "ERROR: "+err.Error())
		s.Screen.Show()
		return
	}

	cmd, vals, err2 := short.CommandAndValues(cs[s.Selected].Name, s.allCuts, nil)
	if err2 != nil {
		s.puts(s.style.code, 10, 10, "ERROR: "+err2.Error())
		s.Screen.Show()
		return
	}

	for k, v := range paramDefs {
		vl := vals[k]
		cp, has := s.currentParameters[k]
		if has {
			vl = cp
		}
		if !has || vals[k] == cp {
			vl = "#" + vl
		}
		lines = append(lines, [3]string{k, vl, v})
	}
	// s.currentParameters = vals

	defaults = vals

	if defaults == nil {
		defaults = map[string]string{}
	}

	sort.Sort(lines)
	return
}

func (s *Screen) CopyAllDefaultsToCurrentParams() {
	cs := s.pagedCuts()

	_, vals, err2 := short.CommandAndValues(cs[s.Selected].Name, s.allCuts, nil)
	if err2 != nil {
		s.puts(s.style.code, 10, 10, "ERROR: "+err2.Error())
		s.Screen.Show()
		return
	}

	for k, v := range vals {
		s.currentParameters[k] = v
	}

}

func (s *Screen) debug(str string) {
	s.puts(tcell.StyleDefault, 50, 50, str)
	s.Screen.Show()
	time.Sleep(time.Millisecond * 300)
}

func (s *Screen) resetFilter() {
	s.filteredCuts = s.All
	s.Selected = 0
	s.First = 0
	s.Search = ""
}

func (s *Screen) clearSearchLine(w int, y int, style tcell.Style) {
	for i := 0; i < w; i++ {
		s.puts(style, i, y, " ")
	}
}

func (s *Screen) fuzzyFind() {
	if s.Search == "" {
		return
	}
	var words []string
	var m = map[string]short.Cut{}

	for _, c := range s.All {
		words = append(words, c.Name)
		m[c.Name] = c
	}

	result := fuzzy.RankFind(s.Search, words) // [{whl cartwheel 6 0} {whl wheel 2 2}]
	sort.Sort(result)

	s.filteredCuts = nil

	for _, r := range result {
		s.filteredCuts = append(s.filteredCuts, m[r.Target])
	}
	s.First = 0
	s.Selected = 0
}

func (s *Screen) SelectedName() string {
	cs := s.pagedCuts()

	return cs[s.Selected].Name
}

func (s *Screen) bark() {
	s.Screen.Clear()
	s.Screen.Show()
	time.Sleep(time.Millisecond * 30)
	s.currentWindow.Print()
}

func (s *Screen) doDeleteParams(ev *tcell.EventKey) (quit bool) {
	s.currentParameters = map[string]string{}
	s.currentWindow.Print()
	return
}

func (s *Screen) doNotImplemented(ev *tcell.EventKey) (quit bool) {
	// s.modalWin = true
	s.Clear()
	s.HideCursor()
	s.puts(tcell.StyleDefault, 10, 2, "not implemented yet")
	s.puts(tcell.StyleDefault, 1, s.height-1, "press ENTER to leave screen")
	s.Show()
	return
}

func (s *Screen) doBark(ev *tcell.EventKey) (quit bool) {
	s.bark()
	return
}

func (s *Screen) doQuit(ev *tcell.EventKey) (quit bool) {
	return true
}

func (s *Screen) doSync(ev *tcell.EventKey) (quit bool) {
	s.Screen.Sync()
	s.width, s.height = s.Screen.Size()
	return
}

func (s *Screen) Run() error {

	quit := make(chan struct{})
	s.currentWindow.Print()

	go func() {
		for {
			ev := s.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventResize:
				s.doSync(nil)
			case *tcell.EventKey:

				fn, has := s.KeyMap[ev.Key()]

				if !has {
					fn = s.currentWindow.KeyOther
				}

				if fn != nil {
					q := fn(ev)

					if q {
						close(quit)
					}
				}
				//fn = s.doNotImplemented

			}
		}
	}()

	<-quit

	if !s.finished {
		s.Fini()
	}

	fmt.Println("quitting")

	return nil
}

func (s *Screen) puts(style tcell.Style, x, y int, str string) {
	i := 0
	var deferred []rune
	dwidth := 0
	for _, r := range str {
		switch runewidth.RuneWidth(r) {
		case 0:
			if len(deferred) == 0 {
				deferred = append(deferred, ' ')
				dwidth = 1
			}
		case 1:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 1
		case 2:
			if len(deferred) != 0 {
				s.SetContent(x+i, y, deferred[0], deferred[1:], style)
				i += dwidth
			}
			deferred = nil
			dwidth = 2
		}
		deferred = append(deferred, r)
	}
	if len(deferred) != 0 {
		s.SetContent(x+i, y, deferred[0], deferred[1:], style)
		i += dwidth
	}
}
