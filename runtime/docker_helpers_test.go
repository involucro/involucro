package runtime

import "testing"

func TestRepoNameAndTagFrom(t *testing.T) {
	cases := []struct {
		source string
		result []string
	}{
		{"foo/bar", []string{"foo/bar", "", "latest"}},
		{"foo/bar:v1", []string{"foo/bar", "v1", "v1"}},
		{"192.168.0.1:5000/foo/bar:v1", []string{"192.168.0.1:5000/foo/bar", "v1", "v1"}},
	}

	for _, el := range cases {
		repo, tag, autotag := repoNameAndTagFrom(el.source)
		if repo != el.result[0] {
			t.Errorf("repo has unexpected value in case %s, it is %s", el.source, repo)
		}
		if tag != el.result[1] {
			t.Errorf("tag has unexpected value in case %s, it is %s", el.source, tag)
		}
		if autotag != el.result[2] {
			t.Errorf("autotag has unexpected value in case %s, it is %s", el.source, autotag)
		}
	}
}
