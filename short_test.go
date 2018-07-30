package short

import (
	"reflect"
	"strings"
	"testing"
)

func TestFindParams(t *testing.T) {
	tests := []struct {
		command  string
		expected [][]string
	}{
		{
			"ssh -P #port:int# root@#host:string#:#dir:string#",
			[][]string{{"#port:int#", "port", "int"}, {"#host:string#", "host", "string"}, {"#dir:string#", "dir", "string"}},
		},
		{
			"ssh root@localhost",
			nil,
		},
	}

	for _, test := range tests {

		if got, want := findParams(test.command), test.expected; !reflect.DeepEqual(got, want) {
			t.Errorf("findParams(%#v) = %v; want %v", test.command, got, want)
		}
	}

}

func TestFindParamsInCommand(t *testing.T) {
	tests := []struct {
		command          string
		paramsDefinition map[string]string
		errors           map[string]string
	}{
		{
			"ssh -P #port:int# root@#host:string#:#dir:string#",
			map[string]string{
				"port": "int",
				"host": "string",
				"dir":  "string",
			},
			map[string]string{},
		},
		{
			"ssh root@localhost",
			map[string]string{},
			map[string]string{},
		},
		{
			"ssh #user:string#@remote:/home/#user:string#",
			map[string]string{
				"user": "string",
			},
			map[string]string{},
		},
		{
			"ssh #user:string#@remote:/home/#user:int#",
			map[string]string{
				"user": "string",
			},
			map[string]string{
				"user": "non matching type definitions for parameter 'user': 'string' vs 'int'",
			},
		},
		{
			"ssh root@#host:strong#",
			map[string]string{},
			map[string]string{
				"host": "invalid type 'strong' for parameter 'host'",
			},
		},
	}

	for _, test := range tests {
		gotParamsDefinition, gotErrors := findParamsInCommand(test.command)

		if got, want := gotErrors, test.errors; !reflect.DeepEqual(got, want) {
			t.Errorf("findParamsInCommand(%#v) = _,%v; want _,%v", test.command, got, want)
		}

		if got, want := gotParamsDefinition, test.paramsDefinition; !reflect.DeepEqual(got, want) {
			t.Errorf("findParamsInCommand(%#v) = %v,_; want %v,_", test.command, got, want)
		}

	}

}

// validateValues

func TestValidateValues(t *testing.T) {
	tests := []struct {
		command string
		params  map[string]string
		error   string
	}{
		{
			"ssh -P #port:int# root@#host:string#:#dir:string#",
			map[string]string{
				"port": "8999",
				"host": "localhost",
				"dir":  "/",
			},
			"",
		},
		{
			"ssh -P #port:int# root@#host:string#:#dir:string#",
			map[string]string{
				"port": "89.99",
				"host": "localhost",
				"dir":  "/",
			},
			"parameter 'port' is not an int",
		},
	}

	for _, test := range tests {
		gotError := validateValues(test.command, test.params)
		got := ""
		if gotError != nil {
			got = gotError.Error()
		}
		if want := test.error; got != want {
			t.Errorf("validateValues(%#v, %v) = %#v; want %#v", test.command, test.params, got, want)
		}

	}

}

func TestFinalValues(t *testing.T) {

	var valGroups = []map[string]string{
		{
			"a": "b",
			"x": "y",
		},
		{
			"a": "c",
			"c": "d",
		},

		{
			"a": "d",
			"c": "e",
			"f": "g",
		},
	}

	vals := finalValues(valGroups...)

	expected := map[string]string{
		"x": "y",
		"a": "b",
		"c": "d",
		"f": "g",
	}

	if got, want := vals, expected; !reflect.DeepEqual(got, want) {
		t.Errorf("finalValues(%#v) = %v; want %v", valGroups, got, want)
	}
}

func TestLoad(t *testing.T) {
	const conf = `[
	   {
			   "name": "ssh.root",
			   "command": "ssh -P #port:int# root@#host:string#:#dir:string#",
			   "defaults": { "port": "22", "dir": "/root"}
			},
			{
				 "name": "ssh.work.root",
				 "command": ":ssh.root",
				 "defaults": { "port": "8999", "dir": "/data"}
			}
	]`

	cuts, err := Load(strings.NewReader(conf))

	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	if got, want := cuts, allCuts; !reflect.DeepEqual(got, want) {
		t.Errorf("Load(%#v) = %v; want %v", conf, got, want)
	}

}

var allCuts = map[string]Cut{
	"ssh.root": {
		"ssh.root",
		"ssh -P #port:int# root@#host:string#:#dir:string#",
		map[string]string{"port": "22", "dir": "/root"},
	},
	"ssh.work.root": {
		"ssh.work.root",
		":ssh.root",
		map[string]string{"port": "8999", "dir": "/data"},
	},
}

func TestCommand(t *testing.T) {

	cmd, err := Command("ssh.work.root", allCuts, map[string]string{"host": "localhost"})

	if err != nil {
		t.Errorf("Error: %v", err)
		return
	}

	expected := "ssh -P 8999 root@localhost:/data"

	if got, want := cmd, expected; got != want {
		t.Errorf(
			"Command(\"ssh.work.root\", allCuts, map[string]string{\"host\": \"localhost\"}) = %v; want %v",
			got, want)
	}
}
