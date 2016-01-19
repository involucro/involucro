package app

import "fmt"

func ExampleConnectToDocker_SettingTcpUrl() {
	client, isremote, err := connectToDocker("tcp://fkdkd.de")
	fmt.Println(err, isremote, client.Endpoint())
	// Output: <nil> true tcp://fkdkd.de
}

func ExampleConnectToDocker_SettingUnixUrl() {
	client, isremote, err := connectToDocker("unix:///var/l.sock")
	fmt.Println(err, isremote, client.Endpoint())
	// Output: <nil> false unix:///var/l.sock
}
