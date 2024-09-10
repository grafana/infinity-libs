package csvframer

import "errors"

var (
	ErrEmptyCsv = errors.New("empty/invalid csv")
	ErrReadingCsvResponse = errors.New("error reading csv response")
)