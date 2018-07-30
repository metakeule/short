package main

import (
	"fmt"
	"github.com/metakeule/short"
	"os"
	"sort"
	"strings"
)

type Field struct {
	Cursor int
	Value  string
}

func (f *Field) Delete() {
	rs := []rune(f.Value)
	var val string
	if len(rs) > 0 {
		if f.Cursor > 0 {
			val = string(rs[:f.Cursor-1])
			if len(rs)+3 > f.Cursor {
				val += string(rs[f.Cursor:])
			}
			f.Value = val
			f.Cursor--
		} else {
			val = string(rs[0 : len(rs)-1])
			f.Cursor = len(f.Value)
		}
	}
}

func (f *Field) Insert(r rune) {
	if f.Cursor < 0 {
		f.Cursor = len(f.Value)
	}

	if f.Cursor > 0 {
		orig := f.Value
		if len(orig)+2 > f.Cursor {
		}
		f.Value = string(r)
		if len(orig)+2 > f.Cursor {
			f.Value = orig[:f.Cursor] + f.Value + orig[f.Cursor:]
		}
		f.Cursor++
	} else {
		f.Value += string(r)
		f.Cursor = len(f.Value)
	}
}

func mapToString(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}

	var lines []string

	for k, v := range m {
		lines = append(lines, fmt.Sprintf("%s=%s", k, v))
	}

	sort.Strings(lines)

	return "[" + strings.Join(lines, ",") + "]"
}

type params [][3]string

func (p params) Len() int {
	return len(p)

}

func (p params) Less(a, b int) bool {
	return p[a][0] < p[b][0]

}

func (p params) Swap(a, b int) {
	p[a], p[b] = p[b], p[a]
}

func paramsStringToMap(s string) (m map[string]string, err error) {
	m = map[string]string{}

	params := strings.Split(s, ",")

	for _, dg := range params {
		dg = strings.TrimSpace(dg)

		if dg == "" {
			continue
		}

		dp := strings.Split(dg, "=")
		if len(dp) != 2 {
			err = fmt.Errorf("invalid default string, use 'a=b,c=1' etc")
			return
		}

		m[strings.TrimSpace(dp[0])] = strings.TrimSpace(dp[1])
	}

	return
}

func save(allCuts map[string]short.Cut) error {
	f, err := os.Create(SHORTCUT_FILE)
	if err != nil && err != os.ErrNotExist {
		return err
	}

	defer f.Close()

	return short.Save(f, allCuts)
}

func load() (allCuts map[string]short.Cut, err error) {
	var f *os.File
	f, err = os.Open(SHORTCUT_FILE)
	if err != nil {
		_, notExist := err.(*os.PathError)
		if notExist {
			err = nil
			allCuts = map[string]short.Cut{}
			return
		}

		return
	}

	defer f.Close()

	return short.Load(f)
}
