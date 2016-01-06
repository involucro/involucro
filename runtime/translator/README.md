
## file/translator

This package translates Lua tables into native Go data structures. It is used
for `Config` and `HostConfig`.

### API

`func ParseImageConfigFromLuaTable(l *lua.State) docker.Config`

`ParseImageConfigFromLuaTable` reads all keys in the currently top-most
table from the stack and applies everything it can to a `docker.Config`.
The comparison is case insensitive by design.


------

`func ParseHostConfigFromLuaTable(l *lua.State) docker.HostConfig`

`ParseHostConfigFromLuaTable` reads all keys in the currently top-most
table from the stack and applies everything it can to a `docker.HostConfig`.
The comparison is case insensitive by design.

