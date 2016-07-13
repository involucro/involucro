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
			t.Errorf("Unexpected result for index %v: %v", index, urls[index])
		}
	}
}

func TestForServerInFile(t *testing.T) {
	unsetEnvVariable(t)

	cases := []struct {
		server   string
		expected string
		found    bool
	}{
		{"blah.de", "[alice/a11c3:test@example.com @ blah.de]", true},
		{"blee.de", "[/: @ ]", false},
		{"", "[a/b: @ https://index.docker.io/v1/]", true},
		{"index.docker.io/v1/", "[a/b: @ https://index.docker.io/v1/]", true},
	}

	for i, el := range cases {
		auth, foundAuthentication, err := forServerInFile(el.server, strings.NewReader(source))
		if err != nil {
			t.Error("in case", i, err)
			continue
		}

		actual := fmt.Sprintf("[%s/%s:%s @ %s]",
			auth.Username, auth.Password, auth.Email, auth.ServerAddress)

		if el.expected != actual {
			t.Errorf("unexpected result for index %v: %v", i, auth)
		}

		if foundAuthentication != el.found {
			t.Errorf("expected found authentication to be %t, but was %t in case %v", el.found, foundAuthentication, i)
		}
	}
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
