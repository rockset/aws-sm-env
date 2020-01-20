package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/secretsmanager"
	"io"
	"log"
	"os/exec"
	"strings"
	"syscall"
	"time"
)

type SecretsInjector struct {
	log  *log.Logger
	name string
}

func NewSecretsInjector(logStream io.Writer, name string) *SecretsInjector {
	return &SecretsInjector{
		log:  log.New(logStream, "", log.LstdFlags),
		name: name,
	}
}

// Exec runs the exec() syscall with the supplied arguments
func (si *SecretsInjector) Exec(args, env []string) error {
	path, filteredArgs, err := filterArgs(args)
	if err != nil {
		return err
	}

	secrets, err := si.getSecrets()
	if err != nil {
		return err
	}

	injectedEnv, err := si.inject(env, secrets)
	if err != nil {
		return err
	}

	return syscall.Exec(path, filteredArgs, injectedEnv)
}

func filterArgs(args []string) (string, []string, error) {
	if len(args) == 1 {
		return "", nil, fmt.Errorf("must specify command to execute")
	}

	filtered := args[1:]
	cmd := filtered[0]
	path, err := exec.LookPath(cmd)
	if err != nil {
		return "", nil, fmt.Errorf("failed to find command %s: %v", cmd, err)
	}

	return path, filtered, nil
}

func (si *SecretsInjector) getSecrets() (map[string]string, error) {
	secrets := make(map[string]string)

	sess, err := session.NewSession()
	if err != nil {
		return secrets, fmt.Errorf("failed to create AWS session: %w", err)
	}

	svc := secretsmanager.New(sess)
	for _, name := range strings.Split(si.name, ",") {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		result, err := svc.GetSecretValueWithContext(ctx, &secretsmanager.GetSecretValueInput{
			SecretId: aws.String(name),
		})
		cancel()

		if err != nil {
			return secrets, fmt.Errorf("failed to get secret %s: %w", name, err)
		}

		if result.SecretString == nil || *result.SecretString == "" {
			return secrets, fmt.Errorf("secret %s does not contain any secret string", name)
		}

		var s map[string]string
		err = json.Unmarshal([]byte(*result.SecretString), &s)
		if err != nil {
			return secrets, fmt.Errorf("failed to unmasrshal secret string: %w", err)
		}

		for k, v := range s {
			secrets[k] = v
		}
	}

	return secrets, err
}

func (si *SecretsInjector) inject(env []string, secrets map[string]string) ([]string, error) {
	filtered := make(map[string]string)
	for _, e := range env {
		tokens := strings.SplitN(e, "=", 2)
		if len(tokens) != 2 {
			si.log.Printf("failed to parse environment string: %s", e)
			continue
		}
		filtered[tokens[0]] = tokens[1]
	}

	for key, value := range secrets {
		if v, found := filtered[key]; found {
			si.log.Printf("replacing environment variable %s (%s) with %s", key, v, value)
		}
		filtered[key] = value
	}

	newEnv := make([]string, len(filtered))
	var i int
	for k, v := range filtered {
		newEnv[i] = fmt.Sprintf("%s=%s", k, v)
		i++
	}

	return newEnv, nil
}
