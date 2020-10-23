<h1 align="center">
  Passline
</h1>

<h4 align="center">
  Password manager for the command line
</h4>

<div align="center">
  <img alt="List" width="70%" src="media/passline.gif">
</div>

<div align="center">
   <a href="https://travis-ci.com/perryrh0dan/passline">
      <img alt="Build Status" src="https://travis-ci.org/perryrh0dan/passline.svg?branch=master" />
   </a>
   <a href="https://codecov.io/gh/perryrh0dan/passline">
      <img src="https://codecov.io/gh/perryrh0dan/passline/branch/master/graph/badge.svg" />
   </a>
   <a href="https://codeclimate.com/github/perryrh0dan/passline/maintainability">
      <img src="https://api.codeclimate.com/v1/badges/83561b59422e2492f9db/maintainability" />
   </a>
   <a href="https://gitter.im/perryrh0danpassline/community">
      <img alt="Build Status" src="https://badges.gitter.im/community.svg" />
   </a>
</div>

## Description

Passline is a command line-based password management system. Thanks to its simple and minimal usage syntax, Passline enables users to effectively manage various passwords across multiple devices within the terminal. All Passwords are stored AES-256 encrypted and can only be encrypted with a global password. Currently data can be stored localy on your computer or in your own firebase database.

Visit the [contributing guidelines](https://github.com/perryrh0dan/passline/blob/master/contributing.md#translating-documentation) to learn more on how to translate this document into more languages.

Come over to [Gitter](https://gitter.im/perryrh0danpassline/community?source=orgpage) or [Twitter](https://twitter.com/perryrh0dan1) to share your thoughts on the project.

## Highlights

- Multiple storage modules (local, firestore)
- Passwords and recovery codes are AES-256 encryped
- Intuitive and fast command line interface
- Filtering allows fast selection of credentials
- Built-in update functionality

## Contents

- [Description](#description)
- [Highlights](#highlights)
- [Contents](#contents)
- [Install](#install)
- [Usage](#usage)
- [Flight Manual](#flight-manual)
- [Development](#development)
- [Team](#team)
- [License](#license)

## Install

### Binary

1. Download the latest release for your platform.
2. Unzip the tar.gz files before running. 
3. Click on the passline.exe file to start the application. 

### Snapcraft

```bash
snap install passline
snap alias passline pl # set alias
```

**Note:** Due to the Snapcraft's strictly confined nature, both storage & configuration files will be saved under the [ `$SNAP_USER_DATA` ](https://docs.snapcraft.io/reference/env) environment variable instead of the generic `$HOME` one.

## Usage

```
> passline --help
NAME:
   Passline - Password manager

USAGE:
   passline [global options] command [command options] [arguments...]

VERSION:
   0.5.3

DESCRIPTION:
   Password manager for the command line

COMMANDS:
   add, a       Adds an existing password for a website
   backup, b    Creates a backup
   delete, d    Deletes an item
   edit, e      Edits an item
   generate, g  Generates a password for an website
   list, ls     Lists all websites/passwords
   password, p  Changes master password
   restore, r   Restores a backup
   update, u    Updates to the newest release
   help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   -- print, -p   Displays password on the terminal (default: false)
   --help, -h     Shows help (default: false)
   --version, -v  Prints the version (default: false)

WEBSITE:
   https://github.com/perryrh0dan/passline
```

### Flight Manual

#### Add

Adds an existing username, password pair.

#### Backup

Creates a backup. 

#### Generate

Generates and stores a password for a website/username combination.

#### Delete

Deletes an existing username, password pair.

#### Edit

Edits an entry

#### List

Lists all websites

#### Restore

Restores a backup file.

#### Update

Updates passline to the latest release on github.

## Development

### Linter

golangci-lint
VS-Code settings

```json
"go.lintTool":"golangci-lint",
"go.lintFlags": [
  "--fast"
]
```

### Build

``` bash
GOOS=windows GOARCH=amd64 go build
```

### Test

``` bash
go test ./...
```

### Icon

Icon is stored under notify/icon.go in dec.
Use this [tool](https://tomeko.net/online_tools/file_to_hex.php?lang=en) to easy convert file to hex and then to dec.

## Team

- Thomas Pöhlmann [(@perryrh0dan)](https://github.com/perryrh0dan)

## License

[MIT](https://github.com/perryrh0dan/passline/blob/master/license.md)

This repository was generated by [tmpo](https://github.com/perryrh0dan/tmpo)