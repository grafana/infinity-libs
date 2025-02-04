package jsonframer

import "errors"

var (
	ErrInvalidRootSelector = errors.New("failed to compile JSONata expression")
	ErrEvaluatingJSONata   = errors.New("error evaluating JSONata expression")
	ErrInvalidJSONContent  = errors.New("invalid/empty JSON")
	ErrInvalidJQSelector   = errors.New("failed to compile jq selector")
	ErrUnMarshalingJSON    = errors.New("error while un-marshaling json")
	ErrMarshalingJSON      = errors.New("error while marshaling json")
	ErrExecutingJQ         = errors.New("error while executing JQ")
)
