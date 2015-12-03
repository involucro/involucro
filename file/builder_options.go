package file

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
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
		switch lua.CheckString(l, -2) {
		case "Hostname":
			conf.Hostname = lua.CheckString(l, -1)
		case "Domainname":
			conf.Domainname = lua.CheckString(l, -1)
		case "User":
			conf.User = lua.CheckString(l, -1)
		case "CPUSet":
			conf.CPUSet = lua.CheckString(l, -1)
		case "StopSignal":
			conf.StopSignal = lua.CheckString(l, -1)
		case "Image":
			conf.Image = lua.CheckString(l, -1)
		case "VolumeDriver":
			conf.VolumeDriver = lua.CheckString(l, -1)
		case "VolumesFrom":
			conf.VolumesFrom = lua.CheckString(l, -1)
		case "WorkingDir":
			conf.WorkingDir = lua.CheckString(l, -1)
		case "MacAddress":
			conf.MacAddress = lua.CheckString(l, -1)

		case "Memory":
			conf.Memory = int64(lua.CheckInteger(l, -1))
		case "MemorySwap":
			conf.MemorySwap = int64(lua.CheckInteger(l, -1))
		case "MemoryReservation":
			conf.MemoryReservation = int64(lua.CheckInteger(l, -1))
		case "KernelMemory":
			conf.KernelMemory = int64(lua.CheckInteger(l, -1))
		case "CPUShares":
			conf.CPUShares = int64(lua.CheckInteger(l, -1))

		case "AttachStdin":
			conf.AttachStdin = checkBoolean(l, -1)
		case "AttachStdout":
			conf.AttachStdout = checkBoolean(l, -1)
		case "AttachStderr":
			conf.AttachStderr = checkBoolean(l, -1)
		case "Tty":
			conf.Tty = checkBoolean(l, -1)
		case "OpenStdin":
			conf.OpenStdin = checkBoolean(l, -1)
		case "StdinOnce":
			conf.StdinOnce = checkBoolean(l, -1)
		case "NetworkDisabled":
			conf.NetworkDisabled = checkBoolean(l, -1)

		case "PortSpecs":
			conf.PortSpecs = checkStringArray(l, -1)
		case "Env":
			conf.Env = checkStringArray(l, -1)
		case "Cmd":
			conf.Cmd = checkStringArray(l, -1)
		case "DNS":
			conf.DNS = checkStringArray(l, -1)
		case "Entrypoint":
			conf.Entrypoint = checkStringArray(l, -1)
		case "SecurityOpts":
			conf.SecurityOpts = checkStringArray(l, -1)
		case "OnBuild":
			conf.OnBuild = checkStringArray(l, -1)
		case "Labels":
			conf.Labels = checkStringMap(l, -1)
		case "ExposedPorts":
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
