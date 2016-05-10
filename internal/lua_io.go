package runtime

import (
	"bufio"
	"github.com/Shopify/go-lua"
	"os"
	goruntime "runtime"
)

func injectIoLib(l *lua.State) {
	tableWith(l, fm{
		//"close":   ioClose,
		//flush":   ioFlush,
		//"input":   ioInput,
		"lines": ioLines,
		//"open":    ioOpen,
		//"output":  ioOutput,
		//"popen":   ioPopen,
		//"read":    ioRead,
		//"tmpfile": ioTmpfile,
		//"type":    ioType,
		//"write":   ioWrite,
	})
	l.SetGlobal("io")
}

type ioFile struct {
	closed bool
	h      *os.File
	sc     *bufio.Scanner
}

func (iof ioFile) Close() error {
	if iof.closed {
		return nil
	}
	iof.closed = true
	return iof.h.Close()
}

func ioLines(l *lua.State) int {
	filename := lua.CheckString(l, -1)

	file, err := os.Open(filename)
	if err != nil {
		lua.Errorf(l, "Unable to open file: %s", err.Error())
		panic("unreachable")
	}

	fh := ioFile{
		h:  file,
		sc: bufio.NewScanner(file),
	}
	goruntime.SetFinalizer(&fh, func(iof *ioFile) {
		iof.Close()
	})

	l.PushGoFunction(fh.readline)
	return 1
}

func (iof *ioFile) readline(l *lua.State) int {
	if !iof.sc.Scan() {
		iof.Close()
		l.PushNil()
		err := iof.sc.Err()
		if err != nil {
			l.PushString(err.Error())
			return 2
		}
		return 1
	}
	l.PushString(iof.sc.Text())
	return 1
}
