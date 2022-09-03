package main

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/hectane/go-acl"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sys/windows"
)

func TestPermissionIs(t *testing.T) {
	filePath := "misc/tests/permTest.txt"
	ioutil.WriteFile(filePath, []byte("testContent"), os.ModePerm)
	defer os.Remove(filePath)

	c := cond{
		Path:  filePath,
		Name:  "BUILTIN\\Users",
		Value: "Read",
	}

	// users deny read access
	acl.Apply(filePath, true, false, acl.DenyName(windows.GENERIC_READ, "BUILTIN\\Users"))
	ok, err := c.PermissionIs()
	assert.Nil(t, err, "Should not have error")
	assert.NotEqual(t, ok, true, "must be Read")

	// users read access
	acl.Apply(filePath, true, false, acl.GrantName(windows.GENERIC_READ, "BUILTIN\\Users"))
	ok, err = c.PermissionIs()
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, ok, true, "must be Read")

	c.Name = "everyone"
	c.Value = "FullControl"

	// everyone deny full access
	acl.Apply(filePath, true, false, acl.DenyName(windows.GENERIC_ALL, "everyone"))
	ok, err = c.PermissionIs()
	assert.Nil(t, err, "Should not have error")
	assert.NotEqual(t, ok, true, "Must be FullControl")

	// everyone full access
	acl.Apply(filePath, true, false, acl.GrantName(windows.GENERIC_ALL, "everyone"))
	ok, err = c.PermissionIs()
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, ok, true, "Must be FullControl")
}
