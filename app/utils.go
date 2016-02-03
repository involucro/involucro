package app

import (
	"fmt"
	"strings"
)

type variablesValue map[string]string

func (v *variablesValue) String() string {
	var items []string
	for k, v := range *v {
		items = append(items, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("[%s]", strings.Join(items, " "))
}

type errInvalidFormatForVariableAssignment string

func (e errInvalidFormatForVariableAssignment) Error() string {
	return fmt.Sprintf("Invalid value [%s], expected value of the form: KEY=VALUE", string(e))
}

func (v *variablesValue) Set(s string) error {
	if *v == nil {
		*v = make(map[string]string)
	}
	ss := strings.SplitN(s, "=", 2)
	if len(ss) < 2 {
		return errInvalidFormatForVariableAssignment(s)
	}
	(*v)[ss[0]] = ss[1]
	return nil
}
