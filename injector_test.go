package main

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

func TestMakeEnv(t *testing.T) {
	sm := NewSecretsInjector(ioutil.Discard, "NAME")
	secrets := map[string]string{
		"foo":  "bar",
		"HOME": "/home/bar",
	}
	env, err := sm.inject([]string{"HOME=/home/foo", "BADENV"}, secrets)
	assert.Nil(t, err)
	assert.Contains(t, env, "foo=bar")
	assert.Contains(t, env, "HOME=/home/bar")
	assert.NotContains(t, env, "BADENV")
}

func TestFilterArgs(t *testing.T) {
	path, args, err := filterArgs([]string{"aws-sm-env", "ls", "foo"})
	assert.Nil(t, err)
	assert.Equal(t, path, "/bin/ls")
	assert.Equal(t, args, []string{"ls", "foo"})
}
