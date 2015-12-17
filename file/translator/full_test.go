package translator

import (
	"github.com/Shopify/go-lua"
	"github.com/fsouza/go-dockerclient"
	"reflect"
	"testing"
)

func TestAllPropertiesConfig(t *testing.T) {
	source := `x = {
		hostname = "host",
		domainname = "dom",
		user = "nobody",
		cpuset = "1",
		stopsignal = "SIGTERM",
		image = "28208",
		volumedriver = "flocker",
		volumesfrom = "data",
		workingdir = "/tmp",
		macaddress = "00:12",

		memory = 100,
		memoryswap = 200,
		memoryreservation = 2010,
		kernelmemory = 201,
		cpushares = 2,

		attachstdin = true,
		attachstdout = true,
		attachstderr = true,
		tty = true,
		openstdin = true,
		stdinonce = true,
		networkdisabled = true,

		portspecs = {"2", "3"},
		env = {"FOO=bar"},
		cmd = {"a", "b", "c"},
		dns = {"8.8.8.8"},
		entrypoint = {"/bin/sh"},
		securityopts = {"a", "b"},
		onbuild = {"RUN echo test"},
		labels = {
			example = "asd"
		},
		exposedports = {"80/tcp"},
		volumes = {"/data"}
	}`

	expected := docker.Config{
		Hostname:     "host",
		Domainname:   "dom",
		User:         "nobody",
		CPUSet:       "1",
		StopSignal:   "SIGTERM",
		Image:        "28208",
		VolumeDriver: "flocker",
		VolumesFrom:  "data",
		WorkingDir:   "/tmp",
		MacAddress:   "00:12",

		Memory:            int64(100),
		MemorySwap:        int64(200),
		MemoryReservation: int64(2010),
		KernelMemory:      int64(201),
		CPUShares:         2,

		AttachStdin:     true,
		AttachStdout:    true,
		AttachStderr:    true,
		Tty:             true,
		OpenStdin:       true,
		StdinOnce:       true,
		NetworkDisabled: true,

		PortSpecs:    []string{"2", "3"},
		Env:          []string{"FOO=bar"},
		Cmd:          []string{"a", "b", "c"},
		DNS:          []string{"8.8.8.8"},
		Entrypoint:   []string{"/bin/sh"},
		SecurityOpts: []string{"a", "b"},
		OnBuild:      []string{"RUN echo test"},
		Labels: map[string]string{
			"example": "asd",
		},
		ExposedPorts: map[docker.Port]struct{}{
			docker.Port("80/tcp"): struct{}{},
		},
		Volumes: map[string]struct{}{
			"/data": struct{}{},
		},
	}

	state := lua.NewState()
	if err := lua.DoString(state, source); err != nil {
		t.Fatal("Error during code execution", err)
	}
	state.Global("x")

	if actual := ParseImageConfigFromLuaTable(state); !reflect.DeepEqual(actual, expected) {
		t.Error("Actual is not equal to expected", actual, expected)
	}
}

func TestAllPropertiesHostConfig(t *testing.T) {
	source := `x = {
			binds = {"/data:asd"},
			capadd = {"CAP_ROOT"},
			capdrop = {"CAP_NOTROOT"},
			GroupAdd = {"users"},
			DNS = {"8.8.8.8"},
			dnsoptions = {"-more-opts"},
			dnssearch = {"a.dns", "b.dns"},
			extrahosts = {"127.0.0.1 service"},
			volumesfrOm = {"data"},
			links = {"db"},
			securityopt = {"asd"},

			containeridfile = "/var/run/pid",
			networkmode = "NAT",
			ipcmode = "0755",
			pidmode = "0755",
			utsmode = "0755",
			cgroupparent = "test",
			cpuset = "2",
			cpusetcpus = "2",
			cpusetmems = "3",
			volumedriver = "flocker",

			cpushares = 82,
			cpuquota = 298,
			cpuperiod = 29,
			blkioweight = 291,
			memory = 291,
			memoryswap = 292,
			memoryswappiness = 293,

			privileged = true,
			publishAllPorts = true,
			ReadOnlyROOTFs = true,
			OomKillDisable = true
		}`

	expected := docker.HostConfig{
		Binds:       []string{"/data:asd"},
		CapAdd:      []string{"CAP_ROOT"},
		CapDrop:     []string{"CAP_NOTROOT"},
		GroupAdd:    []string{"users"},
		DNS:         []string{"8.8.8.8"},
		DNSOptions:  []string{"-more-opts"},
		DNSSearch:   []string{"a.dns", "b.dns"},
		ExtraHosts:  []string{"127.0.0.1 service"},
		VolumesFrom: []string{"data"},
		Links:       []string{"db"},
		SecurityOpt: []string{"asd"},

		ContainerIDFile: "/var/run/pid",
		NetworkMode:     "NAT",
		IpcMode:         "0755",
		PidMode:         "0755",
		UTSMode:         "0755",
		CgroupParent:    "test",
		CPUSet:          "2",
		CPUSetCPUs:      "2",
		CPUSetMEMs:      "3",
		VolumeDriver:    "flocker",

		CPUShares:        int64(82),
		CPUQuota:         int64(298),
		CPUPeriod:        int64(29),
		BlkioWeight:      int64(291),
		Memory:           int64(291),
		MemorySwap:       int64(292),
		MemorySwappiness: int64(293),

		Privileged:      true,
		PublishAllPorts: true,
		ReadonlyRootfs:  true,
		OOMKillDisable:  true,
	}
	state := lua.NewState()
	if err := lua.DoString(state, source); err != nil {
		t.Fatal("Error during code execution", err)
	}
	state.Global("x")

	if actual := ParseHostConfigFromLuaTable(state); !reflect.DeepEqual(actual, expected) {
		t.Error("Actual is not equal to expected", actual, expected)
	}
}
