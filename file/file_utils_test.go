package file

import (
	. "github.com/smartystreets/goconvey/convey"
	duk "gopkg.in/olebedev/go-duktape.v2"
	"testing"
)

func TestRequireString(t *testing.T) {
	Convey("RequireString fails when asked to retrieve a number", t, func() {
		ctx := duk.New()
		ctx.PushNumber(5)
		So(func() { requireStringOrFailGracefully(ctx, -1, "test") }, ShouldPanic)
	})

	Convey("RequireString succeeds when asked to retrieve a string", t, func() {
		ctx := duk.New()
		ctx.PushString("asd")
		So(func() { requireStringOrFailGracefully(ctx, -1, "test") }, ShouldNotPanic)
	})
}
