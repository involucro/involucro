package runtime

// Step represents one action taken by the tool.
type Step interface {
	Take(i *Runtime) error

	// ShowStartInfo displays some information on the default logger that identifies the step
	ShowStartInfo()
}
