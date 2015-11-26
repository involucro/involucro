package main

import (
	. "github.com/smartystreets/goconvey/convey"
	duk "gopkg.in/olebedev/go-duktape.v2"
	"testing"
)

func TestTaskDefinition(t *testing.T) {
	var env *duk.Context

	Convey("Given an empty runtime environment", t, func() {
		env = InstantiateRuntimeEnv()

		Convey("Defining an empty Task succeeds", func() {
			So(func() { env.EvalString(`inv.task('test')`) }, ShouldNotPanic)
		})

		Convey("Defining a Task with using/run succeeds", func() {
			So(func() { env.EvalString(`inv.task('test').using('blah').run('test')`) }, ShouldNotPanic)
		})

		Convey("Defining a Task with a number as ID panics", func() {
			So(func() { env.EvalString(`inv.task(5)`) }, ShouldPanic)
		})
	})
}
