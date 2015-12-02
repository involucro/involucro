package run

import "io"
import "fmt"

// AsShellCommandOn prints sh compatible commands into the given writer, that
// accomplish the funciontality encoded in this step
func (img ExecuteImage) AsShellCommandOn(w io.Writer) {
	fmt.Fprintf(w, "docker run -t --rm %s\n", img.Config.Image)
}
