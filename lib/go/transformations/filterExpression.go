package transformations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/framesql"
	"gopkg.in/Knetic/govaluate.v3"
)

type FilterExpressionOptions struct {
	Expression string `json:"expression,omitempty"`
}

func FilterExpression(input []*data.Frame, options FilterExpressionOptions) ([]*data.Frame, error) {
	output := []*data.Frame{}
	for _, frame := range input {
		filteredFrame, err := ApplyFilter(frame, options.Expression)
		if err != nil {
			return output, errors.New("unable to apply filter")
		}
		output = append(output, filteredFrame)
	}
	return output, nil
}

func ApplyFilter(frame *data.Frame, filterExpression string) (*data.Frame, error) {
	if frame == nil {
		return frame, nil
	}
	if frame.Rows() == 0 {
		return frame, nil
	}
	if strings.TrimSpace(filterExpression) == "" {
		return frame, nil
	}
	filteredFrame := frame.EmptyCopy()
	filteredFrame.Meta = frame.Meta
	rowLen, err := frame.RowLen()
	if err != nil {
		return frame, err
	}
	parsedExpression, err := govaluate.NewEvaluableExpressionWithFunctions(filterExpression, ExpressionFunctions)
	if err != nil {
		if checkIfInvalidFilterExpression(err) {
			return frame, fmt.Errorf("%w. %w", ErrInvalidFilterExpression, err)
		}
		return frame, err
	}
	fieldKeys := map[string]bool{}
	for _, field := range frame.Fields {
		if fieldKeys[field.Name] {
			return frame, ErrNotUniqueFieldNames
		}
		fieldKeys[field.Name] = true
	}
	for inRowIdx := 0; inRowIdx < rowLen; inRowIdx++ {
		var match *bool
		var err error
		parameters := map[string]any{"frame": frame, "null": nil, "nil": nil, "rowIndex": inRowIdx}
		for _, field := range frame.Fields {
			var v any = nil
			if fieldOriginalValue, ok := field.ConcreteAt(inRowIdx); fieldOriginalValue != nil && ok {
				v = framesql.GetValue(fieldOriginalValue)
			}
			parameters[framesql.SlugifyFieldName(field.Name)] = v
			parameters[field.Name] = v
		}
		result, err := parsedExpression.Evaluate(parameters)
		if err != nil {
			return frame, errors.Join(ErrEvaluatingFilterExpression, fmt.Errorf("error: %w. row %d. Not applying filter", err, inRowIdx))
		}
		if currentMatch, ok := result.(bool); ok {
			match = &currentMatch
		}
		if currentMatch, ok := result.(*bool); ok {
			match = currentMatch
		}
		if match == nil {
			return frame, fmt.Errorf("filter expression for row %d didn't produce binary result. Not applying filter", inRowIdx)
		}
		if !*match {
			continue
		}
		filteredFrame.AppendRow(frame.RowCopy(inRowIdx)...)
	}
	return filteredFrame, nil
}

func checkIfInvalidFilterExpression(err error) bool {
	// Check if the error is due to an invalid token
	// This error is not exported by the govaluate library, so we need to check the error message
	return strings.Contains(err.Error(), "Invalid token:")
}
