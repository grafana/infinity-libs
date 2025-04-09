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
		return nil, ErrEvaluatingFilterExpressionWithEmptyFrame
	}
	if strings.TrimSpace(filterExpression) == "" {
		return frame, nil
	}
	filteredFrame := frame.EmptyCopy()
	for i := range filteredFrame.Fields {
		if frame.Fields[i].Labels == nil {
			filteredFrame.Fields[i].Labels = nil
		}
	}
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
			v := framesql.GetValue(field.At(inRowIdx))
			if field.Nullable() && field.NilAt(inRowIdx) {
				v = nil
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
			return frame, fmt.Errorf("filter expression for row %d didn't produce binary result. error: %w. Not applying filter", inRowIdx, err)
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
