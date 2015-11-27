package file

import (
	wrap "github.com/thriqon/involucro/steps/wrap"
)

func fileWrapping(c *InvContext) int {
	sourceDir := requireStringOrFailGracefully(c.duk, -1, "wrap")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "sourceDir", sourceDir)
	defineFuncOnObject(c.duk, idx, "inImage", c.asCallback(fileInImage))
	defineFuncOnObject(c.duk, idx, "at", c.asCallback(fileAt))
	defineFuncOnObject(c.duk, idx, "as", c.asCallback(fileAs))

	return 1
}

func fileInImage(c *InvContext) int {
	inImage := requireStringOrFailGracefully(c.duk, -1, "inImage")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "parentImage", inImage)
	return 1
}

func fileAt(c *InvContext) int {
	targetDir := requireStringOrFailGracefully(c.duk, -1, "at")
	idx := pushDerivedFromThis(c.duk)

	defineStringOnObject(c.duk, idx, "targetDir", targetDir)
	return 1
}

func fileAs(i *InvContext) int {
	c := i.duk
	newName := requireStringOrFailGracefully(c, -1, "as")
	pushDerivedFromThis(c)

	c.GetPropString(-1, "taskId")
	taskID := requireStringOrFailGracefully(c, -1, "run:task_id")
	c.Pop()

	c.GetPropString(-1, "parentImage")
	parentImage := requireStringOrFailGracefully(c, -1, "as:inImage")
	c.Pop()

	c.GetPropString(-1, "targetDir")
	targetDir := requireStringOrFailGracefully(c, -1, "as:at")
	c.Pop()

	c.GetPropString(-1, "sourceDir")
	sourceDir := requireStringOrFailGracefully(c, -1, "as:wrap")
	c.Pop()

	wi := wrap.AsImage{
		SourceDir:         sourceDir,
		TargetDir:         targetDir,
		ParentImage:       parentImage,
		NewRepositoryName: newName,
	}

	i.Tasks[taskID] = append(i.Tasks[taskID], wi)

	return 1
}
