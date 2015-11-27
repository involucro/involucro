package file

import (
	log "github.com/Sirupsen/logrus"
)

func (i *InvContext) RunFile(fileName string) error {
	log.WithFields(log.Fields{"fileName": fileName}).Debug("Run file")
	return i.duk.PevalFile(fileName)
}
