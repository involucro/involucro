package app

import "flag"

var (
	dockerUrl       string
	controlFile     string
	controlScript   string
	logLevel        int
	relativeWorkDir string
	variables       variablesValue
	showTasks       bool

	remoteWrapTask string
)

const (
	defaultDockerUrl   = "unix:///var/run/docker.sock"
	defaultControlFile = "invfile.lua"
)

var flags *flag.FlagSet

func initializeFlagSet() {
	flags = flag.NewFlagSet("involucro", flag.ExitOnError)

	flags.StringVar(&dockerUrl, "H", defaultDockerUrl, "Set the URL of the Docker instance")
	flags.StringVar(&dockerUrl, "host", defaultDockerUrl, "Long form for -H")

	flags.StringVar(&controlFile, "f", defaultControlFile, "Set the control file")
	flags.StringVar(&controlScript, "e", "", "Evaluate the given script directly, not evaluating the control file")

	flags.IntVar(&logLevel, "l", -1, "Set minimum log level, -3 logs everything.")

	flags.StringVar(&relativeWorkDir, "w", ".", "Set working dir, being the base for all operations. Also settable via environment variable $INVOLUCRO_WORKDIR")

	flags.Var(&variables, "set", "Used as KEY=VALUE, makes VAR[KEY] available with value VALUE in Lua script")
	flags.Var(&variables, "s", "Shorthand for --set")

	flags.BoolVar(&showTasks, "tasks", false, "Show available tasks and then exit")
	flags.BoolVar(&showTasks, "T", false, "Shorthand for --tasks")

	flags.StringVar(&remoteWrapTask, "wrap", "", "Execute encoded wrap task")
}

func isControlFileOverriden() bool {
	return controlFile != defaultControlFile
}
