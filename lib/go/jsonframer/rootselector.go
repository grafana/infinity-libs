package jsonframer

import (
	"encoding/json"
	"errors"

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
	r := gjson.Get(string(jsonString), rootSelector)
	if r.Exists() {
		return r.String(), nil
	}
	return jsonString, ErrInvalidRootSelector
}

func ApplyRootSelectorUsingJQ(jsonString string, rootSelector string) (string, error) {
	query, err := gojq.Parse(rootSelector)
	if err != nil {
		return "", errors.Join(ErrInvalidJQSelector, err)
	}
	var data any
	err = json.Unmarshal([]byte(jsonString), &data)
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
	// if the result is array of 1 element array, ignore the outer array and return that single element from inner array
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

func ApplyRootSelectorUsingJSONata(jsonString string, rootSelector string) (string, error) {
	expr, err := jsonata.Compile(rootSelector)
	if err != nil {
		return "", errors.Join(ErrInvalidRootSelector, err)
	}
	if expr == nil {
		return "", errors.Join(ErrInvalidRootSelector)
	}
	var data any
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	res, err := expr.Eval(data)
	if err != nil {
		return "", errors.Join(ErrEvaluatingJSONata, err)
	}
	r2, err := json.Marshal(res)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	return string(r2), nil
}

func ApplyRootSelectorUsingWithGuess(jsonString string, rootSelector string) (string, error) {
	r := gjson.Get(string(jsonString), rootSelector)
	if r.Exists() {
		return r.String(), nil
	}
	expr, err := jsonata.Compile(rootSelector)
	if err != nil {
		return "", errors.Join(ErrInvalidRootSelector, err)
	}
	if expr == nil {
		return "", errors.Join(ErrInvalidRootSelector)
	}
	var data any
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	res, err := expr.Eval(data)
	if err != nil {
		return "", errors.Join(ErrEvaluatingJSONata, err)
	}
	r2, err := json.Marshal(res)
	if err != nil {
		return "", errors.Join(ErrInvalidJSONContent, err)
	}
	return string(r2), nil
}
