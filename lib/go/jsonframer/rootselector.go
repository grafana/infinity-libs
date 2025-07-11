package jsonframer

import (
	"encoding/json"
	"errors"
	"strings"
	"github.com/itchyny/gojq"
	"github.com/tidwall/gjson"
	jsonata "github.com/xiatechs/jsonata-go"
)

func ApplyRootSelector(jsonString string, rootSelector string, framerType FramerType) (string, error) {
	if rootSelector == "" {
		return jsonString, nil
	}
	if framerType == FramerTypeGJSON {
		return ApplyRootSelectorUsingGJSON(jsonString, rootSelector)
	}
	if framerType == FramerTypeJsonata {
		return ApplyRootSelectorUsingJSONata(jsonString, rootSelector)
	}
	if framerType == FramerTypeJQ {
		return ApplyRootSelectorUsingJQ(jsonString, rootSelector)
	}
	return ApplyRootSelectorUsingWithGuess(jsonString, rootSelector)
}

func ApplyRootSelectorUsingGJSON(jsonString string, rootSelector string) (string, error) {
	r := gjson.Get(jsonString, rootSelector)
	if r.Exists() {
		return r.String(), nil
	}
	return jsonString, ErrInvalidRootSelector
}

func ApplyRootSelectorUsingJSONata(jsonString string, rootSelector string) (string, error) {
	expr, err := jsonata.Compile(rootSelector)
	if err != nil {
		return "", errors.Join(ErrInvalidRootSelector, err)
	}
	if expr == nil {
		return "", errors.Join(ErrInvalidRootSelector)
	}
	return ApplyRootSelectorUsingJSONataExpression(jsonString, expr)
}

func ApplyRootSelectorUsingJSONataExpression(jsonString string, expr *jsonata.Expr) (string, error) {
	var data any
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	res, err := expr.Eval(data)
	if err != nil {
		if strings.Contains(err.Error(), "argument") { // TODO: Fix error handling in jsonata-go library when dealing with missed argument
			return "", ErrEvaluatingJSONata
		}
		return "", errors.Join(ErrEvaluatingJSONata, err)
	}
	r2, err := json.Marshal(res)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	return string(r2), nil
}

func ApplyRootSelectorUsingJQ(jsonString string, rootSelector string) (string, error) {
	query, err := gojq.Parse(rootSelector)
	if err != nil {
		return "", errors.Join(ErrInvalidJQSelector, err)
	}
	return ApplyRootSelectorUsingJQQuery(jsonString, query)
}

func ApplyRootSelectorUsingJQQuery(jsonString string, query *gojq.Query) (string, error) {
	var data any
	err := json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", errors.Join(ErrUnMarshalingJSON, err)
	}
	iter := query.Run(data)
	out := []any{}
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			if err, ok := err.(*gojq.HaltError); ok && err.Value() == nil {
				break
			}
			return "", errors.Join(ErrExecutingJQ, err)
		}
		out = append(out, v)
	}
	if len(out) == 1 {
		if v, ok := out[0].([]any); ok {
			outStr, err := json.Marshal(v)
			if err != nil {
				return "", errors.Join(ErrMarshalingJSON, err)
			}
			return string(outStr), nil
		}
	}
	outStr, err := json.Marshal(out)
	if err != nil {
		return "", errors.Join(ErrMarshalingJSON, err)
	}
	return string(outStr), nil
}

// ApplyRootSelectorUsingWithGuess method try to guess the root selector (for legacy reasons)
// First it will attempt to GJSON based selector (used by many users)
// If the GJSON based selection failed, it will default fallback to JSONata based selection
// Intentionally we are not falling back to JQ yet
func ApplyRootSelectorUsingWithGuess(jsonString string, rootSelector string) (string, error) {
	r := gjson.Get(jsonString, rootSelector)
	if r.Exists() {
		return r.String(), nil
	}
	return ApplyRootSelectorUsingJSONata(jsonString, rootSelector)
}
