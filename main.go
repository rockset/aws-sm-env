package main

import (
	"io/ioutil"
	"log"
	"os"
)

const (
	envSecretsManagerPath = "SECRETS_MANAGER_PATH"
	envAssumeRoleArn      = "ASSUME_ROLE_ARN"
)

func main() {
	name := os.Getenv(envSecretsManagerPath)
	if name == "" {
		log.Fatalf("%s environment variable required", envSecretsManagerPath)
	}

	logStream := ioutil.Discard
	if debug := os.Getenv("DEBUG"); debug != "" {
		logStream = os.Stderr
	}

	roleArn := os.Getenv("envAssumeRoleArn")

	si := NewSecretsInjector(logStream, name)
	if err := si.Exec(roleArn, os.Args, os.Environ()); err != nil {
		log.Fatalln(err)
	}
}
