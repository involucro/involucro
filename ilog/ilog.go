// Package ilog provides a simple and lightweight logging framework for
// involucro. It is a Go-lang-y abbreviation for involucro logging.
//
// All public methods in this package are thread-safe.
package ilog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

// StdLog is the default logger that is availabe in the global space.
var StdLog = New()

// Logger provides functionality to log strings. It provides a ln variant which
// works like Sprint, and a f variant like Sprintf.
type Logger interface {
	level() int
	Logln(a ...interface{})
	Logf(f string, a ...interface{})
}

var (
	// Debug is a level not printed by default with prefix "DEBU".
	Debug = ForLevelPrefix(-2, "DEBU")

	// Info is a level not printed by default with prefix "INFO".
	Info = ForLevelPrefix(-1, "INFO")

	// Warn is a level printed by default with prefix "WARN".
	Warn = ForLevelPrefix(0, "WARN")

	// Error is a level printed by default with prefix "ERRO".
	Error = ForLevelPrefix(1, "ERRO")
)

// A Bough is an entry in a log.
type Bough struct {
	Level   int
	Prefix  string
	Message string
}

// PrintFunc is a function that is supposed to handle a bough, for example
// printing it on the terminal.
type PrintFunc func(b Bough)

// ForLevelPrefix gives a Logger that logs with level l and prefix p. It logs
// on the default logging context.
func ForLevelPrefix(l int, prefix string) Logger {
	return StdLog.ForLevelPrefix(l, prefix)
}

// Ilog is the context for all loggers. There is a default instance which is
// sufficient for most use cases.
type Ilog struct {
	mut sync.Mutex // protects all below

	print         PrintFunc
	minPrintLevel int
}

// ColorfulPrintFunc is a PrintFunc that formats the Bough nicely with colors
// (if supported).
func ColorfulPrintFunc(b Bough) {
	prefixPrinter := fmt.Sprintf
	switch {
	case b.Level == 0:
		prefixPrinter = color.YellowString
	case b.Level > 0:
		prefixPrinter = color.RedString
	case b.Level == -1:
		prefixPrinter = color.BlueString
	}

	fmt.Fprintf(os.Stderr, "[%s] %s %s\n", time.Now().Format(time.Stamp), prefixPrinter(b.Prefix), b.Message)
}

// New creates a new context.
func New() *Ilog {
	return &Ilog{}
}

// Send handles a Bough and delivers it to the print function if the level of
// the bough is at least the minimum print level.
func (i *Ilog) Send(b Bough) {
	i.mut.Lock()
	defer i.mut.Unlock()

	if b.Level >= i.minPrintLevel {
		f := i.print
		if f == nil {
			f = ColorfulPrintFunc
		}
		f(b)
	}
}

// SetPrintFunc replaces the print function with f.
func (i *Ilog) SetPrintFunc(f PrintFunc) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.print = f
}

// PrintFunc gives the function currently used for printing.
func (i *Ilog) PrintFunc() PrintFunc {
	i.mut.Lock()
	defer i.mut.Unlock()
	return i.print
}

// SetMinPrintLevel sets the minimum required level for a Bough to be actually
// printed. Level 2 messages are only printed if the level is 2 or lower.
func (i *Ilog) SetMinPrintLevel(level int) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.minPrintLevel = level
}

// MinPrintLevel returns the minimum required print level.
func (i *Ilog) MinPrintLevel() int {
	i.mut.Lock()
	defer i.mut.Unlock()
	return i.minPrintLevel
}

// bougher is an implementation of Logger. It stores all relevant information
// in itself.
type bougher struct {
	i      *Ilog
	l      int
	prefix string
}

func (b bougher) Logln(a ...interface{}) {
	b.i.Send(Bough{b.l, b.prefix, fmt.Sprint(a...)})
}

func (b bougher) Logf(f string, a ...interface{}) {
	b.i.Send(Bough{b.l, b.prefix, fmt.Sprintf(f, a...)})
}

func (b bougher) level() int {
	return b.l
}

// ForLevelPrefix gives a logger for level l and prefix p.
func (i *Ilog) ForLevelPrefix(l int, prefix string) Logger {
	return bougher{i, l, prefix}
}
