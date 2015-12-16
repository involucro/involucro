package app

import "fmt"

func ExampleConnectToDocker_WithoutArguments() {
	client, isremote, err := connectToDocker(argumentsMap{})
	fmt.Println(err, isremote, client.Endpoint())
	// Output: <nil> false unix:///var/run/docker.sock
}

func ExampleConnectToDocker_SettingTcpUrl() {
	client, isremote, err := connectToDocker(argumentsMap{"--host": "tcp://fkdkd.de"})
	fmt.Println(err, isremote, client.Endpoint())
	// Output: <nil> true tcp://fkdkd.de
}

func ExampleConnectToDocker_SettingUnixUrl() {
	client, isremote, err := connectToDocker(argumentsMap{"--host": "unix:///var/l.sock"})
	fmt.Println(err, isremote, client.Endpoint())
	// Output: <nil> false unix:///var/l.sock
}
