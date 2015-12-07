package app

import "github.com/docopt/docopt-go"

func parseArguments(argv []string, exit bool) (map[string]interface{}, error) {

	usage := `Involucro.

Usage:
  involucro -h | --help
  involucro --version
  involucro [-w <path>] [ -H <url> | --host=<url> ] [-v [-v]] [-f <file> | -e <script>] [--] <task>...
  involucro (-n | -s) [-v [-v]] [-f <file> | -e <script> ] <task>...
  involucro --wrap=<source-dir> --into-image=<parent-image> --at=<target-dir> --as=<image-id>

Options:
  -h --help               Show this screen.
  -H, --host=<url>        Set the URL for Docker [default: unix:///var/run/docker.sock].
  --version               Show version.
  -f <file>               Set the control file [default: invfile.lua].
  -e <script>             Evaluate the given script directly.
  -v                      Increase verbosity (use twice for even more messages).
  -n                      Do not really execute commands in Docker, just show them.
  -s                      Instead of executing the commands against Docker, print equivalent shell commands.
  -w <path>               Set working dir, being the base for all scoping operations. [default: .]
  --wrap=<source-dir>    
	--into-image=<parent-image>
	--at=<target-dir>
	--as=<image-id>
`
	return docopt.Parse(usage, argv, true, "Involucro 0.1", false, exit)
}
