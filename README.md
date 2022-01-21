# sm - AWS Secrets Manager CLI Tool

[![License: MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://github.com/clok/sm/blob/master/LICENSE)
[![Go Report Card](https://goreportcard.com/badge/clok/sm)](https://goreportcard.com/report/clok/sm)
[![Coverage Status](https://coveralls.io/repos/github/clok/sm/badge.svg)](https://coveralls.io/github/clok/sm)
[![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white)](https://pkg.go.dev/github.com/clok/sm?tab=overview)

> Please see [the docs for details on the commands.](./docs/sm.md)

```text
$ sm --help
NAME:
   sm - AWS Secrets Manager CLI Tool

USAGE:
   sm [global options] command [command options] [arguments...]

COMMANDS:
   get, view    select from list or pass in specific secret
   edit, e      interactive edit of a secret String Value
   create, c    create new secret in Secrets Manager
   put          non-interactive update to a specific secret
   delete, del  delete a specific secret
   list         display table of all secrets with meta data
   describe     print description of secret to `STDOUT`
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help (default: false)
   --version, -v  print the version (default: false)

COPYRIGHT:
   (c) 2021 Derek Smith
```

- [Documentation](./docs/sm.md)
- [Authentication](#authentication)
- [Installation](#installation)
   - [Homebrew](#homebrewhttpsbrewsh-for-macos-users)
   - [curl binary](#curl-binary)
   - [docker](#dockerhttpswwwdockercom)
- [Development](#development)
- [Versioning](#versioning)
- [Authors](#authors)
- [License](#license)

## Authentication

This tool uses the [awssession](https://github.com/clok/awssession) module for creating authenticated sessions. This will
use AWS Instance Role, Environment Variables or AWS CLI configuration files to generate a session. This tool should obey
all [AWS CLI Environment Variables](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-envvars.html).

## Installation

### [Homebrew](https://brew.sh) (for macOS users)

```
brew tap clok/sm
brew install sm
```

### curl binary

```
$ curl https://i.jpillora.com/clok/sm! | bash
```

### [docker](https://www.docker.com/)

The compiled docker images are maintained on [GitHub Container Registry (ghcr.io)](https://github.com/orgs/clok/packages/container/package/sm).
We maintain the following tags:

- `edge`: Image that is build from the current `HEAD` of the main line branch.
- `latest`: Image that is built from the [latest released version](https://github.com/clok/sm/releases)
- `x.y.z` (versions): Images that are build from the tagged versions within Github.

```bash
docker pull ghcr.io/clok/sm
docker run -v "$PWD":/workdir ghcr.io/clok/sm --version
```

### man page

To install `man` page:

```
$ sm install-manpage
```

## Development

1. Fork the [clok/sm](https://github.com/clok/sm) repo
1. Use `go >= 1.16`
1. Branch & Code
1. Run linters :broom: `golangci-lint run`
    - The project uses [golangci-lint](https://golangci-lint.run/usage/install/#local-installation)
1. Commit with a Conventional Commit
1. Open a PR

## Versioning

We employ [git-chglog](https://github.com/git-chglog/git-chglog) to manage the [CHANGELOG.md](CHANGELOG.md). For the
versions available, see the [tags on this repository](https://github.com/clok/sm/tags).

## Authors

* **Derek Smith** - [@clok](https://github.com/clok)

See also the list of [contributors](https://github.com/clok/sm/contributors) who participated in this project.

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details