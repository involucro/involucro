package file

import (
	log "github.com/Sirupsen/logrus"
	duk "gopkg.in/olebedev/go-duktape.v2"
)

// InstantiateRuntimeEnv creates a new InvContext and returns it. This new
// context uses the working dir that is passed as a parameter.  After
// instantiation, the context will be ready to load additional files.
func InstantiateRuntimeEnv(workingDir string) InvContext {
	m := InvContext{
		duk:        duk.New(),
		Tasks:      make(map[string][]Step),
		WorkingDir: workingDir,
	}

	global := m.duk.PushObject()

	idx := m.duk.PushObject()

	defineFuncOnObject(m.duk, idx, "task", m.asCallback(func(i *InvContext) int {
		taskID := requireStringOrFailGracefully(i.duk, -1, "task")
		log.WithFields(log.Fields{"taskId": taskID}).Info("defined task")

		retobj := pushDerivedFromThis(i.duk)

		defineStringOnObject(i.duk, retobj, "taskId", taskID)
		defineFuncOnObject(i.duk, retobj, "using", m.asCallback(fileUsing))
		defineFuncOnObject(i.duk, retobj, "wrap", m.asCallback(fileWrapping))

		return 1
	}))

	m.duk.PutPropString(global, "inv")

	//idxLogger := m.duk.PushObject()

	m.duk.SetGlobalObject()

	return m
}
