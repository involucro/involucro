package file

import (
	"github.com/Shopify/go-lua"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestRequireString(t *testing.T) {
	Convey("RequireString fails when asked to retrieve a table", t, func() {
		ctx := lua.NewState()
		ctx.NewTable()
		So(func() { requireStringOrFailGracefully(ctx, -1, "test") }, ShouldPanic)
	})

	Convey("RequireString succeeds when asked to retrieve a string", t, func() {
		ctx := lua.NewState()
		ctx.PushString("asd")
		So(func() { requireStringOrFailGracefully(ctx, -1, "test") }, ShouldNotPanic)
	})
}
