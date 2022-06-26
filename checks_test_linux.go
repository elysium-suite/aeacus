package main

import "testing"

func TestCommand(t *testing.T) {
	c := cond{
		Cmd: "echo 'hello, world!'",
	}

	// Should pass: command ran
	out, err := c.Command()
	if err != nil || out != true {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command execution fails
	c.Cmd = "commanddoesntexist"
	out, err = c.Command()
	if err == nil || out != false {
		t.Error(c, "failed:", out, err)
	}

	// Should fail: command returns error
	c.Cmd = "cat /etc/file/doesnt/exist"
	out, err = c.Command()
	if err == nil || out != false {
		t.Error(c, "failed:", out, err)
	}
}
