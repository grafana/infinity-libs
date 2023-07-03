package transformations

import (
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/yesoreyeram/grafana-plugins/lib/go/framesql"
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
	if strings.TrimSpace(filterExpression) == "" {
		return frame, nil
	}
	filteredFrame := frame.EmptyCopy()
	rowLen, err := frame.RowLen()
	if err != nil {
		return frame, err
	}
	parsedExpression, err := govaluate.NewEvaluableExpressionWithFunctions(filterExpression, ExpressionFunctions)
	if err != nil {
		return frame, err
	}
	fieldKeys := map[string]bool{}
	for _, field := range frame.Fields {
		if fieldKeys[field.Name] {
			return frame, errors.New("field names are not unique. Not applying filter")
		}
		fieldKeys[field.Name] = true
	}
	for inRowIdx := 0; inRowIdx < rowLen; inRowIdx++ {
		var match *bool
		var err error
		parameters := map[string]any{"frame": frame, "null": nil, "nil": nil, "rowIndex": inRowIdx}
		for _, field := range frame.Fields {
			v := framesql.GetValue(field.At(inRowIdx))
			parameters[framesql.SlugifyFieldName(field.Name)] = v
			parameters[field.Name] = v
		}
		result, err := parsedExpression.Evaluate(parameters)
		if err != nil {
			return frame, fmt.Errorf("error evaluating filter expression for row %d. error: %w. Not applying filter", inRowIdx, err)
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
