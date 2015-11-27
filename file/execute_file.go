package file

import (
	log "github.com/Sirupsen/logrus"
)

// RunFile runs the file with the given filename in this context
func (i *InvContext) RunFile(fileName string) error {
	log.WithFields(log.Fields{"fileName": fileName}).Debug("Run file")
	return i.duk.PevalFile(fileName)
}

// RunString runs the given parameter directly
func (i *InvContext) RunString(script string) error {
	log.WithFields(log.Fields{"script": script}).Debug("Run script")
	return i.duk.PevalString(script)
}
