package runtime

import "testing"

func TestServerOfRepo(t *testing.T) {
	cases := []struct {
		repo   string
		server string
	}{
		{"", ""},
		{"alpine", ""},
		{"thriqon/hello", ""},
		{"gcr.io/thriqon/hello", "gcr.io"},
		{"test.local:5000/internal/too", "test.local:5000"},
	}

	for _, el := range cases {
		if s := serverOfRepo(el.repo); s != el.server {
			t.Errorf("%s is not recognized as server of %s, answer was %s", el.server, el.repo, s)
		}
	}
}
