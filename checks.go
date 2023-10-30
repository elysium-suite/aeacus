// checks.go contains checks that are identical for both Linux and Windows.

package main

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"

	"github.com/safinsingh/aeaconf2"
)

var funcRegistry map[string]reflect.Type

func init() {
	funcRegistry = make(map[string]reflect.Type)

	funcRegistry["AutoCheckUpdatesEnabled"] = reflect.TypeOf(AutoCheckUpdatesEnabled{})
	funcRegistry["DirContains"] = reflect.TypeOf(DirContains{})
	funcRegistry["FileContains"] = reflect.TypeOf(FileContains{})
	funcRegistry["PathExists"] = reflect.TypeOf(PathExists{})

	aeaconf2.CheckFunctionRegistry(funcRegistry)
}

// requireArgs is a convenience function that prints a warning if any required
// parameters for a given condition are not provided.
func (c cond) requireArgs(args ...interface{}) {
	// Don't process internal calls -- assume the developers know what they're
	// doing. This also prevents extra errors being printed when they don't pass
	// required arguments.
	if c.Type == "" {
		return
	}

	v := reflect.ValueOf(c)
	vType := v.Type()
	for i := 0; i < v.NumField(); i++ {
		if vType.Field(i).Name == "Type" || vType.Field(i).Name == "regex" {
			continue
		}

		// Ignore hint fields, they only show up in the scoring report
		if vType.Field(i).Name == "Hint" {
			continue
		}

		required := false
		for _, a := range args {
			if vType.Field(i).Name == a {
				required = true
				break
			}
		}

		if required {
			if v.Field(i).String() == "" {
				fail(c.Type+":", "missing required argument '"+vType.Field(i).Name+"'")
			}
		} else if v.Field(i).String() != "" {
			warn(c.Type+":", "specifying unused argument '"+vType.Field(i).Name+"'")
		}
	}
}

func (c cond) String() string {
	output := ""
	v := reflect.ValueOf(c)
	typeOfS := v.Type()

	for i := 0; i < v.NumField(); i++ {
		if v.Field(i).String() == "" {
			continue
		}
		output += fmt.Sprintf("\t%s: %v\n", typeOfS.Field(i).Name, v.Field(i).String())
	}
	return output
}

func handleReflectPanic(condFunc string) {
	if r := recover(); r != nil {
		fail("Check type does not exist: "+condFunc, "("+r.(*reflect.ValueError).Error()+")")
	}
}

// runCheck executes a single condition check.
func runCheck(cond aeaconf2.Condition) bool {
	if err := deobfuscateCond(cond); err != nil {
		fail(err.Error())
	}
	defer obfuscateCond(cond)
	debug("Running condition:\n", cond)

	not := "Not"
	regex := "Regex"
	condFunc := ""
	negation := false
	cond.regex = false

	// Ensure that condition type is a valid length
	if len(cond.Type) <= len(regex) {
		fail(`Condition type "` + cond.Type + `" is not long enough to be valid. Do you have a "type = 'CheckTypeHere'" for all check conditions?`)
		return false
	}
	condFunc = cond.Type
	if condFunc[len(condFunc)-len(not):] == not {
		negation = true
		condFunc = condFunc[:len(condFunc)-len(not)]
	}
	if condFunc[len(condFunc)-len(regex):] == regex {
		cond.regex = true
		condFunc = condFunc[:len(condFunc)-len(regex)]
	}

	// Catch panic if check type doesn't exist
	defer handleReflectPanic(condFunc)

	// Using reflection to find the correct function to call.
	vals := reflect.ValueOf(cond).MethodByName(condFunc).Call([]reflect.Value{})
	result := vals[0].Bool()
	err := vals[1]

	if negation {
		debug("Result is", !result, "(was", result, "before negation) and error is", err)
		return err.IsNil() && !result
	}

	debug("Result is", result, "and error is", err)

	if verboseEnabled && !err.IsNil() {
		warn(condFunc, "returned an error:", err)
	}

	return err.IsNil() && result
}

// CommandContains checks if a given shell command contains a certain string.
// This check will always fail if the command returns an error.
func (c cond) CommandContains() (bool, error) {
	c.requireArgs("Cmd", "Value")
	out, err := shellCommandOutput(c.Cmd)
	if err != nil {
		return false, err
	}
	if c.regex {
		outTrim := strings.TrimSpace(out)
		return regexp.Match(c.Value, []byte(outTrim))
	}
	return strings.Contains(strings.TrimSpace(out), c.Value), err
}

// CommandOutput checks if a given shell command produces an exact output.
// This check will always fail if the command returns an error.
func (c cond) CommandOutput() (bool, error) {
	c.requireArgs("Cmd", "Value")
	out, err := shellCommandOutput(c.Cmd)
	return strings.TrimSpace(out) == c.Value, err
}

// DirContains returns true if any file in the directory contains the string value provided.
type DirContains struct {
	aeaconf2.BaseCondition
	Path  string
	Value string
	Regex bool
}

func (d *DirContains) Score() (bool, error) {
	result, err := (&PathExists{Path: d.Path}).Score()
	if err != nil {
		return false, err
	}
	if !result {
		return false, errors.New("path does not exist")
	}

	var files []string
	err = filepath.Walk(d.Path, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, path)
		}
		if len(files) > 10000 {
			return errors.New("attempted to index too many files in recursive search")
		}
		return nil
	})

	if err != nil {
		return false, err
	}

	for _, file := range files {
		result, err := (&FileContains{Path: file, Value: d.Value}).Score()
		if os.IsPermission(err) {
			return false, err
		}
		if result {
			return result, nil
		}
	}
	return false, nil
}

// FileContains determines whether a file contains a given regular expression.
//
// Newlines in regex may not work as expected, especially on Windows. It's
// best to not use these (ex. ^ and $).
type FileContains struct {
	aeaconf2.BaseCondition
	Path  string
	Value string
	Regex bool
}

func (f *FileContains) Score() (bool, error) {
	fileContent, err := readFile(f.Path)
	if err != nil {
		return false, err
	}
	found := false
	for _, line := range strings.Split(fileContent, "\n") {
		if f.Regex {
			found, err = regexp.Match(f.Value, []byte(line))
			if err != nil {
				return false, err
			}
		} else {
			found = strings.Contains(line, f.Value)
		}
		if found {
			break
		}
	}
	return found, err
}

// FileEquals calculates the SHA256 sum of a file and compares it with the hash
// provided in the check.
func (c cond) FileEquals() (bool, error) {
	c.requireArgs("Path", "Value")
	fileContent, err := readFile(c.Path)
	if err != nil {
		return false, err
	}
	hasher := sha256.New()
	_, err = hasher.Write([]byte(fileContent))
	if err != nil {
		return false, err
	}
	hash := hex.EncodeToString(hasher.Sum(nil))
	return hash == c.Value, nil
}

// PathExists is a wrapper around os.Stat and os.IsNotExist, and determines
// whether a file or folder exists.
type PathExists struct {
	aeaconf2.BaseCondition
	Path string
}

func (p *PathExists) Score() (bool, error) {
	_, err := os.Stat(p.Path)
	if err != nil && os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}
