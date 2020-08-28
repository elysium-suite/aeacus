package cmd

import "testing"

func boolToString(b bool) string {
	if b == true {
		return "true"
	}
	return "false"
}

func TestCommandOutput(t *testing.T) {
	out, err := commandOutput(`echo 1`, "1")
	if err != nil || out != true {
		t.Error("commandOutput(`echo 1`, \"1\") got " + boolToString(out) + ", want `true`")
	}
}
func TestCommandContains(t *testing.T) {
	out, err := commandContains(`echo hello world`, "hello")
	if err != nil || out != true {
		t.Error("commandContains(`echo hello world!`, \"hello\") got " + boolToString(out) + ", want `true`")
	}
}

func TestPathExists(t *testing.T) {
	out, err := pathExists("/")
	if err != nil || out != true {
		t.Error("pathExists(\"/\") got " + boolToString(out) + ", want `true`")
	}
}

func TestFileContains(t *testing.T) {
	out, err := fileContains("../misc/tests/TestFileContains.txt", "world")
	if err != nil || out != true {
		t.Error("fileContains(\"../misc/tests/TestFileContains.txt\", \"hello\") got " + boolToString(out) + ", want `true`")
	}
}

func TestFileContainsRegex(t *testing.T) {
	out, err := fileContainsRegex("../misc/tests/TestFileContains.txt", "^hello")
	if err != nil || out != true {
		t.Error("fileContainsRegex(\"../misc/tests/TestFileContains.txt\", \"^hello\") got " + boolToString(out) + ", want `true`. Error: " + err.Error())
	}
}

func TestDirContainsRegex(t *testing.T) {
	out, err := dirContainsRegex("../misc/tests/dir", "^efgh")
	if err != nil || out != true {
		t.Error("dirContainsRegex(\"../misc/tests/dir\", \"^efgh\") got " + boolToString(out) + err.Error() + ", want `true`")
	}
}
