package cmd

import "testing"

func TestCommandOutput(t *testing.T) {
	out, err := commandOutput(`echo 1`, "1")
	if err != nil || out != true {
		t.Error("commandOutput(`echo 1`, \"1\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
	}
}

func TestCommandContains(t *testing.T) {
	out, err := commandContains(`echo hello world`, "hello")
	if err != nil || out != true {
		t.Error("commandContains(`echo hello world!`, \"hello\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
	}
}

func TestPathExists(t *testing.T) {
	out, err := pathExists("/")
	if err != nil || out != true {
		t.Error("pathExists(\"/\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
	}
}

func TestFileContains(t *testing.T) {
	out, err := fileContains("../misc/tests/TestFileContains.txt", "world")
	if err != nil || out != true {
		t.Error("fileContains(\"../misc/tests/TestFileContains.txt\", \"hello\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
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
		t.Error("dirContainsRegex(\"../misc/tests/dir\", \"^efgh\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
	}
}

func TestFileEquals(t *testing.T) {
	fileName := "../misc/tests/TestFileContains.txt"
	hashText := "a594c7d7d9065ad2e37a57cc065a8b5c75a6cd90"
	out, err := fileEquals(fileName, hashText)
	if err != nil {
		t.Error("fileEquals for file", fileName, "hash", hashText, err.Error())
	} else if out != true {
		t.Error("fileEquals for file", fileName, "hash", hashText, "got false, want true")
	}
}
