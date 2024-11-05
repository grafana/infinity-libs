package transformations_test

import (
	"errors"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/transformations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApplyFilter(t *testing.T) {
	sampleDataFrame := data.NewFrame("hello", data.NewField("group", nil, []string{"A", "B", "A"}), data.NewField("id", nil, []int64{3, 4, 5}), data.NewField("value", nil, []int64{6, 7, 8})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"})
	A := "a"
	B := "b"
	zero := int64(0)
	one := int64(1)
	tests := []struct {
		name             string
		frame            *data.Frame
		filterExpression string
		want             *data.Frame
		wantErr          error
	}{
		{
			name:             "nil frame should return nil frame",
			filterExpression: "group =='A'",
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
			want:             data.NewFrame("hello", data.NewField("group", data.Labels{}, []string{"A", "A"}), data.NewField("id", data.Labels{}, []int64{3, 5}), data.NewField("value", data.Labels{}, []int64{6, 8})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "frame with data and with filter should filter the data without matching condition",
			filterExpression: "id == 1",
			frame:            sampleDataFrame,
			want:             data.NewFrame("hello", data.NewField("group", data.Labels{}, []string{}), data.NewField("id", data.Labels{}, []int64{}), data.NewField("value", data.Labels{}, []int64{})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "frame with data and with filter should filter the data with incorrect matching condition",
			filterExpression: "group == 3",
			frame:            sampleDataFrame,
			want:             data.NewFrame("hello", data.NewField("group", data.Labels{}, []string{}), data.NewField("id", data.Labels{}, []int64{}), data.NewField("value", data.Labels{}, []int64{})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "null value filter",
			filterExpression: "value != nil",
			frame:            data.NewFrame("hello", data.NewField("name", nil, []*string{&A, &B}), data.NewField("value", nil, []*string{&A, nil})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want:             data.NewFrame("hello", data.NewField("name", data.Labels{}, []*string{&A}), data.NewField("value", data.Labels{}, []*string{&A})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:             "null value filter with number",
			filterExpression: "value != nil",
			frame:            data.NewFrame("hello", data.NewField("name", nil, []*string{&A, &B, &A}), data.NewField("value", nil, []*int64{&zero, &one, nil})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
			want:             data.NewFrame("hello", data.NewField("name", data.Labels{}, []*string{&A, &B}), data.NewField("value", data.Labels{}, []*int64{&zero, &one})).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
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
