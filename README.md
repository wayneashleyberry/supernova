> Supernova is a Go command that backs up your GitHub stars offline and removes them from GitHub.com.

I've been on a privacy binge lately, and having amassed thousands of GitHub stars over the years I did not like all of that data being public. This script provides simple helpers to read your list of stars, so that you can back them up, and then remove them all from GitHub.

### Installation

```sh
go install github.com/wayneashleyberry/supernova@latest
```

### Usage

```
Usage:
  supernova [flags]
  supernova [command]

Available Commands:
  delete      Unstar everything on GitHub
  env         Print the required environment variables
  help        Help about any command
  read        Print a list of your GitHub stars

Flags:
  -h, --help   help for stars

Use "stars [command] --help" for more information about a command.
```
