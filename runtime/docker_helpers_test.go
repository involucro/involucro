package runtime

import "fmt"

func ExampleRepoNameAndTagFrom() {
	repo, tag, autotag := repoNameAndTagFrom("foo/bar")
	fmt.Printf("[%s] [%s] [%s]", repo, tag, autotag)
	// Output: [foo/bar] [] [latest]
}

func ExampleRepoNameAndTagFrom_SpecifiedTag() {
	repo, tag, autotag := repoNameAndTagFrom("foo/bar:v1")
	fmt.Printf("[%s] [%s] [%s]", repo, tag, autotag)
	// Output: [foo/bar] [v1] [v1]
}

func ExampleRepoNameAndTagFrom_WithPrivateRepository() {
	repo, tag, autotag := repoNameAndTagFrom("192.168.0.1:5000/foo/bar:v1")
	fmt.Printf("[%s] [%s] [%s]", repo, tag, autotag)
	// Output: [192.168.0.1:5000/foo/bar] [v1] [v1]
}
