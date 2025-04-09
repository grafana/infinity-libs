package transformations_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/transformations"
	"github.com/grafana/infinity-libs/lib/go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFilter(t *testing.T) {
	devicesFrame := data.NewFrame(
		"devices",
		data.NewField("id", nil, []int64{0, 1, 2, 3, 4, 5}),
		data.NewField("name", nil, []string{"iPhone 6S", "iPhone 5S", "MacBook", "MacBook Air", "MacBook Air 2013", "MacBook Air 2012"}),
		data.NewField("price", nil, []int64{799, 349, 1499, 999, 599, 499}),
	)
	tests := []struct {
		name             string
		frame            *data.Frame
		filterExpression string
		want             *data.Frame
		wantErr          error
	}{
		{
			name:    "empty frame",
			want:    data.NewFrame("test"),
			wantErr: transformations.ErrEvaluatingFilterExpressionWithEmptyFrame,
		},
		{
			name:             "numeric filter expression",
			frame:            data.NewFrame("test", data.NewField("num", nil, []*int64{utils.P(int64(1)), utils.P(int64(2)), utils.P(int64(3)), utils.P(int64(4)), utils.P(int64(5))})),
			filterExpression: "num > 2 && num < 5",
			want:             data.NewFrame("test", data.NewField("num", nil, []*int64{utils.P(int64(3)), utils.P(int64(4))})),
		},
		{
			name:             "string filter expression",
			frame:            data.NewFrame("test", data.NewField("user", nil, []string{"foo", "bar", "baz"})),
			filterExpression: "user == 'foo' || user == 'baz'",
			want:             data.NewFrame("test", data.NewField("user", nil, []string{"foo", "baz"})),
		},
		{
			name:             "boolean filter expression",
			frame:            data.NewFrame("test", data.NewField("user", nil, []string{"foo", "bar", "baz"}), data.NewField("active", nil, []bool{true, false, true})),
			filterExpression: "active == true",
			want:             data.NewFrame("test", data.NewField("user", nil, []string{"foo", "baz"}), data.NewField("active", nil, []bool{true, true})),
		},
		{
			name:             "null value filter expression",
			frame:            data.NewFrame("test", data.NewField("user", nil, []string{"foo", "bar", "baz"}), data.NewField("salary", nil, []*int64{utils.P(int64(300)), nil, utils.P(int64(400))})),
			filterExpression: "salary == null",
			want:             data.NewFrame("test", data.NewField("user", nil, []string{"bar"}), data.NewField("salary", nil, []*int64{nil})),
		},
		{
			name:             "nil value filter expression",
			frame:            data.NewFrame("test", data.NewField("user", nil, []string{"foo", "bar", "baz"}), data.NewField("salary", nil, []*int64{utils.P(int64(300)), nil, utils.P(int64(400))})),
			filterExpression: "salary == nil",
			want:             data.NewFrame("test", data.NewField("user", nil, []string{"bar"}), data.NewField("salary", nil, []*int64{nil})),
		},
		{
			name:             "decices with numeric filter expression",
			frame:            devicesFrame,
			filterExpression: "price > 500",
			want: data.NewFrame(
				"devices",
				data.NewField("id", nil, []int64{0, 2, 3, 4}),
				data.NewField("name", nil, []string{"iPhone 6S", "MacBook", "MacBook Air", "MacBook Air 2013"}),
				data.NewField("price", nil, []int64{799, 1499, 999, 599}),
			),
		},
		{
			name:             "decices with multi filter expression",
			frame:            devicesFrame,
			filterExpression: "name != 'MacBook' && price > 400",
			want: data.NewFrame(
				"devices",
				data.NewField("id", nil, []int64{0, 3, 4, 5}),
				data.NewField("name", nil, []string{"iPhone 6S", "MacBook Air", "MacBook Air 2013", "MacBook Air 2012"}),
				data.NewField("price", nil, []int64{799, 999, 599, 499}),
			),
		},
		{
			name:             "decices with IN filter expression",
			frame:            devicesFrame,
			filterExpression: "!(name IN ('MacBook','MacBook Air'))",
			want: data.NewFrame(
				"devices",
				data.NewField("id", nil, []int64{0, 1, 4, 5}),
				data.NewField("name", nil, []string{"iPhone 6S", "iPhone 5S", "MacBook Air 2013", "MacBook Air 2012"}),
				data.NewField("price", nil, []int64{799, 349, 599, 499}),
			),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			frame := tt.frame
			got, err := transformations.ApplyFilter(frame, tt.filterExpression)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.Nil(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
