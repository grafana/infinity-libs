package utils

import (
	"errors"
	"strings"

	"github.com/tidwall/gjson"
)

func ValidateJson(jsonString string) (err error) {
	if strings.TrimSpace(jsonString) == "" {
		return errors.New("empty json received")
	}
	if !gjson.Valid(jsonString) {
		return errors.New("invalid json response received")
	}
	return err
}
