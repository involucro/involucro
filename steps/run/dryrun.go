package run

import log "github.com/Sirupsen/logrus"

// DryRun runs this task without doing anything, but logging an indication of
// what would have been done
func (img ExecuteImage) DryRun() {
	log.WithFields(log.Fields{"dry": true, "image": img.Config.Image}).Info("RUN")
}
