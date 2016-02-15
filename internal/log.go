package runtime

import "github.com/thriqon/involucro/ilog"

var (
	logProgress = ilog.ForLevelPrefix(-3, "PRGS")
	logStdout   = ilog.ForLevelPrefix(-2, "SOUT")
	logStderr   = ilog.ForLevelPrefix(-2, "SERR")
	logTask     = ilog.ForLevelPrefix(-1, "STEP")
)
