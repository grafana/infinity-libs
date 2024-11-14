package framesql_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/framesql"
	"github.com/grafana/infinity-libs/lib/go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExpressionFunctions(t *testing.T) {
	t.Run("count", func(t *testing.T) {
		t.Run("all values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["count"](data.NewField("values", nil, []bool{true, false, true}))
			require.Equal(t, float64(3), out)
			require.Nil(t, err)
		})
		t.Run("no values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["count"](data.NewField("values", nil, []int64{}))
			require.Equal(t, float64(0), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["count"](data.NewField("values", nil, []*bool{utils.P(true), nil, utils.P(false)}))
			require.Equal(t, float64(3), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer at first position", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["count"](data.NewField("values", nil, []*bool{nil, utils.P(true), utils.P(false)}))
			require.Equal(t, float64(3), out)
			require.Nil(t, err)
		})
		t.Run("all int64", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["count"](data.NewField("values", nil, []int64{1, 2, 3}))
			require.Equal(t, float64(3), out)
			require.Nil(t, err)
		})
	})
	t.Run("first", func(t *testing.T) {
		t.Run("all values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []bool{true, false, true}))
			require.Equal(t, true, out)
			require.Nil(t, err)
		})
		t.Run("no values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []int64{}))
			require.Equal(t, nil, out)
			require.Nil(t, err)
		})
		t.Run("single value", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []int64{12}))
			require.Equal(t, int64(12), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []*bool{utils.P(true), nil, utils.P(false)}))
			require.Equal(t, utils.P(true), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer at first position", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []*bool{nil, utils.P(true), utils.P(false)}))
			require.Equal(t, nil, out)
			require.Nil(t, err)
		})
		t.Run("all int64", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["first"](data.NewField("values", nil, []int64{1, 2, 3}))
			require.Equal(t, int64(1), out)
			require.Nil(t, err)
		})
	})
	t.Run("last", func(t *testing.T) {
		t.Run("all values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []bool{true, false, true}))
			require.Equal(t, true, out)
			require.Nil(t, err)
		})
		t.Run("no values", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []int64{}))
			require.Equal(t, nil, out)
			require.Nil(t, err)
		})
		t.Run("single value", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []int64{12}))
			require.Equal(t, int64(12), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []*bool{utils.P(true), nil, utils.P(false)}))
			require.Equal(t, utils.P(false), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer string", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []*string{utils.P("hello"), nil, utils.P("world")}))
			require.Equal(t, utils.P("world"), out)
			require.Nil(t, err)
		})
		t.Run("with nil pointer at last position", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []*bool{utils.P(false), utils.P(true), nil}))
			require.Equal(t, nil, out)
			require.Nil(t, err)
		})
		t.Run("all int64", func(t *testing.T) {
			out, err := framesql.ExpressionFunctions["last"](data.NewField("values", nil, []int64{1, 2, 3}))
			require.Equal(t, int64(3), out)
			require.Nil(t, err)
		})
	})
}

func TestEvaluateInFrame(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		input      *data.Frame
		want       any
		wantErr    error
	}{
		{
			wantErr: framesql.ErrEmptySummarizeExpression,
		},
		{
			input:      data.NewFrame("test", data.NewField("sample", nil, []*float64{toFP(1), toFP(2), toFP(0.5), toFP(1.5)})),
			expression: `count(sample)`,
			want:       float64(4),
		},
		{
			input:      data.NewFrame("test", data.NewField("sample", nil, []*float64{utils.P(float64(1)), utils.P(float64(2)), utils.P(float64(0.5)), utils.P(float64(1.5))})),
			expression: `sum(sample)`,
			want:       float64(5),
		},
		{
			input:      data.NewFrame("test", data.NewField("sample", nil, []*float64{utils.P(float64(1)), utils.P(float64(2)), utils.P(float64(0.5)), utils.P(float64(1.5))})),
			expression: `min(sample)`,
			want:       float64(0.5),
		},
		{
			input:      data.NewFrame("test", data.NewField("sample", nil, []*float64{utils.P(float64(1)), utils.P(float64(2)), utils.P(float64(0.5)), utils.P(float64(1.5))})),
			expression: `max(sample)`,
			want:       float64(2),
		},
		{
			input:      data.NewFrame("test", data.NewField("hello 1st world!", nil, []*float64{utils.P(float64(1)), utils.P(float64(2)), utils.P(float64(0.5)), utils.P(float64(1.5))})),
			expression: `mean(hello_1st_world_)`,
			want:       float64(1.25),
		},
		{
			input:      data.NewFrame("test", data.NewField("grafana infinity 數據源!", nil, []*float64{utils.P(float64(1)), utils.P(float64(2)), utils.P(float64(0.5)), utils.P(float64(1.5))})),
			expression: `mean([grafana infinity 數據源!])`,
			want:       float64(1.25),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := framesql.EvaluateInFrame(tt.expression, tt.input)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr, err)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}

func toFP(v float64) *float64 {
	return &v
}
