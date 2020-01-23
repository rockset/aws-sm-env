# `aws-sm-env`

This tool lets you inject AWS Secrets Manager secrets as environment variables in a command.

```
SECRETS_MANAGER_PATH=buildkite/token,buildkite/packer aws-sm-env packer build buildkite.json
```

The `SECRETS_MANAGER_PATH` is a comma-separated list of Secrets Manager paths, and it the same
secret exist in multiple paths, the last one will be used.

If `SECRETS_MANAGER_PATH` is not set, `aws-sm-env` will act as a pass-though and just call
`syscall.Exec()` with the arguments.

It also lets you assume an IAM role after fetching the secrets, if `ASSUME_ROLE_ARN` is set.

```
ASSUME_ROLE_ARN=arn:aws:iam::216640736862:role/buildkite-agent SECRETS_MANAGER_PATH=packer/buildkite aws-sm-env packer build buildkite.json
```

## Release

This repo uses Circle CI to automatically build and release new version when a SemVer tag is detected.
