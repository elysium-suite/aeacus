package cmd

import (
	"testing"
)

func TestRegistryKeyExists(t *testing.T) {
	keyName := `HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\CrashControl`
	keyValue := `Overwrite`
	out, err := registryKey(keyName, keyValue, true)
	if err != nil {
		t.Error(`registryKey with`, keyName, keyValue, "error "+err.Error())
	} else if out != true {
		t.Error(`registryKey with`, keyName, keyValue, "got false, want true")
	}
}
