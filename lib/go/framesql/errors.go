package framesql

import "errors"

var (
	ErrEmptySummarizeExpression = errors.New("empty/invalid summarize expression")
	ErrExpressionNotFoundInFields = errors.New("expression not found in fields")
)
