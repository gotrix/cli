# cli [![GoDoc](https://godoc.org/github.com/gotrix/cli?status.svg)](http://godoc.org/github.com/gotrix/cli) [![Go Report](https://goreportcard.com/badge/github.com/gotrix/cli)](https://goreportcard.com/report/github.com/gotrix/cli)
GOTRIX command line helper.

## Installation
```
go get -u github.com/gotrix/cli/...
```

## Usage
```
Usage:

  gotrix [options] command

Commands:

  create-app
        Create a new application.

  create-component
        Create a new component.

  build-components
        Build components so files. Path can be modified by -path
        option. The defaut is "./components".

Options:

  -help
        Show this help.

  -no-color
        Do not colorize output (default false).

  -path
        Path to contents directory (default depends on the command).

  -quiet
        Do not print any output (default false).
```

## Commands
List of available cli commands.

### create-app
Create a new application.

### create-component
Create a new component.

### build-components
Build components so files. Path can be modified by -path option. The defaut is "./components".

## Arguments

### -help
Show help.

### -no-color
Do not colorize output (default false).

### -path
Path to contents directory (default depends on the command).

### -quiet
Do not print any output (default false).
