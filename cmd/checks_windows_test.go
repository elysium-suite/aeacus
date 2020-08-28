package cmd

import (
	"testing"
)

func TestRegistryKeyExists(t *testing.T) {
	out, err := registryKey(`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\CrashControl`, "Overwrite", true)
	if err != nil || out != true {
		t.Error("registryKey(`HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Control\CrashControl`, \"Overwrite\") got " + boolToString(out) + ", want `true`. Error " + err.Error())
	}
}
