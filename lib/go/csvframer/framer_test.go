package csvframer_test

import (
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/infinity-libs/lib/go/csvframer"
	"github.com/grafana/infinity-libs/lib/go/gframer"
	"github.com/grafana/infinity-libs/lib/go/jsonframer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCsvStringToFrame(t *testing.T) {
	tests := []struct {
		name      string
		csvString string
		options   csvframer.FramerOptions
		wantError error
	}{
		{
			name:      "empty csv should return error",
			wantError: csvframer.ErrEmptyCsv,
		},
		{
			name:      "valid csv should not return error",
			csvString: strings.Join([]string{`a,b,c`, `1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
		},
		{
			name:      "valid csv without headers should not return error",
			csvString: strings.Join([]string{`1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
			options:   csvframer.FramerOptions{NoHeaders: true},
		},
		{
			name:      "framer options should be respected",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12	13`, `21	22	23`}, "\n"),
			options: csvframer.FramerOptions{FrameName: "foo", Delimiter: "\t", RelaxColumnCount: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch"},
			}},
		},
		{
			name:      "relax column count",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12`, `21	22	23`}, "\n"),
			options: csvframer.FramerOptions{FrameName: "foo", Delimiter: "\t", SkipLinesWithError: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch"},
			}},
		},
		{
			name:      "Skip empty lines",
			csvString: strings.Join([]string{`a	b	c`, `1	2	3`, ``, `21	22	23`}, "\n"),
			options: csvframer.FramerOptions{FrameName: "foo", Delimiter: "\t", Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "timestamp_epoch_s"},
			}},
		},
		{
			name:      "relax column count",
			csvString: strings.Join([]string{`a;b;c`, `1;2;3`, `11;13`, `21;22;23`}, "\n"),
			options: csvframer.FramerOptions{FrameName: "foo", Delimiter: ";", RelaxColumnCount: true, Columns: []gframer.ColumnSelector{
				{Selector: "a", Alias: "A", Type: "number"},
				{Selector: "b", Alias: "b", Type: "string"},
				{Selector: "c", Type: "string"},
			}},
		},
		{
			name:      "comment",
			csvString: strings.Join([]string{`# foo`, `a,b,c`, `#01,02,03`, `1,2,3`, `11,12,13`, `21,22,23`, `#`}, "\n"),
			options:   csvframer.FramerOptions{Comment: "#"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotFrame, err := csvframer.ToFrame(tt.csvString, tt.options)
			if tt.wantError != nil {
				require.NotNil(t, err)
				assert.Equal(t, tt.wantError, err)
				return
			}
			require.Nil(t, err)
			require.NotNil(t, gotFrame)
			if tt.wantError == nil {
				goldenFileName := strings.Replace(t.Name(), "TestCsvStringToFrame/", "", 1)
				experimental.CheckGoldenJSONFrame(t, "testdata", goldenFileName, gotFrame, false)
			}
		})
	}
	t.Run("with root selector and JSONata framer", func(t *testing.T) {
		tests := []struct {
			name      string
			csvString string
			options   csvframer.FramerOptions
			wantError error
		}{
			{
				name:      "empty csv should return error",
				wantError: csvframer.ErrEmptyCsv,
			},
			{
				name:      "valid csv should not return error",
				csvString: strings.Join([]string{`a,b,c`, `1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
			},
			{
				name:      "valid csv without headers should not return error",
				csvString: strings.Join([]string{`1,2,3`, `11,12,13`, `21,22,23`}, "\n"),
				options: csvframer.FramerOptions{
					NoHeaders:  true,
					FramerType: jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
						"A": $number($lookup($v,"1")),
						"B": $number($lookup($v,"2")),
						"C": $number($lookup($v,"3"))
					}})`},
			},
			{
				name:      "framer options should be respected",
				csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12	13`, `21	22	23`}, "\n"),
				options: csvframer.FramerOptions{
					FrameName:        "foo",
					Delimiter:        "\t",
					RelaxColumnCount: true,
					FramerType:       jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
							"a": $number($lookup($v,"a")),
							"b": $number($lookup($v,"b")),
							"c": $number($lookup($v,"c"))
					}})`,
					Columns: []gframer.ColumnSelector{
						{Selector: "a", Alias: "A", Type: "number"},
						{Selector: "b", Alias: "b", Type: "string"},
						{Selector: "c", Type: "timestamp_epoch"},
					}},
			},
			{
				name:      "relax column count",
				csvString: strings.Join([]string{`a	b	c`, `1	2	3`, `11	12`, `21	22	23`}, "\n"),
				options: csvframer.FramerOptions{
					FrameName:          "foo",
					Delimiter:          "\t",
					SkipLinesWithError: true,
					FramerType:         jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
							"a": $number($lookup($v,"a")) + 100,
							"b": $number($lookup($v,"b")),
							"c": $number($lookup($v,"c"))
					}})`,
					Columns: []gframer.ColumnSelector{
						{Selector: "a", Alias: "A", Type: "number"},
						{Selector: "b", Alias: "b", Type: "string"},
						{Selector: "c", Type: "timestamp_epoch"},
					}},
			},
			{
				name:      "Skip empty lines",
				csvString: strings.Join([]string{`a	b	c`, `1	2	3`, ``, `21	22	23`}, "\n"),
				options: csvframer.FramerOptions{
					FrameName:  "foo",
					Delimiter:  "\t",
					FramerType: jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
							"a": $number($lookup($v,"a")) + 100,
							"b": $number($lookup($v,"b")),
							"c": $number($lookup($v,"c"))
					}})`,
					Columns: []gframer.ColumnSelector{
						{Selector: "a", Alias: "A", Type: "number"},
						{Selector: "b", Alias: "b", Type: "string"},
						{Selector: "c", Type: "timestamp_epoch_s"},
					}},
			},
			{
				name:      "relax column count with custom delimiter",
				csvString: strings.Join([]string{`a;b;c`, `1;2;3`, `11;13`, `21;22;23`}, "\n"),
				options: csvframer.FramerOptions{
					FrameName:        "foo",
					Delimiter:        ";",
					RelaxColumnCount: true,
					FramerType:       jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
							"a": $number($lookup($v,"a")) +100,
							"b": $number($lookup($v,"b")),
							"c": $exists($number($lookup($v,"c"))) ? $number($lookup($v,"c")): null
					}})`,
					Columns: []gframer.ColumnSelector{
						{Selector: "a", Alias: "A", Type: "number"},
						{Selector: "b", Alias: "b", Type: "string"},
						{Selector: "c", Type: "string"},
					},
				},
				wantError: jsonframer.ErrEvaluatingJSONata, // TODO: bug in jsonata-go library handling missing column
			},
			{
				name:      "comment",
				csvString: strings.Join([]string{`# foo`, `a,b,c`, `#01,02,03`, `1,2,3`, `11,12,13`, `21,22,23`, `#`}, "\n"),
				options: csvframer.FramerOptions{
					Comment:    "#",
					FramerType: jsonframer.FramerTypeJsonata,
					RootSelector: `$map($,function($v){{
							"a": $number($lookup($v,"a")) + 100,
							"b": $number($lookup($v,"b")),
							"c": $number($lookup($v,"c"))
					}})`},
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				gotFrame, err := csvframer.ToFrame(tt.csvString, tt.options)
				if tt.wantError != nil {
					require.NotNil(t, err)
					assert.Equal(t, tt.wantError, err)
					return
				}
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				if tt.wantError == nil {
					goldenFileName := strings.Replace(t.Name(), "TestCsvStringToFrame/", "", 1)
					experimental.CheckGoldenJSONFrame(t, "testdata", goldenFileName, gotFrame, true)
				}
			})
		}
	})
}
