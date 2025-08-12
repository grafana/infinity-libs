package framesql

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	"github.com/casbin/govaluate"
	"github.com/grafana/grafana-plugin-sdk-go/data"
	dataFrameField "github.com/grafana/infinity-libs/lib/go/framesql/field"
)

func EvaluateInFrame(expression string, input *data.Frame) (any, error) {
	if strings.TrimSpace(expression) == "" {
		return nil, ErrEmptySummarizeExpression
	}
	parsedExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expression, ExpressionFunctions)
	if err != nil {
		return nil, err
	}
	frameLen, _ := input.RowLen()
	parameters := map[string]any{"frame": input, "recordsCount": frameLen}
	if input != nil {
		for _, field := range input.Fields {
			parameters[SlugifyFieldName(field.Name)] = field
			parameters[field.Name] = field
		}
	}
	result, err := parsedExpression.Evaluate(parameters)
	if err != nil {
		if checkIfExpressionFoundInFields(err) {
			return result, errors.Join(err, ErrExpressionNotFoundInFields)
		}
	}
	return result, err
}

var ExpressionFunctions = map[string]govaluate.ExpressionFunction{
	"count": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("count", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Count(field)
	},
	"first": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("first", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.First(field)
	},
	"last": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("last", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Last(field)
	},
	"sum": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("sum", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Sum(field)
	},
	"min": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("min", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Min(field)
	},
	"max": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("max", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Max(field)
	},
	"mean": func(arguments ...any) (any, error) {
		field, err := getFieldFromArguments("mean", arguments...)
		if err != nil {
			return nil, err
		}
		return dataFrameField.Mean(field)
	},
}

func SlugifyFieldName(input string) string {
	re, _ := regexp.Compile(`[^\w]`)
	input = strings.TrimSpace(re.ReplaceAllString(strings.ToLower(strings.TrimSpace(input)), "_"))
	return input
}

func checkIfExpressionFoundInFields(err error) bool {
	r := regexp.MustCompile(`No parameter '(.+)' found\.`)
	// Check if the message matches the pattern
	matches := r.FindStringSubmatch(err.Error())
	return len(matches) > 0
}

func getFieldFromArguments(functionName string, arguments ...any) (*data.Field, error) {
	if len(arguments) < 1 {
		return nil, fmt.Errorf("invalid arguments to %s method", functionName)
	}
	field, ok := arguments[0].(*data.Field)
	if !ok {
		return nil, fmt.Errorf("first argument is not a valid field to %s method", functionName)
	}
	return field, nil
}
