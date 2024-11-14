package transformations_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/transformations"
	"github.com/grafana/infinity-libs/lib/go/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSummaryFrame(t *testing.T) {
	sampleDataFrame := data.NewFrame(
		"hello",
		data.NewField("group", nil, []string{"A", "B", "A"}),
		data.NewField("id", nil, []int64{3, 4, 5}),
		data.NewField("value", nil, []int64{6, 7, 8}),
	).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"})
	sampleDataFrameWithNull := data.NewFrame(
		"hello",
		data.NewField("group", nil, []string{"A", "B", "A"}),
		data.NewField("id", nil, []int64{3, 4, 5}),
		data.NewField("value", nil, []*int64{utils.P(int64(6)), utils.P(int64(7)), nil}),
	).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"})
	tests := []struct {
		name       string
		frame      *data.Frame
		expression string
		by         string
		alias      string
		want       *data.Frame
		wantErr    bool
	}{
		{
			name: "summarize",
			frame: data.NewFrame(
				"response",
				data.NewField("sex", nil, []string{"m", "m", "f", "f", "m"}),
				data.NewField("mass", nil, []*float64{toFP(1), toFP(2), toFP(3), toFP(4), toFP(5)}),
			),
			expression: "mean(mass)",
			want: data.NewFrame(
				"response",
				data.NewField("mean(mass)", nil, []*float64{toFP(3)}),
			),
		},
		{
			name: "summarize with null",
			frame: data.NewFrame(
				"response",
				data.NewField("sex", nil, []string{"m", "m", "f", "f", "m"}),
				data.NewField("mass", nil, []*float64{toFP(1), nil, toFP(3), toFP(4), toFP(5)}),
			),
			expression: "mean(mass)",
			alias:      "summary",
			want: data.NewFrame(
				"response",
				data.NewField("summary", nil, []*float64{toFP(2.6)}),
			),
		},
		{
			name: "summarize with empty data",
			frame: data.NewFrame(
				"hello",
			),
			expression: "mean(value)",
			want: data.NewFrame(
				"hello",
				data.NewField("mean(value)", nil, []*float64{}),
			),
		},
		{
			name: "summarize with empty data and alias",
			frame: data.NewFrame(
				"hello",
			),
			expression: "mean(value)",
			alias:      "summary",
			want: data.NewFrame(
				"hello",
				data.NewField("summary", nil, []*float64{}),
			),
		},
		{
			name: "summarize with empty data and summarize by",
			frame: data.NewFrame(
				"hello",
			),
			expression: "mean(value)",
			by:         "world",
			want: data.NewFrame(
				"hello",
				data.NewField("mean(value)", nil, []*float64{}),
				data.NewField("world", nil, []*string{}),
			),
		},
		{
			name: "summarize with empty data and alias and summarize by",
			frame: data.NewFrame(
				"hello",
			),
			expression: "mean(value)",
			alias:      "summary",
			by:         "world",
			want: data.NewFrame(
				"hello",
				data.NewField("summary", nil, []*float64{}),
				data.NewField("world", nil, []*string{}),
			),
		},
		{
			name: "summarize with empty frame and summarize by",
			frame: data.NewFrame(
				"world population",
				data.NewField("country", nil, []*string{}),
				data.NewField("population", nil, []*float64{}),
			),
			expression: "mean(population)",
			by:         "country",
			want: data.NewFrame(
				"world population",
				data.NewField("mean(population)", nil, []*float64{}),
				data.NewField("country", nil, []*string{}),
			),
		},
		{
			name: "summarize with empty frame and alias and summarize by",
			frame: data.NewFrame(
				"world population",
				data.NewField("country", nil, []*string{}),
				data.NewField("population", nil, []*float64{}),
			),
			expression: "mean(population)",
			alias:      "mean population",
			by:         "country",
			want: data.NewFrame(
				"world population",
				data.NewField("mean population", nil, []*float64{}),
				data.NewField("country", nil, []*string{}),
			),
		},
		{
			name:       "actual data with summarize",
			frame:      sampleDataFrame,
			expression: "sum(value)",
			want: data.NewFrame(
				"hello",
				data.NewField("sum(value)", nil, []*float64{utils.P(float64(21))}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:       "actual data with summarize and by",
			frame:      sampleDataFrame,
			expression: "sum(value)",
			by:         "group",
			alias:      "sum of value",
			want: data.NewFrame(
				"summary",
				data.NewField("group", nil, []*string{utils.P(string("A")), utils.P(string("B"))}),
				data.NewField("sum of value", nil, []*float64{utils.P(float64(14)), utils.P(float64(7))}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
		{
			name:       "actual data with summarize and by",
			frame:      sampleDataFrameWithNull,
			expression: "sum(value)",
			by:         "group",
			alias:      "sum of value",
			want: data.NewFrame(
				"summary",
				data.NewField("group", nil, []*string{utils.P(string("A")), utils.P(string("B"))}),
				data.NewField("sum of value", nil, []*float64{utils.P(float64(6)), utils.P(float64(7))}),
			).SetMeta(&data.FrameMeta{PreferredVisualizationPluginID: "text"}),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := transformations.GetSummaryFrame(tt.frame, tt.expression, tt.by, tt.alias)
			require.Nil(t, err)
			require.NotNil(t, got)
			assert.Equal(t, tt.want, got)
		})
	}
}
