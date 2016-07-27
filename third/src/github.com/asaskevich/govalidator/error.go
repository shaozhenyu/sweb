package govalidator

import (
	"fmt"
	"strings"
)

type Errors []error

func (es Errors) Error() string {
	errs := make([]string, len(es))
	for i, e := range es {
		errs[i] = e.Error()
	}
	return "[" + strings.Join(errs, ",") + "]"
}

type Error struct {
	Name string
	Err  error
}

func (e Error) Error() string {
	return fmt.Sprintf(`{"field":"%s","message":"%s"}`, e.Name, e.Err.Error())
}
