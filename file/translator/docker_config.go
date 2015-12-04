package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"strings"
)

// ParseImageConfigFromLuaTable reads all keys in the currently top-most
// table from the stack and applies everything it can to a docker.Config.
// The comparison is case insensitive by design.
func ParseImageConfigFromLuaTable(l *lua.State) docker.Config {
	lua.CheckType(l, -1, lua.TypeTable)

	conf := docker.Config{}

	l.PushNil()
	for l.Next(-2) {
		switch strings.ToLower(lua.CheckString(l, -2)) {
		case "hostname":
			conf.Hostname = lua.CheckString(l, -1)
		case "domainname":
			conf.Domainname = lua.CheckString(l, -1)
		case "user":
			conf.User = lua.CheckString(l, -1)
		case "cpusetet":
			conf.CPUSet = lua.CheckString(l, -1)
		case "stopsignal":
			conf.StopSignal = lua.CheckString(l, -1)
		case "image":
			conf.Image = lua.CheckString(l, -1)
		case "volumedriver":
			conf.VolumeDriver = lua.CheckString(l, -1)
		case "volumesfrom":
			conf.VolumesFrom = lua.CheckString(l, -1)
		case "workingdir":
			conf.WorkingDir = lua.CheckString(l, -1)
		case "macaddress":
			conf.MacAddress = lua.CheckString(l, -1)

		case "memory":
			conf.Memory = int64(lua.CheckInteger(l, -1))
		case "memoryswap":
			conf.MemorySwap = int64(lua.CheckInteger(l, -1))
		case "memoryreservation":
			conf.MemoryReservation = int64(lua.CheckInteger(l, -1))
		case "kernelmemory":
			conf.KernelMemory = int64(lua.CheckInteger(l, -1))
		case "cpushares":
			conf.CPUShares = int64(lua.CheckInteger(l, -1))

		case "attachatdin":
			conf.AttachStdin = checkBoolean(l, -1)
		case "attachatdout":
			conf.AttachStdout = checkBoolean(l, -1)
		case "attachatderr":
			conf.AttachStderr = checkBoolean(l, -1)
		case "tty":
			conf.Tty = checkBoolean(l, -1)
		case "openstdin":
			conf.OpenStdin = checkBoolean(l, -1)
		case "stdinonce":
			conf.StdinOnce = checkBoolean(l, -1)
		case "networkdisabled":
			conf.NetworkDisabled = checkBoolean(l, -1)

		case "portspecs":
			conf.PortSpecs = checkStringArray(l, -1)
		case "env":
			conf.Env = checkStringArray(l, -1)
		case "cmd":
			conf.Cmd = checkStringArray(l, -1)
		case "dns":
			conf.DNS = checkStringArray(l, -1)
		case "entrypoint":
			conf.Entrypoint = checkStringArray(l, -1)
		case "securityopts":
			conf.SecurityOpts = checkStringArray(l, -1)
		case "onbuild":
			conf.OnBuild = checkStringArray(l, -1)
		case "labels":
			conf.Labels = checkStringMap(l, -1)
		case "exposedports", "ports", "expose":
			conf.ExposedPorts = parseExposedPorts(l, -1)
		case "volumes":
			conf.Volumes = checkStringSet(l, -1)
		}
		l.Pop(1)
	}

	return conf
}

/*
not implemented:
	Mounts            []Mount
*/
