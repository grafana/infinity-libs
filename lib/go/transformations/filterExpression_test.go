package transformations_test

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/transformations"
	"github.com/grafana/infinity-libs/lib/go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFilter(t *testing.T) {
	sampleDataFrame := data.NewFrame("hello", data.NewField("group", nil, []string{"A", "B", "A"}), data.NewField("id", nil, []int64{3, 4, 5}), data.NewField("value", nil, []int64{6, 7, 8})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"})
	A := "a"
	B := "b"
	zero := int64(0)
	one := int64(1)
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
			name:             "nil frame should throw error",
			filterExpression: "group =='A'",
			wantErr:          transformations.ErrEvaluatingFilterExpressionWithEmptyFrame,
		},
		{
			name:             "frame without fields should return the same",
			filterExpression: "group =='A'",
			frame:            data.NewFrame("hello").SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want:             data.NewFrame("hello").SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "frame with emtpy fields should return the same",
			filterExpression: "group =='A'",
			frame:            data.NewFrame("hello", data.NewField("field1", nil, []int64{})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want:             data.NewFrame("hello", data.NewField("field1", nil, []int64{})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:  "frame with data and without filter should return the same",
			frame: sampleDataFrame,
			want:  sampleDataFrame,
		},
		{
			name:             "frame with data and with filter should filter the data with matching condition",
			filterExpression: "group =='A'",
			frame:            sampleDataFrame,
			want: data.NewFrame("hello",
				data.NewField("group", nil, []string{"A", "A"}),
				data.NewField("id", nil, []int64{3, 5}),
				data.NewField("value", nil, []int64{6, 8}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "frame with data and with filter should filter the data without matching condition",
			filterExpression: "id == 1",
			frame:            sampleDataFrame,
			want: data.NewFrame("hello",
				data.NewField("group", nil, []string{}),
				data.NewField("id", nil, []int64{}),
				data.NewField("value", nil, []int64{}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "frame with data and with filter should filter the data with incorrect matching condition",
			filterExpression: "group == 3",
			frame:            sampleDataFrame,
			want: data.NewFrame("hello",
				data.NewField("group", nil, []string{}),
				data.NewField("id", nil, []int64{}),
				data.NewField("value", nil, []int64{}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "null value filter",
			filterExpression: "value != nil",
			frame:            data.NewFrame("hello", data.NewField("name", nil, []*string{&A, &B}), data.NewField("value", nil, []*string{&A, nil})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want: data.NewFrame("hello",
				data.NewField("name", nil, []*string{&A}),
				data.NewField("value", nil, []*string{&A}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "null value filter with number",
			filterExpression: "value != nil",
			frame:            data.NewFrame("hello", data.NewField("name", nil, []*string{&A, &B, &A}), data.NewField("value", nil, []*int64{&zero, &one, nil})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want:             data.NewFrame("hello", data.NewField("name", nil, []*string{&A, &B}), data.NewField("value", nil, []*int64{&zero, &one})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "invalid filter should throw error",
			filterExpression: "group ==='A'",
			frame:            sampleDataFrame,
			wantErr:          errors.New("invalid filter expression. Invalid token: '==='"),
		},
		{
			name:             "non binary filter should throw error",
			filterExpression: "1 + 2",
			frame:            sampleDataFrame,
			wantErr:          errors.New("filter expression for row 0 didn't produce binary result. Not applying filter"),
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
		{
			// Without Filter: 	2024-01-01, 2025-01-01, 2026-01-01, 2027-01-01
			// With Filter:		2025-01-01, 2026-01-01
			name:             "non nullable date field expression",
			frame:            data.NewFrame("test", data.NewField("time", nil, []time.Time{time.Unix(1704067200, 0).UTC(), time.Unix(1735689600, 0).UTC(), time.Unix(1767225600, 0).UTC(), time.Unix(1798761600, 0).UTC()})),
			filterExpression: "time > 1704067200 && time < 1798761600",
			want:             data.NewFrame("test", data.NewField("time", nil, []time.Time{time.Unix(1735689600, 0).UTC(), time.Unix(1767225600, 0).UTC()})),
		},
		{
			// Without Filter: 	2024-01-01, 2025-01-01, 2026-01-01, 2027-01-01, 2028-01-01
			// With Filter:		2025-01-01, 2026-01-01
			name:             "nullable date field expression",
			frame:            data.NewFrame("test", data.NewField("time", nil, []*time.Time{utils.P(time.Unix(1704067200, 0).UTC()), utils.P(time.Unix(1735689600, 0).UTC()), nil, utils.P(time.Unix(1798761600, 0).UTC()), utils.P(time.Unix(1830297600, 0).UTC())})),
			filterExpression: "time > 1704067200 && time < 1798761600",
			want:             data.NewFrame("test", data.NewField("time", nil, []*time.Time{utils.P(time.Unix(1735689600, 0).UTC())})),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformations.ApplyFilter(tt.frame, tt.filterExpression)
			if tt.wantErr != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				return
			}
			require.Nil(t, err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_GetNormalizedValueForExpressionEvaluation(t *testing.T) {
	tests := []struct {
		field *data.Field
		wantV any
	}{
		{field: data.NewField("int64", nil, []int64{1}), wantV: float64(1)},
		{field: data.NewField("int32", nil, []int32{1}), wantV: float64(1)},
		{field: data.NewField("int16", nil, []int16{1}), wantV: float64(1)},
		{field: data.NewField("int8", nil, []int8{1}), wantV: float64(1)},
		{field: data.NewField("float64", nil, []float64{1}), wantV: float64(1)},
		{field: data.NewField("float32", nil, []float32{1}), wantV: float64(1)},
		{field: data.NewField("bool", nil, []bool{true}), wantV: true},
		{field: data.NewField("string", nil, []string{"test"}), wantV: "test"},
		{field: data.NewField("time", nil, []time.Time{time.Unix(1704067200, 0).UTC()}), wantV: time.Unix(1704067200, 0).UTC().Unix()},

		{field: data.NewField("nullable int64", nil, []*int64{utils.P(int64(1))}), wantV: float64(1)},
		{field: data.NewField("nullable int32", nil, []*int32{utils.P(int32(1))}), wantV: float64(1)},
		{field: data.NewField("nullable int16", nil, []*int16{utils.P(int16(1))}), wantV: float64(1)},
		{field: data.NewField("nullable int8", nil, []*int8{utils.P(int8(1))}), wantV: float64(1)},
		{field: data.NewField("nullable float64", nil, []*float64{utils.P(float64(1))}), wantV: float64(1)},
		{field: data.NewField("nullable float32", nil, []*float32{utils.P(float32(1))}), wantV: float64(1)},
		{field: data.NewField("nullable bool", nil, []*bool{utils.P(true)}), wantV: true},
		{field: data.NewField("nullable string", nil, []*string{utils.P("test")}), wantV: "test"},
		{field: data.NewField("nullable time", nil, []*time.Time{utils.P(time.Unix(1704067200, 0).UTC())}), wantV: time.Unix(1704067200, 0).UTC().Unix()},

		{field: data.NewField("null int64", nil, []*int64{nil}), wantV: nil},
		{field: data.NewField("null int32", nil, []*int32{nil}), wantV: nil},
		{field: data.NewField("null int16", nil, []*int16{nil}), wantV: nil},
		{field: data.NewField("null int8", nil, []*int8{nil}), wantV: nil},
		{field: data.NewField("null float64", nil, []*float64{nil}), wantV: nil},
		{field: data.NewField("null float32", nil, []*float32{nil}), wantV: nil},
		{field: data.NewField("null bool", nil, []*bool{nil}), wantV: nil},
		{field: data.NewField("null string", nil, []*string{nil}), wantV: nil},
		{field: data.NewField("null time", nil, []*time.Time{nil}), wantV: int(0)},
	}
	for ti, tt := range tests {
		t.Run(fmt.Sprintf("%v %v", ti, tt.field.Name), func(t *testing.T) {
			require.Equal(t, tt.wantV, transformations.GetNormalizedValueForExpressionEvaluation(tt.field, 0))
		})
	}
}
