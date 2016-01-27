package integrationtest

import "github.com/thriqon/involucro/app"

func ExampleTaskListWithTAndDirectScript() {
	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')",
		"-T",
	}); err != nil {
		panic(err)
	}
	// Output:
	// a
	// b
}

func ExampleTaskListWithTasksAndDirectScript() {
	if err := app.Main([]string{
		"involucro", "-e",
		"inv.task('a').using('busybox').run('x'); inv.task('b').using('busybox').run('z')",
		"--tasks",
	}); err != nil {
		panic(err)
	}
	// Output:
	// a
	// b
}
