package main

import (
	"io/ioutil"
	"log"
	"os"
)

const envVarName = "SECRETS_MANAGER_PATH"

func main() {
	name := os.Getenv(envVarName)
	if name == "" {
		log.Fatalf("%s environment variable required", envVarName)
	}

	logStream := ioutil.Discard
	if debug := os.Getenv("DEBUG"); debug != "" {
		logStream = os.Stderr
	}

	si := NewSecretsInjector(logStream, name)
	if err := si.Exec(os.Args, os.Environ()); err != nil {
		log.Fatalln(err)
	}
}
