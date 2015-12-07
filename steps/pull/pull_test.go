package pull

import (
	"errors"
	"github.com/fsouza/go-dockerclient"
	"testing"
)

type mockPullable struct {
	lastPulled docker.PullImageOptions
	err        error
}

func (m *mockPullable) PullImage(opts docker.PullImageOptions, _ docker.AuthConfiguration) error {
	m.lastPulled = opts
	m.lastPulled.OutputStream.Write([]byte("{}"))
	m.lastPulled.OutputStream.Write([]byte("{\"status\": 5}"))
	return m.err
}

func TestPull(t *testing.T) {
	var m mockPullable

	err := Pull(&m, "test/asd")
	if err != nil {
		t.Fatal("Err was not nil")
	}
	if m.lastPulled.Repository != "test/asd" {
		t.Fatalf("Pulled %s instead of test/asd", m.lastPulled.Repository)
	}

	m.err = errors.New("Mocked error")
	err = Pull(&m, "test/asd2")

	if err == nil || err.Error() != m.err.Error() {
		t.Fatalf("err != m.err: actual %s, expected %s", err, m.err)
	}
}
