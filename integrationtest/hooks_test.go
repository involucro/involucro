package integrationtest

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestIoLines(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	if err := ioutil.WriteFile("iolines_test.txt", []byte("ASD\nDSA"), 0755); err != nil {
		t.Fatal(err)
	}
	defer os.Remove("iolines_test.txt")

	script := `
for line in io.lines('iolines_test.txt') do
  inv.task('do' .. line)
	  .using('busybox')
			.run('echo', 'TESTOK')
end`

	cases := []string{"ASD", "DSA"}

	for _, el := range cases {
		assertStdoutContainsFlag([]string{
			"-e", script, "do" + el,
		}, "TESTOK", t)
	}
}

func TestHooksWithIoLines(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
	}
	ioutil.WriteFile("37_lines.txt", []byte{}, 0755)

	defer func() {
		os.Remove("37_lines.txt")
	}()
	script := `
inv.task('write')
  .using('busybox')
  .run('/bin/sh', '-c', 'echo A >> 37_lines.txt && echo B >> 37_lines.txt && echo C >> 37_lines.txt')

inv.task('modify')
  .hook(function ()
      for line in io.lines("37_lines.txt") do
        inv.task('do' .. line).using('busybox').run('echo', line)
      end
    end)
`

	assertStdoutContainsFlag([]string{"-e", script, "write", "modify", "doB"}, "B", t)
}
