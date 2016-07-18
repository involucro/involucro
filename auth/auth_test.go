package auth

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"testing"
)

var source = `{
		"auths": [
			"https://user:pw@blubb.de",
			"https://alice:a11c3@blah.de?email=test@example.com",
			"https://a:b@index.docker.io/v1/"
		]
	}`

func unsetEnvVariable(t *testing.T) {
	if err := os.Unsetenv(ENV_NAME); err != nil {
		t.Fatal(err)
	}
}

func TestParseUserInfo(t *testing.T) {
	unsetEnvVariable(t)

	u, _ := url.Parse("https://test+asd:pw@quay.io")
	if u.User.Username() != "test+asd" {
		t.Error("user is wrong", u.User.Username())
	}
	if pw, set := u.User.Password(); !set || pw != "pw" {
		t.Error("password is wrong", pw)
	}
}

func TestGetAllFrom(t *testing.T) {
	unsetEnvVariable(t)

	expected := []string{
		"[user/pw: @ blubb.de]",
		"[alice/a11c3:test@example.com @ blah.de]",
		"[a/b: @ index.docker.io/v1/]",
	}

	urls, err := getAllFrom(strings.NewReader(source))
	if err != nil {
		t.Fatal(err)
	}

	if len(expected) != len(urls) {
		t.Fatal("length not equal")
	}

	for index := range expected {
		actual := fmt.Sprintf("[%s/%s:%s @ %s]",
			urls[index].Username, urls[index].Password, urls[index].Email, urls[index].ServerAddress)

		if expected[index] != actual {
			t.Errorf("Unexpected result for index %v: %v (expected %v)", index, urls[index], expected[index])
		}
	}
}

type authcase struct {
	server   string
	expected string
	found    bool
}

func runCases(cases []authcase, t *testing.T) {
	for i, el := range cases {
		auth, foundAuthentication, err := forServerInFile(el.server, strings.NewReader(source))
		if err != nil {
			t.Error("in case", i, err)
			continue
		}

		actual := fmt.Sprintf("[%s/%s:%s @ %s]",
			auth.Username, auth.Password, auth.Email, auth.ServerAddress)

		if foundAuthentication != el.found {
			t.Errorf("expected found authentication to be %t, but was %t in case %v", el.found, foundAuthentication, i)
		}

		if !foundAuthentication {
			continue
		}

		if el.expected != actual {
			t.Errorf("Unexpected result for index %v: %v (expected %v)", i, actual, el.expected)
		}
	}
}

func TestForServerInFile(t *testing.T) {
	unsetEnvVariable(t)

	cases := []authcase{
		{"blah.de", "[alice/a11c3:test@example.com @ blah.de]", true},
		{"blee.de", "[/: @ ]", false},
		{"", "[a/b: @ https://index.docker.io/v1/]", true},
		{"index.docker.io/v1/", "[a/b: @ https://index.docker.io/v1/]", true},
	}
	runCases(cases, t)
}

func TestWithFailingURLs(t *testing.T) {
	unsetEnvVariable(t)

	cases := []string{
		"\"http://withoutuser.com\"",
		"5",
		"\":example.com:90/\"",
	}

	for _, el := range cases {
		file := `{"auths": [` + el + `]}`
		_, _, err := forServerInFile("", strings.NewReader(file))
		if err == nil {
			t.Error("expected error in case ", el)
		}
	}
}

func TestUseEnvironmentVariable(t *testing.T) {
	unsetEnvVariable(t)

	if err := os.Setenv(ENV_NAME, "https://a:b_override@index.docker.io/v1/ https://x:y@example_env.com https://user:p2@blubb.de"); err != nil {
		t.Fatal(err)
	}

	cases := []authcase{
		// in none
		{"test.com", "", false},
		// only in config file
		{"blah.de", "[alice/a11c3:test@example.com @ blah.de]", true},
		// overriden in env variable
		{"index.docker.io/v1/", "[a/b_override: @ https://index.docker.io/v1/]", true},
		{"blubb.de", "[user/p2: @ blubb.de]", true},
		// only in env
		{"example_env.com", "[x/y: @ example_env.com]", true},
	}

	runCases(cases, t)
}
