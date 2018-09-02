package main

import (
	"github.com/lithammer/fuzzysearch/fuzzy"
	"github.com/stretchr/powerwalk"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type findFilesAndDirs struct {
	files         chan string
	errors        chan error
	filesOnly     bool
	wd            string
	search        string
	callbackFile  func(f string)
	callbackError func(err error)
	includeHidden bool
	searchDirs    bool
}

func FindFiles(includeHidden bool, searchDirs bool, wd, search string, callbackFile func(f string), callbackError func(err error)) {
	f := &findFilesAndDirs{}
	f.files = make(chan string, 20)
	f.errors = make(chan error, 20)
	f.wd = wd         // , _ = os.Getwd()
	f.search = search // "*.go" //os.Args[1]
	f.callbackError = callbackError
	f.callbackFile = callbackFile
	f.searchDirs = searchDirs
	f.includeHidden = includeHidden
	f.Find()
}

func (f *findFilesAndDirs) Find() {
	// panic("find called")
	var wg sync.WaitGroup
	wg.Add(3)
	// var fin = make(chan bool)
	go func() {
		for fl := range f.files {
			//fmt.Printf("%#v\n", fl)
			f.callbackFile(fl)
		}
		wg.Done()
	}()

	go func() {
		for e := range f.errors {
			// fmt.Printf("ERROR: %#v\n", e.Error())
			if f.callbackError != nil {

				f.callbackError(e)
			}
		}
		wg.Done()
	}()

	powerwalk.WalkLimit(f.wd, f.walkFn, 20)
	wg.Done()
	// panic("find done")
	close(f.files)
	close(f.errors)
	wg.Wait()
	// fin <- true
}

func (f *findFilesAndDirs) isHidden(path string) bool {
	path, _ = filepath.Abs(path)
	p := strings.Split(path, "/")

	for _, pa := range p {
		if len(pa) > 0 && pa[0] == '.' {
			return true
		}
	}

	return false

}

func (f *findFilesAndDirs) walkFn(path string, info os.FileInfo, err error) error {
	if err != nil {
		f.errors <- err
		return nil
	}

	if info.IsDir() == f.searchDirs {

		// ok, er := filepath.Match(f.search, info.Name())
		/*
			if er != nil {
				f.errors <- er
				return nil
			}
		*/

		if !f.includeHidden && f.isHidden(path) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// panic(info.Name())
		// files, er := filepath.Glob(f.search, info.Name())
		/*
		   if er != nil {
		   				f.errors <- er
		   				return nil
		   			}

		   			file
		*/
		_ = fuzzy.Match

		//	if fuzzy.Match(f.search, info.Name()) {
		if ok, _ := filepath.Match(f.search, info.Name()); ok {
			//fmt.Println("- " + info.Name())
			p, _ := filepath.Rel(f.wd, path)
			f.files <- p
		}
		/*
		 */
	}
	return nil
}
