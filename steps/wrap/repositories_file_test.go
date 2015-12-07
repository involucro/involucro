package wrap

import (
	"fmt"
)

func ExampleRepositoriesFile() {
	_, buf := repositoriesFile("test/gcc:latest", "283028")
	fmt.Printf("%s\n", buf)
	// Output: {"test/gcc":{"latest":"283028"}}
}

func ExampleRepoNameAndTagFrom() {
	repo, tag := repoNameAndTagFrom("foo/bar")
	fmt.Printf("%s %s", repo, tag)
	// Output: foo/bar latest
}

func ExampleRepoNameAndTagFrom_SpecifiedTag() {
	repo, tag := repoNameAndTagFrom("foo/bar:v1")
	fmt.Printf("%s %s", repo, tag)
	// Output: foo/bar v1
}
