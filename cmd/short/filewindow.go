package main

import (
	"github.com/metakeule/pager"
	"io/ioutil"
	"log"
	"sort"
	// "path"
	"path/filepath"
	"strings"
	// "fmt"
	"github.com/gdamore/tcell"
	"github.com/lithammer/fuzzysearch/fuzzy"
	// "sort"
	"os"
	// "sync"
)

type FileWindow struct {
	*ModalWindow
	fileSearch string
	wd         string
	files      []string
	dir        string
	search     string
	// mx         sync.Mutex
	// Selected         int
	searchDirs       bool
	includeHidden    bool
	paramName        string
	selectedShortCut int
	pager            pager.Pager
	paramsWindow     *ParamsWindow
	logger           *log.Logger
	absolutePath     bool
}

func NewFileWindow(p *ParamsWindow, paramName string, value string, selectedShortCut int) *FileWindow {
	f := &FileWindow{ModalWindow: NewModalWindow(p.s)}
	f.paramsWindow = p
	f.s.Screen.HideCursor()
	f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	f.wd, _ = os.Getwd()
	f.includeHidden = false
	f.searchDirs = false
	f.dir = f.wd
	f.paramName = paramName
	f.fileSearch = strings.TrimSpace(value)
	f.selectedShortCut = selectedShortCut
	f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	f.pager = pager.New(p.s.height-3, 0)
	_ = os.Create
	// out, _ := os.Create(filepath.Join(f.wd, "files.log"))

	// f.logger = log.New(out, "short-files", log.Lshortfile)
	if f.fileSearch != "" {
		f.findFiles()
	}
	return f
}

func (f *FileWindow) KeyEnter(ev *tcell.EventKey) (quit bool) {

	f.s.currentParameters[f.paramName] = f.fileSearch
	f.s.switchWindow(f.paramsWindow)
	return false
}

func (f *FileWindow) KeyEscape(ev *tcell.EventKey) (quit bool) {
	f.s.switchWindow(f.paramsWindow)
	return false
}

func (f *FileWindow) KeyBackspace(ev *tcell.EventKey) (quit bool) {

	// f.mx.Lock()
	rs := []rune(f.fileSearch)
	if len(rs) > 0 {
		rs = rs[0 : len(rs)-1]
		f.fileSearch = string(rs)
		f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	}

	f.Print()
	if strings.TrimSpace(f.fileSearch) != "" {
		f.findFiles()
	}
	// f.mx.Unlock()

	return
}

func (f *FileWindow) KeyOther(ev *tcell.EventKey) (quit bool) {
	// f.mx.Lock()
	f.fileSearch += string(ev.Rune())
	// f.logger.Printf("filesearch: %#v\n", f.fileSearch)
	f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	/*
		m.s.fuzzyFind()
		m.Print()
	*/
	f.Print()
	if strings.TrimSpace(f.fileSearch) != "" {

		f.findFiles()
	}
	// f.mx.Unlock()
	return
}

func (f *FileWindow) Print() {
	f.s.Clear()

	/*
		files := f.files

		if len(files) > 40 {
			files = files[:40]
		}
	*/

	f.s.puts(f.s.style.name, 1, 1, "Search: ")
	f.s.puts(f.s.style.highlighted, 9, 1, f.fileSearch)

	abs, _ := filepath.Abs(f.dir)
	f.s.puts(f.s.style.name, 1, 2, "("+abs+")") // +" ["+f.search+"]")

	/*
		if len(files) > f.Selected {
			f.s.puts(f.s.style.highlighted, 1, 2, files[f.Selected])
		}
	*/

	/*
		f.s.puts(f.s.style.highlighted, 1, 3, f.dir)
		f.s.puts(f.s.style.highlighted, 1, 4, f.search)
	*/

	from, to, selected := f.pager.Indexes()

	/*
		f.Selected = selected

		if f.Selected != -1 {
			f.Selected += from
		}
	*/

	if from > -1 {

		for i, file := range f.files[from:to] {
			style := tcell.StyleDefault
			//if i == f.Selected {
			if i == selected {
				style = f.s.style.selected

				for x := 0; x < f.s.width; x++ {
					f.s.puts(style, x, 3+i, " ")
				}

			}
			f.s.puts(style, 5, 3+i, file)
		}
	}

	f.s.Show()
}

func (f *FileWindow) KeyUp(ev *tcell.EventKey) (quit bool) {
	//f.Up()
	if !f.pager.Prev() {
		f.s.bark()
		return
	}

	f.Print()
	return
}

func (f *FileWindow) absolutify() {

}

func (f *FileWindow) KeyF2(ev *tcell.EventKey) (quit bool) {
	if len(f.fileSearch) == 0 {
		f.s.bark()
		return
	}

	var postfix string
	if f.fileSearch[len(f.fileSearch)-1] == '/' {
		postfix = "/"
	}

	f.absolutePath = !f.absolutePath

	if f.absolutePath {
		abs, err := filepath.Abs(f.tildeToHome(f.fileSearch))
		if err == nil {
			f.fileSearch = abs + postfix
		} else {
			f.s.bark()
			return
		}
	} else {
		if filepath.IsAbs(f.fileSearch) {
			rel, err := filepath.Rel(f.wd, f.tildeToHome(f.fileSearch))
			if err == nil {
				f.fileSearch = rel + postfix
			} else {
				f.s.bark()
				return
			}
		}
	}

	f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	f.findFiles()
	f.Print()
	// } else {
	// f.s.bark()
	// }
	return
}

func (f *FileWindow) KeyTab(ev *tcell.EventKey) (quit bool) {
	/*
		files := f.files

		if len(files) > 40 {
			files = files[:40]
		}
	*/

	from, to, selected := f.pager.Indexes()

	if selected == -1 {
		f.s.bark()
		return
	}

	// if len(files) > f.Selected {

	//f.fileSearch = files[f.Selected]
	f.fileSearch = f.files[from:to][selected]

	f.s.Screen.ShowCursor(9+len(f.fileSearch), 1)
	f.findFiles()
	f.Print()
	// } else {
	// f.s.bark()
	// }
	return
}

func (f *FileWindow) KeyLeft(ev *tcell.EventKey) (quit bool) {
	// f.Down()
	if !f.pager.PageUp() {
		f.s.bark()
		return
	}

	f.Print()
	return
}

func (f *FileWindow) KeyRight(ev *tcell.EventKey) (quit bool) {
	// f.Down()
	if !f.pager.PageDown() {
		f.s.bark()
		return
	}

	f.Print()
	return
}

func (f *FileWindow) KeyDown(ev *tcell.EventKey) (quit bool) {
	// f.Down()
	if !f.pager.Next() {
		f.s.bark()
		return
	}

	f.Print()
	return
}

func (f *FileWindow) tildeToHome(path string) string {
	return strings.Replace(path, "~", os.Getenv("HOME"), 1)
}

func (f *FileWindow) homeToTilde(path string) string {
	return strings.Replace(path, os.Getenv("HOME"), "~", 1)
}

/*
func (f *FileWindow) Up() {
	f.mx.Lock()
	if f.Selected == 0 {
		f.s.bark()
		return
	}
	f.Selected--
	f.Print()
	f.mx.Unlock()
}

func (f *FileWindow) Down() {
	f.mx.Lock()
	defer f.mx.Unlock()
	// TODO: page down if we are not on the last page
	if f.Selected < len(f.files)-1 {
		f.Selected++
		f.Print()
		return
	}
	f.s.bark()

	//s.Print()
}
*/

func (f *FileWindow) findFiles() {
	// f.logger.Printf("-------Start---------\n")
	f.files = nil
	// f.Selected = 0

	var abs bool
	var hasTilde bool
	// abs = true

	if len(f.fileSearch) > 0 && (f.fileSearch[0] == '/') {
		f.dir = filepath.Dir(f.fileSearch)
		abs = true
	} else {
		if len(f.fileSearch) > 0 && f.fileSearch[0] == '~' {
			// f.logger.Println("jo")
			hasTilde = true
			f.dir = filepath.Dir(f.tildeToHome(f.fileSearch))
			abs = true
		} else {

			f.dir = filepath.Join(f.wd, filepath.Dir(f.fileSearch))
		}
	}

	// hasTilde := strings.Index(f.dir, "~") == 0

	// home := os.Getenv("HOME")
	//f.dir = filepath.Dir(f.fileSearch)

	if len(f.fileSearch) == 0 {

	} else {

		idx := strings.LastIndex(f.fileSearch, "/")

		if idx == -1 {
			f.search = f.fileSearch
		} else {
			if idx == len(f.fileSearch)-1 {
				f.search = "."
			} else {
				f.search = f.fileSearch[idx+1:]
			}
		}
	}
	/*
		if f.fileSearch[len(f.fileSearch)-1] == '/' {
			f.search = "."
			//f.search = "*"
		} else {
			f.search = filepath.Base(f.fileSearch)
		}
	*/

	//go func() {
	//files, _ := filepath.Glob(f.dir + "/*")
	files, _ := ioutil.ReadDir(f.dir)
	// fuzzy.Find(source, targets)
	// f.logger.Printf("dir: %#v search: %#v\n", f.dir, f.search)

	for _, fl := range files {
		name := fl.Name()

		//if pth[len(pth)-1] != '.' {
		if name != "." && name != ".." {
			// f.logger.Printf("file: %#v\n", fl.Name())
			//if pth[len(pth)-1] != '.' && fuzzy.Match(f.search, filepath.Base(pth)) {
			_ = fuzzy.Match
			//if f.search == "." || f.search == ".." || fuzzy.Match(f.search, strings.ToLower(name)) {
			if f.search == "." || f.search == ".." || strings.Contains(strings.ToLower(name), strings.ToLower(f.search)) {
				dir := f.dir

				if !abs {
					dir, _ = filepath.Rel(f.wd, dir)
				}

				/*
					if hasTilde {
						dir, _ = filepath.Rel(os.Getenv("HOME"), dir)
						dir = "~/" + dir
					}
				*/

				if hasTilde {
					//dir = strings.Replace(dir, home, "~", 1)
					dir = f.homeToTilde(dir)
				}

				if name != "." && name != ".." {
					pth := filepath.Join(dir, name)
					if fl.IsDir() {
						pth += "/"
					}
					f.files = append(f.files, pth)
				}
			}
		}
	}
	// f.logger.Printf("-------END---------\n")
	sort.Strings(f.files)

	f.pager = pager.New(f.s.height-3, len(f.files))
	f.Print()
	//}()

	/*
	   	p := strings.Split(f.fileSearch, "/")

	   dir,
	*/

	//go FindFiles(f.includeHidden, f.searchDirs, filepath.Join(f.wd, f.dir), f.search, f.onFile, f.onError)

	//go FindFiles(f.includeHidden, f.searchDirs, f.dir, f.search, f.onFile, f.onError)
}
