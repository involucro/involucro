package translator

import (
	"strings"

	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"github.com/involucro/involucro/ilog"
)

// ParseHostConfigFromLuaTable reads all keys in the currently top-most
// table from the stack and applies everything it can to the given default
// docker.HostConfig.
// The comparison is case insensitive by design.
func ParseHostConfigFromLuaTable(l *lua.State, conf docker.HostConfig) docker.HostConfig {
	lua.CheckType(l, -1, lua.TypeTable)

	l.PushNil()
	for l.Next(-2) {
		switch strings.ToLower(lua.CheckString(l, -2)) {
		case "binds":
			conf.Binds = checkStringArray(l, -1)
		case "capadd":
			conf.CapAdd = checkStringArray(l, -1)
		case "capdrop":
			conf.CapDrop = checkStringArray(l, -1)
		case "groupadd":
			conf.GroupAdd = checkStringArray(l, -1)
		case "dns":
			conf.DNS = checkStringArray(l, -1)
		case "dnsoptions":
			conf.DNSOptions = checkStringArray(l, -1)
		case "dnssearch":
			conf.DNSSearch = checkStringArray(l, -1)
		case "extrahosts":
			conf.ExtraHosts = checkStringArray(l, -1)
		case "volumesfrom":
			conf.VolumesFrom = checkStringArray(l, -1)
		case "links":
			conf.Links = checkStringArray(l, -1)
		case "securityopt":
			conf.SecurityOpt = checkStringArray(l, -1)

		case "containeridfile":
			conf.ContainerIDFile = lua.CheckString(l, -1)
		case "networkmode":
			conf.NetworkMode = lua.CheckString(l, -1)
		case "ipcmode":
			conf.IpcMode = lua.CheckString(l, -1)
		case "pidmode":
			conf.PidMode = lua.CheckString(l, -1)
		case "utsmode":
			conf.UTSMode = lua.CheckString(l, -1)
		case "cgroupparent":
			conf.CgroupParent = lua.CheckString(l, -1)
		case "cpuset":
			conf.CPUSet = lua.CheckString(l, -1)
		case "cpusetcpus":
			conf.CPUSetCPUs = lua.CheckString(l, -1)
		case "cpusetmems":
			conf.CPUSetMEMs = lua.CheckString(l, -1)
		case "volumedriver":
			conf.VolumeDriver = lua.CheckString(l, -1)

		case "cpushares":
			conf.CPUShares = int64(lua.CheckInteger(l, -1))
		case "cpuquota":
			conf.CPUQuota = int64(lua.CheckInteger(l, -1))
		case "cpuperiod":
			conf.CPUPeriod = int64(lua.CheckInteger(l, -1))
		case "blkioweight":
			conf.BlkioWeight = int64(lua.CheckInteger(l, -1))
		case "memory":
			conf.Memory = int64(lua.CheckInteger(l, -1))
		case "memoryswap":
			conf.MemorySwap = int64(lua.CheckInteger(l, -1))
		case "memoryswappiness":
			conf.MemorySwappiness = int64(lua.CheckInteger(l, -1))

		case "privileged":
			conf.Privileged = checkBoolean(l, -1)
		case "publishallports":
			conf.PublishAllPorts = checkBoolean(l, -1)
		case "readonlyrootfs":
			conf.ReadonlyRootfs = checkBoolean(l, -1)
		case "oomkilldisable":
			conf.OOMKillDisable = checkBoolean(l, -1)
		default:
			ilog.Warn.Logf("Unrecognized setting [%s] in config, ignoring", lua.CheckString(l, -2))
		}
		l.Pop(1)
	}

	return conf
}

/*
LxcConf          []KeyValuePair
PortBindings     map[Port][]PortBinding
RestartPolicy    RestartPolicy
Devices          []Device
LogConfig        LogConfig
Ulimits          []ULimit
*/
