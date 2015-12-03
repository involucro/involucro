package file

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"strings"
)

func (wbs wrapBuilderState) withConfig(l *lua.State) int {
	wbs.baseConf = parseImageConfigFromLuaTable(l)
	return wrapTable(l, &wbs)
}

func parseImageConfigFromLuaTable(l *lua.State) docker.Config {
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
		}
		l.Pop(1)
	}

	return conf
}

/*
	Volumes           map[string]struct{}
	Mounts            []Mount
*/

func checkBoolean(l *lua.State, index int) bool {
	lua.CheckType(l, index, lua.TypeBoolean)
	return l.ToBoolean(index)
}

func checkStringArray(l *lua.State, index int) []string {
	lua.CheckType(l, index, lua.TypeTable)

	items := make([]string, 0)

	l.PushNil()
	for l.Next(-2) {
		items = append(items, lua.CheckString(l, -1))
		l.Pop(1)
	}

	return items
}

func checkStringMap(l *lua.State, index int) map[string]string {
	lua.CheckType(l, index, lua.TypeTable)
	items := make(map[string]string)

	l.PushNil()
	for l.Next(-2) {
		key := lua.CheckString(l, -2)
		val := lua.CheckString(l, -1)
		items[key] = val
		l.Pop(1)
	}

	return items
}

type emptyStruct struct {
}

func parseExposedPorts(l *lua.State, index int) map[docker.Port]struct{} {
	items := map[docker.Port]struct{}{}

	marker := emptyStruct{}
	for _, el := range checkStringArray(l, index) {
		items[docker.Port(el)] = marker
	}

	return items
}
