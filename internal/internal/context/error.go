package context

import "errors"

var (
	errNilContext        = errors.New("error nil context")
	errMissingAppContext = errors.New("error missing app context")
)
