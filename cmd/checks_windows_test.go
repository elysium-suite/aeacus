package cmd

import (
	// "fmt"
	"io/ioutil"
	"os"
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/hectane/go-acl"
	"golang.org/x/sys/windows"
)

func Test_filePermission(t *testing.T){
	filePath := "test1.txt"
	ioutil.WriteFile(filePath, []byte("test"), os.ModePerm)
	defer os.Remove(filePath)

	// users deny read access
	acl.Apply(filePath, true, false, acl.DenyName(windows.GENERIC_READ, "BUILTIN\\Users"))
	ok, err := filePermission(filePath, "BUILTIN\\Users", "Read")
	assert.Nil(t, err, "Should not have error")
	assert.NotEqual(t, ok, true, "must be Read")

	// users read access
	acl.Apply(filePath, true, false, acl.GrantName(windows.GENERIC_READ, "BUILTIN\\Users"))
	ok, err = filePermission(filePath, "BUILTIN\\Users", "Read")
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, ok, true, "must be Read")
	
	// everyone deny full access
	acl.Apply(filePath, true, false, acl.DenyName(windows.GENERIC_ALL, "everyone"))
	ok, err = filePermission(filePath, "everyone", "FullControl")
	assert.Nil(t, err, "Should not have error")
	assert.NotEqual(t, ok, true, "must be fullControl")

	// everyone full access
	acl.Apply(filePath, true, false, acl.GrantName(windows.GENERIC_ALL, "everyone"))
	ok, err = filePermission(filePath, "everyone", "FullControl")
	assert.Nil(t, err, "Should not have error")
	assert.Equal(t, ok, true, "must be fullControl")
}