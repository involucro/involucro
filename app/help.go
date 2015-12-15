package app

import "github.com/docopt/docopt-go"

func parseArguments(argv []string, exit bool) (map[string]interface{}, error) {

	usage := `Involucro.

Usage:
  involucro -h | --help
  involucro --version
  involucro [-w <path>] [ -H <url> | --host=<url> ] [-v [-v]] [-f <file> | -e <script>] [-s KEY=VALUE | --set KEY=VALUE]... [--] <task>...
  involucro -n [-v [-v]] [-f <file> | -e <script> ] <task>...
  involucro --wrap=<source-dir> --into-image=<parent-image> --at=<target-dir> --as=<image-id>

Options:
  -h --help               Show this screen.
  -H, --host=<url>        Set the URL for Docker [default: unix:///var/run/docker.sock].
  --version               Show version.
  -f <file>               Set the control file [default: invfile.lua].
  -e <script>             Evaluate the given script directly.
  -v                      Increase verbosity (use twice for even more messages).
  -w <path>               Set working dir, being the base for all scoping operations. [default: .]
  -s, --set KEY=VALUE     Makes VAR[KEY] available with value VALUE in the Lua script.
  --wrap=<source-dir>    
	--into-image=<parent-image>
	--at=<target-dir>
	--as=<image-id>
`
	return docopt.Parse(usage, argv, true, "Involucro 0.1", false, exit)
}
