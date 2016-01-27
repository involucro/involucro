package ilog

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/fatih/color"
)

var (
	StdLogger = New()
	Debug     = ForLevelPrefix(-2, "DEBU")
	Info      = ForLevelPrefix(-1, "INFO")
	Warn      = ForLevelPrefix(0, "WARN")
	Error     = ForLevelPrefix(1, "ERRO")
)

type Bough struct {
	Level   int
	Prefix  string
	Message string
}

type PrintFunc func(b Bough)

type Ilog struct {
	mut sync.Mutex // protects all below

	print         PrintFunc
	minPrintLevel int
}

type Logger interface {
	Level() int
	Logln(a ...interface{})
	Logf(f string, a ...interface{})
}

func DefaultPrintFunc(b Bough) {
	fmt.Fprintf(os.Stderr, "%s  %s\n", b.Prefix, b.Message)
}

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

func New() *Ilog {
	return &Ilog{}
}

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

func (i *Ilog) SetPrintFunc(f PrintFunc) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.print = f
}

func (i *Ilog) PrintFunc() PrintFunc {
	i.mut.Lock()
	defer i.mut.Unlock()
	return i.print
}

func (i *Ilog) SetMinPrintLevel(level int) {
	i.mut.Lock()
	defer i.mut.Unlock()
	i.minPrintLevel = level
}

func (i *Ilog) MinPrintLevel() int {
	i.mut.Lock()
	defer i.mut.Unlock()
	return i.minPrintLevel
}

type bougher struct {
	i      *Ilog
	level  int
	prefix string
}

func (b bougher) Logln(a ...interface{}) {
	b.i.Send(Bough{b.level, b.prefix, fmt.Sprint(a...)})
}

func (b bougher) Logf(f string, a ...interface{}) {
	b.i.Send(Bough{b.level, b.prefix, fmt.Sprintf(f, a...)})
}

func (b bougher) Level() int {
	return b.level
}

func (i *Ilog) ForLevelPrefix(l int, prefix string) Logger {
	return bougher{i, l, prefix}
}

func ForLevelPrefix(l int, prefix string) Logger {
	return StdLogger.ForLevelPrefix(l, prefix)
}
