// checks_test.go is responsible for testing all non-platform dependent checks.
package main

import (
	"testing"
)

func TestCommandContains(t *testing.T) {
	c := cond{
		Cmd:   "echo 'hello, world!'",
		Value: "hello, world!",
	}

	// Should pass: exact match
	out, err := c.CommandContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should pass: substring
	c.Value = "hello"
	out, err = c.CommandContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: not substring
	c.Value = "bye"
	out, err = c.CommandContains()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command execution fails
	c.Value = ""
	c.Cmd = "commanddoesntexist"
	out, err = c.CommandContains()
	if err == nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command returns error
	c.Cmd = "cat /etc/file/doesnt/exist"
	out, err = c.CommandContains()
	if err == nil || out != true {
		t.Error(c, "failed:", out, err)
	}
}

func TestCommandOutput(t *testing.T) {
	c := cond{
		Cmd:   "echo 'hello, world!'",
		Value: "hello, world!",
	}

	// Should pass: exact match
	out, err := c.CommandOutput()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: just substring
	c.Value = "hello"
	out, err = c.CommandOutput()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: not exact or substring
	c.Value = "bye"
	out, err = c.CommandOutput()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command execution fails
	c.Value = ""
	c.Cmd = "commanddoesntexist"
	out, err = c.CommandOutput()
	if err == nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command returns error
	c.Cmd = "cat /etc/file/doesnt/exist"
	out, err = c.CommandOutput()
	if err == nil || out != true {
		t.Error(c, "failed:", out, err)
	}
}

func TestDirContains(t *testing.T) {
	c := cond{
		Path:  "misc/tests/dir",
		Value: "^efgh",
	}
	out, err := c.DirContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	c.Value = "^efghabcd$"
	out, err = c.DirContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	c.Value = "^aaaaaa$"
	out, err = c.DirContains()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}

	c.Value = `spaces\s+in\s+it\s+[0-9]*\s+nums`
	out, err = c.DirContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	c.Value = `spaces\s+in\s+it\s+[1-5]*\s+nums`
	out, err = c.DirContains()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}

}

func TestFileContains(t *testing.T) {
	c := cond{
		Path:  "misc/tests/TestFileContains.txt",
		Value: "^hello",
	}
	out, err := c.FileContains()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	c.Value = "nothere"
	out, err = c.FileContains()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}
}

func TestPathExists(t *testing.T) {
	c := cond{
		Path: "misc/tests/",
	}
	out, err := c.PathExists()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	c.Path = "misc/doesntexist"
	out, err = c.PathExists()
	if err != nil || out != false {
		t.Error(c, "failed:", out, err)
	}
}
