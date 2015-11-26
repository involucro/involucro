package main

import "github.com/docopt/docopt-go"

func parse() map[string]interface{} {

	usage := `Involucro.

Usage:
  involucro -h | --help
  involucro --version
  involucro [ -H <url> | --host=<url> ] [-v [-v]] [ -f <file> ] [--] <task>...
	involucro (-n | -s) [-v [-v]] [-f <file>] [--] <task>...

Options:
  -h --help               Show this screen.
  -H, --host=<url>        Set the URL for Docker [default: unix:///var/run/docker.sock].
  --version               Show version.
  -f <file>               Set the control file [default: invfile.js].
  -v                      Increase verbosity (use twice for even more messages).
	-n                      Do not really execute commands in Docker, just show them.
	-s                      Instead of executing the commands against Docker, print equivalent shell commands.
`
	arguments, _ := docopt.Parse(usage, nil, true, "Involucro 0.1", false)
	return arguments
}
