package short

import (
	"encoding/json"
	"fmt"
	"github.com/metakeule/fmtdate"
	"io"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

/*
this package takes a table/data-driven approach to something parametrized shell aliases (called "short.Cut").
*/

type Cut struct {
	Name    string `json:name`    // names may contain dots to create namespaces (to allow grouping)
	Command string `json:command` // a command for the shell. placeholders are defined by inserting something like #Name: Type#, where
	// Name is the name of the parameter and Type is one of string, int, uint, float, date, time
	// placeholders may have defaults and are asked when the shortcut is executed
	// if command starts with a colon, then the followong name is the name of the parent command of which this command is a special case
	Defaults map[string]string `json:defaults` // default values
}

/*
example
[
	   {
			   "name": "ssh.root",
			   "command": "ssh -P #port:int# root@#host:string#:#dir:path#",
			   "defaults": { "port": "22", "dir": "/root"}
			},
			{
				 "name": "ssh.work.root",
				 "command": ":ssh.root",
				 "defaults": { "port": "8999", "dir": "/data"}
			}
	]
*/

var Types = []string{"string", "int", "uint", "float", "date", "time"}

const ParamRegExp = `#([.-_a-z]+):([a-z]+)#`

var regexParamDef = regexp.MustCompilePOSIX(ParamRegExp)

func findParams(command string) [][]string {
	return regexParamDef.FindAllStringSubmatch(command, -1)
}

func replaceParamsInCommand(command string, params map[string]string) string {
	found := findParams(command)
	done := map[string]bool{}

	for i := 0; i < len(found); i++ {
		repl := found[i][0]
		key := found[i][1]
		if _, has := params[key]; has && !done[repl] {
			command = strings.Replace(command, repl, params[key], -1)
			done[repl] = true
		}
	}

	return command
}

func Params(cutName string, allCuts map[string]Cut) (paramsDefinition map[string]string, err error) {

	c := allCuts[cutName]

	cmd, _, err2 := c.commandAndValues(allCuts)

	if err2 != nil {
		err = err2
		return
	}

	var errors map[string]string
	paramsDefinition, errors = findParamsInCommand(cmd)

	if len(errors) > 0 {
		err = fmt.Errorf("ERROR in params definition: %v", errors)
	}
	return
}

/*
func Params(cutName string, allCuts map[string]Cut) (paramsDefinition map[string]string, err error) {
	c := allCuts[cutName]
	cmd := c.Command

	for cmd[0] == ':' {

		c, has := allCuts[cmd[1:]]

		if !has {
			err = fmt.Errorf("unknown cut: %#v, defined as parent for cut %#v", cmd[1:], cutName)
			return
		}

		cutName = cmd[1:]
		cmd = c.Command
	}

	var errors map[string]string
	paramsDefinition, errors = findParamsInCommand(cmd)

	if len(errors) > 0 {
		err = fmt.Errorf("ERROR in params definition: %v", errors)
	}
	return
}
*/

func findParamsInCommand(command string) (paramsDefinition map[string]string, errors map[string]string) {
	paramsDefinition = map[string]string{}
	errors = map[string]string{}
	// the same param may appear several times, but only with the same type
	found := findParams(command)

	for i := 0; i < len(found); i++ {
		key := found[i][1]
		val := found[i][2]
		validType := false
		for _, t := range Types {
			if t == val {
				validType = true
			}
		}

		if !validType {
			errors[key] = "invalid type '" + val + "' for parameter '" + key + "'"
			continue
		}

		if v, has := paramsDefinition[key]; has {
			if val != v {
				errors[key] = "non matching type definitions for parameter '" + key + "': " + "'" + v + "' vs '" + val + "'"
			}
			continue
		}

		paramsDefinition[key] = val

	}
	return
}

func validateValues(command string, params map[string]string) error {
	pDef, defErr := findParamsInCommand(command)

	if len(defErr) > 0 {
		return fmt.Errorf("error in definition of command %#v: %v", command, defErr)
	}

	for k, v := range params {
		def, has := pDef[k]
		if !has {
			continue
			//return fmt.Errorf("parameter '%s' does not exist in command %#v", k, command)
		}

		switch def {
		case "string":
		case "int":
			_, err := strconv.Atoi(v)
			if err != nil {
				return fmt.Errorf("parameter '%s' is not an int", k)
			}
		case "uint":
			i, err := strconv.Atoi(v)
			if err != nil || i < 0 {
				return fmt.Errorf("parameter '%s' is not an uint", k)
			}
		case "float":
			_, err := strconv.ParseFloat(v, 64)
			if err != nil {
				return fmt.Errorf("parameter '%s' is not a float", k)
			}
		case "date":
			_, err := fmtdate.ParseDate(v)
			if err != nil {
				return fmt.Errorf("parameter '%s' is not a date", k)
			}
		case "time":
			_, err := fmtdate.ParseTime(v)
			if err != nil {
				return fmt.Errorf("parameter '%s' is not a time", k)
			}
		default:
			panic("unreachable") // should have been checked via findParamsInCommand
		}

	}

	return nil
}

// valueGroups is a slice of values that is applied in order in that the value in the last map never overwrite
// previous values of the same key. This is to have multiple defaults like with parent Cuts
// the given command is expected to be no parent link, but the command of the final Cut
// the valueGroups should be ordered in such a way, that the first map is the parameters passed at runtime,
// the next the defaults for the chosen cut, the next the defaults of its parent and so on until the final Cut
// that has no parents
func finalValues(valueGroups ...map[string]string) (finals map[string]string) {

	finals = map[string]string{}

	for _, vals := range valueGroups {
		if vals != nil {
			for k, v := range vals {
				if _, has := finals[k]; !has {
					finals[k] = v
				}
			}
		}

	}

	return
}

type Cuts []Cut

func (c Cuts) Len() int {
	return len(c)
}

func (c Cuts) Swap(a, b int) {
	c[a], c[b] = c[b], c[a]
}

func Sort(allCuts map[string]Cut) (cs Cuts) {
	for _, c := range allCuts {
		cs = append(cs, c)
	}

	sort.Sort(cs)
	return
}

func (c Cuts) Less(a, b int) bool {
	return c[a].Name < c[b].Name
}

func Exec(cutName string, allCuts map[string]Cut, runtimeParams map[string]string) (err error) {
	var cmd string
	cmd, err = Command(cutName, allCuts, runtimeParams)

	if err != nil {
		return
	}

	c := exec.Command("/bin/sh", "-c", cmd)
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	return c.Run()
}

func Add(c Cut, allCuts map[string]Cut) error {
	allCuts[c.Name] = c
	_, err := Command(c.Name, allCuts, nil) // validate all the default values
	return err
}

func Save(configJSON io.Writer, allCuts map[string]Cut) (err error) {
	var cuts []Cut

	for _, c := range allCuts {
		cuts = append(cuts, c)
	}

	b, err := json.MarshalIndent(cuts, "", "  ")
	if err != nil {
		return err
	}

	_, err = configJSON.Write(b)
	return
}

func Load(configJSON io.Reader) (allCuts map[string]Cut, err error) {

	var cuts []Cut
	err = json.NewDecoder(configJSON).Decode(&cuts)

	if err != nil {
		err = fmt.Errorf("invalid json: %s", err.Error())
		return
	}

	allCuts = map[string]Cut{}

	for _, c := range cuts {
		if _, has := allCuts[c.Name]; has {
			err = fmt.Errorf("more than one definition of cut %#v", c.Name)
			return
		}
		allCuts[c.Name] = c
	}

	return
}

func (c Cut) commandAndValues(allCuts map[string]Cut) (cmd string, vals []map[string]string, err error) {

	vals = []map[string]string{c.Defaults}

	if c.Command[0] != ':' {
		cmd = c.Command
		return
	}

	refname := c.Command[1:]

	idx := strings.Index(refname, " ")
	var appendix string
	if idx > 0 {
		appendix = strings.TrimSpace(refname[idx:])
		refname = strings.TrimSpace(refname[:idx])
	}

	c, has := allCuts[refname]

	if !has {
		err = fmt.Errorf("unknown cut: %#v, defined as parent for cut %#v", refname, c.Name)
		return
	}

	var grps []map[string]string
	var res string
	res, grps, err = c.commandAndValues(allCuts)

	if err != nil {
		return
	}
	vals = append(vals, grps...)
	cmd = res
	if appendix != "" {
		cmd += " " + appendix

	}

	return
}

func CommandAndValues(cutName string, allCuts map[string]Cut, runtimeParams map[string]string) (cmd string, vals map[string]string, err error) {
	c, has := allCuts[cutName]

	if !has {
		err = fmt.Errorf("unknown cut: %#v", cutName)
		return
	}

	var valueGroups []map[string]string
	cmd, valueGroups, err = c.commandAndValues(allCuts)
	if err != nil {
		return
	}

	vals = finalValues(append([]map[string]string{runtimeParams}, valueGroups...)...)
	err = validateValues(cmd, vals)
	return
}

func _CommandAndValues(cutName string, allCuts map[string]Cut, runtimeParams map[string]string) (cmd string, vals map[string]string, err error) {
	c, has := allCuts[cutName]

	if !has {
		err = fmt.Errorf("unknown cut: %#v", cutName)
		return
	}

	var valueGroups []map[string]string
	valueGroups = append(valueGroups, runtimeParams)

	valueGroups = append(valueGroups, c.Defaults)

	cmd = c.Command

	// TODO: allow appendix after space

	for cmd[0] == ':' {

		c, has := allCuts[cmd[1:]]

		if !has {
			err = fmt.Errorf("unknown cut: %#v, defined as parent for cut %#v", cmd[1:], cutName)
			return
		}

		valueGroups = append(valueGroups, c.Defaults)
		cutName = cmd[1:]
		cmd = c.Command
	}

	vals = finalValues(valueGroups...)
	err = validateValues(cmd, vals)
	return

}

func Command(cutName string, allCuts map[string]Cut, runtimeParams map[string]string) (cmd string, err error) {
	var vals map[string]string
	cmd, vals, err = CommandAndValues(cutName, allCuts, runtimeParams)
	if err != nil {
		return
	}
	return replaceParamsInCommand(cmd, vals), nil
}
