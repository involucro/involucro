package main

import (
	"flag"
	"os"

	"github.com/involucro/involucro/app"
	"github.com/involucro/involucro/ilog"
)

func main() {
	err := app.Main(os.Args)

	switch err {
	case flag.ErrHelp:
		os.Exit(0)
	case nil:
		os.Exit(0)
	default:
		ilog.Error.Logf("Task processing failed: %s", err)
		os.Exit(1)
	}
}
