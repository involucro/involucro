package app

import (
	"flag"
	"fmt"
	"strings"
)

var (
	dockerUrl       string
	controlFile     string
	controlScript   string
	verbose         bool
	silent          bool
	relativeWorkDir string
	variables       variablesValue
	showTasks       bool

	remoteWrapTask string
)

const (
	defaultDockerUrl   = "unix:///var/run/docker.sock"
	defaultControlFile = "invfile.lua"
)

func init() {
	flag.StringVar(&dockerUrl, "H", defaultDockerUrl, "Set the URL of the Docker instance")
	flag.StringVar(&dockerUrl, "host", defaultDockerUrl, "Long form for -H")

	flag.StringVar(&controlFile, "f", defaultControlFile, "Set the control file")
	flag.StringVar(&controlScript, "e", "", "Evaluate the given script directly, not evaluating the control file")

	flag.BoolVar(&silent, "silent", false, "Descrease verbosity, only shows errors")
	flag.BoolVar(&verbose, "v", false, "Log debug output messages")

	flag.StringVar(&relativeWorkDir, "w", ".", "Set working dir, being the base for all operations")

	flag.Var(&variables, "set", "Used as KEY=VALUE, makes VAR[KEY] available with value VALUE in Lua script")
	flag.Var(&variables, "s", "Shorthand for --set")

	flag.BoolVar(&showTasks, "tasks", false, "Show available tasks and then exit")
	flag.BoolVar(&showTasks, "T", false, "Shorthand for --tasks")

	flag.StringVar(&remoteWrapTask, "wrap", "", "Execute encoded wrap task")
}

type variablesValue map[string]string

func (v *variablesValue) String() string {
	items := make([]string, 0)
	for k, v := range *v {
		items = append(items, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("[%s]", strings.Join(items, " "))
}

type ErrInvalidFormatForVariableAssignment string

func (e ErrInvalidFormatForVariableAssignment) Error() string {
	return fmt.Sprintf("Invalid value [%s], expected value of the form: KEY=VALUE", string(e))
}

func (v *variablesValue) Set(s string) error {
	if *v == nil {
		*v = make(map[string]string)
	}
	ss := strings.SplitN(s, "=", 2)
	if len(ss) < 2 {
		return ErrInvalidFormatForVariableAssignment(s)
	}
	(*v)[ss[0]] = ss[1]
	return nil
}

func isControlFileOverriden() bool {
	return controlFile != defaultControlFile
}
