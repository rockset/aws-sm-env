# `aws-sm-env`

A tool to inject AWS Secrets Manager secrets as environment variables in a command.
It also lets you assume a IAM role after fetching the secrets.

```
ASSUME_ROLE_ARN=arn:aws:iam::216640736862:role/buildkite-agent SECRETS_MANAGER_PATH=packer/buildkite aws-sm-env packer build buildkite.json
```

## Release

This repo uses Circle CI to automatically build and release new version when a SemVer tag is detected.
