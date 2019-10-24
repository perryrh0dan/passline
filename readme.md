<h1 align="center">
  Passline
</h1>

<h4 align="center">
  Password manager for the command line
</h4>

<div align="center">
  <img alt="List" width="70%" src="media/passline.gif">
</div>

## Description

By utilizing a simple and minimal usage syntax, that requires a flat leaning curve, Passline enables you to effectively manage you password accross multiple devices within your terminal. All Password are stored AES-256 encrypted and can only be encrypted with a your global password. Currently data can be stored localy on your computer or in your own firebase database.

## Highlights


## Contents

* [Description](#description)
* [Highlights](#highlights)
* [Install](#install)
* [Usage](#usage)

## Install
### Binary

### Snapcraft

``` bash
snap install passline
```

**Note:** Due to the snap's strictly confined nature, both the storage & configuration files will be saved under the [ `$SNAP_USER_DATA` ](https://docs.snapcraft.io/reference/env) environment variable instead of the generic `$HOME` one.

## Development
### Linter
golangci-lint
VS-Code settings
``` json
"go.lintTool":"golangci-lint",
"go.lintFlags": [
  "--fast"
]
```