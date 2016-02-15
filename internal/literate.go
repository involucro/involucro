package runtime

// Literate Programming is a technique invented by Donald Knuth. It allows the
// mixture of prose and code, allowing the author/developer to express the
// chain of thoughts that went into the code.
//
// One common approach is use a file containing Markdown source with code
// comments. These can be interpreted by compatible tools and used just as
// "normal" code.
//
// This file provides the possibility for the invfile's to be written in such
// way. This is useful since they often are the result of long and
// time-consuming development efforts.

import (
	"bufio"
	"bytes"
	"os"
	"strings"

	"github.com/Shopify/go-lua"
)

// RunLiterateFile interprets the given file as literate Markdown with
// intermixed Lua code and runs that code.
func (inv *Runtime) RunLiterateFile(filename string) error {

	// Using the supplied filename, we open the file for reading.
	fileReader, err := os.Open(filename)
	if err != nil {
		return err
	}

	// We construct a scanner around the file reader.  Note that the default
	// split algorithm for scanners is splitting by line (which is what we need).
	scanner := bufio.NewScanner(fileReader)

	// The easiest thing to do is to keep all the source code lines in memory
	// until they can be handed of to Lua for execution. Another possibility is
	// to use a Goroutine for loading the lines and filtering them, and using a
	// pipe to send the interesting ones to Lua.

	lines := bytes.Buffer{}

	// Space-prefixed rows need an empty row directly before them. We record if
	// the precedeing row was empty, and set this variable at the end of every
	// loop. If the file just starts, we also have to assume the line before tat
	// was empty.
	preceedingRowEmtpy := true

	// the Scan() method returns a boolean indicating whether a next line is
	// available.
	for scanner.Scan() {
		// This line can be retrieved using Text(). It does not contain the line
		// ending character.
		line := scanner.Text()

		var trimmed string
		// This check tests whether this line starts with a bird track prefix: "> "
		// or with four spaces. If it does, the rest of the lines is stored in the
		// lines buffer.
		switch {
		case strings.HasPrefix(line, "> "):
			trimmed = strings.TrimPrefix(line, "> ")
		case strings.HasPrefix(line, "    ") && preceedingRowEmtpy:
			trimmed = strings.TrimPrefix(line, "    ")
		case line == "":
			preceedingRowEmtpy = true
			continue
		default:
			preceedingRowEmtpy = false
			continue
		}

		// Add the chopped line ending character.
		trimmed += "\n"

		// Write the result into the buffer
		lines.Write([]byte(trimmed))
	}

	// The scanner may have received an error, which we have to check.
	if err := scanner.Err(); err != nil {
		return err
	}

	// LoadBuffer parses the given string and pushes an anonymous, parameterless
	// function to the top of the stack. Errors are reported using the filename
	// in the third parameter.
	if err := lua.LoadBuffer(inv.lua, string(lines.Bytes()), filename, ""); err != nil {
		return err
	}

	// This anonymous function can be called with ProtectedCall. ProtectedCall
	// doesn't panic if an error  occurs, and instead returns the error as-is.
	return inv.lua.ProtectedCall(0, lua.MultipleReturns, 0)
}
