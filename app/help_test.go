package app

import "testing"

func TestParseArguments(t *testing.T) {
	args, err := parseArguments([]string{"--encoded-state", "--socket", "/asd"}, false)

	if err != nil {
		t.Fatal("Failed parsing arguments")
	}

	if !args["--encoded-state"].(bool) {
		t.Fatal("Didn't parse introduction")
	}

	if args["--socket"].(string) != "/asd" {
		t.Fatal("Didn't parse socket option")
	}
}
