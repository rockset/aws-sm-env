# `aws-sm-env`

A tool to inject AWS Secrets Manager secrets as environment variables in a command

```
SECRETS_MANAGER_PATH=packer/buildkite aws-sm-env packer build buildkite.json
```

## Release

This repo uses Circle CI to automatically build and release new version when a SemVer tag is detected.
