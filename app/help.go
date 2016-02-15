package app

import "flag"

var (
	dockerURL       string
	controlFile     string
	controlScript   string
	logLevel        int
	relativeWorkDir string
	variables       variablesValue
	showTasks       bool
	showVersion     bool

	remoteWrapTask string
)

var version = "vDev"

const (
	defaultDockerURL   = "unix:///var/run/docker.sock"
	defaultControlFile = "invfile.lua"
)

var flags *flag.FlagSet

func initializeFlagSet() {
	flags = flag.NewFlagSet("involucro", flag.ExitOnError)

	flags.StringVar(&dockerURL, "H", defaultDockerURL, "Set the URL of the Docker instance")
	flags.StringVar(&dockerURL, "host", defaultDockerURL, "Long form for -H")

	flags.StringVar(&controlFile, "f", defaultControlFile, "Set the control file")
	flags.StringVar(&controlScript, "e", "", "Evaluate the given script directly, not evaluating the control file")

	flags.IntVar(&logLevel, "l", -1, "Set minimum log level, -3 logs everything.")

	flags.StringVar(&relativeWorkDir, "w", ".", "Set working dir, being the base for all operations. Also settable via environment variable $INVOLUCRO_WORKDIR")

	flags.Var(&variables, "set", "Used as KEY=VALUE, makes VAR[KEY] available with value VALUE in Lua script")
	flags.Var(&variables, "s", "Shorthand for --set")

	flags.BoolVar(&showTasks, "tasks", false, "Show available tasks and then exit")
	flags.BoolVar(&showTasks, "T", false, "Shorthand for --tasks")

	flags.StringVar(&remoteWrapTask, "wrap", "", "Execute encoded wrap task")

	flags.BoolVar(&showVersion, "version", false, "Show version and the exit")
}

func isControlFileOverriden() bool {
	return controlFile != defaultControlFile
}

var versionNotice = version + `

Licensed under the Apache 2.0 License, full text available at http://www.apache.org/licenses/LICENSE-2.0

This software contains (all licensed under the MIT License):

https://github.com/fatih/color            by Fatih Arslan
https://github.com/fsouza/go-dockerclient by Francisco Souza and contributors
https://github.com/mattn/go-colorable     by Yasuhiro Matsumoto
https://github.com/mattn/go-isatty        by Yasuhiro Matsumoto
https://github.com/Shopify/go-lua         by Shopify Inc
`
