package transformations_test

import (
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/transformations"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetSummaryFrame(t *testing.T) {
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
