package gframer_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/infinity-libs/lib/go/gframer"
	"github.com/stretchr/testify/require"
)

func TestToDataFrame(t *testing.T) {
	updateGoldenText := false
	t.Run("nil", func(t *testing.T) {
		var input interface{}
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: ""}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("string", func(t *testing.T) {
		input := `foo`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "foo"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
		require.Nil(t, err)
	})
	t.Run("integer", func(t *testing.T) {
		input := 21
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "21"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("float", func(t *testing.T) {
		input := 21.43
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "21.43"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("bool", func(t *testing.T) {
		input := true
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: "true"}
		gotFrame, err := gframer.ToDataFrame(input, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("object", func(t *testing.T) {
		input := `{ "name":"foo", "age": 12, "hobbies":["cricket","music"], "isPrimeUser": true, "fullname": { "first": "foo", "last":"bar" } }`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("string-array", func(t *testing.T) {
		input := `["foo","bar"]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("number-array", func(t *testing.T) {
		input := `[12,14.56,0,30]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("array-inside-array", func(t *testing.T) {
		input := `[["one","two"],["three"]]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
	t.Run("all-null-array", func(t *testing.T) {
		input := `[null,null]`
		options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: input}
		var out interface{}
		err := json.Unmarshal([]byte(input), &out)
		require.Nil(t, err)
		gotFrame, err := gframer.ToDataFrame(out, options)
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata", "structs/"+strings.ReplaceAll(t.Name(), "TestToDataFrame/", ""), gotFrame, updateGoldenText)
	})
}

func TestToDataFrameSlices(t *testing.T) {
	updateGoldenText := false
	files, err := os.ReadDir("./testdata/slices")
	if err != nil {
		require.Nil(t, err)
	}
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".json") {
			t.Run(f.Name(), func(t *testing.T) {
				fileContent, err := os.ReadFile("./testdata/slices/" + f.Name())
				require.Nil(t, err)
				options := gframer.FramerOptions{FrameName: t.Name(), ExecutedQueryString: ""}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := gframer.ToDataFrame(out, options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				goldenFileName := strings.TrimSuffix(f.Name(), ".json")
				experimental.CheckGoldenJSONFrame(t, "testdata/slices", goldenFileName, gotFrame, updateGoldenText)
			})
		}
	}
}

func TestJsonFieldType(t *testing.T) {
	t.Run("mixed json content without null values", func(t *testing.T) {
		gotFrame, err := gframer.ToDataFrame([]any{
			map[string]any{"num": int64(1), "str": "two", "json": []int{3}, "bool": false},
			map[string]any{"num": int64(11), "str": "two-two", "json": map[string]any{"something": "else"}, "bool": true},
		}, gframer.FramerOptions{Columns: []gframer.ColumnSelector{{Selector: "num", Type: "number"}, {Selector: "str", Type: "string"}, {Selector: "json", Type: "json"}, {Selector: "bool", Type: "boolean"}}})
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata/jsonfield", strings.ReplaceAll(t.Name(), "TestJsonFieldType/", ""), gotFrame, true)
	})
	t.Run("mixed json content without null values and override columns", func(t *testing.T) {
		gotFrame, err := gframer.ToDataFrame([]any{
			map[string]any{"num": int64(1), "str": "two", "json": []int{3}, "bool": false},
			map[string]any{"num": int64(11), "str": "two-two", "json": map[string]any{"something": "else"}, "bool": true},
		}, gframer.FramerOptions{OverrideColumns: []gframer.ColumnSelector{{Selector: "num", Type: "json"}, {Selector: "str", Type: "json"}, {Selector: "json", Type: "json"}, {Selector: "bool", Type: "json"}}})
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata/jsonfield", strings.ReplaceAll(t.Name(), "TestJsonFieldType/", ""), gotFrame, true)
	})
	t.Run("mixed json content with null values", func(t *testing.T) {
		gotFrame, err := gframer.ToDataFrame([]any{
			map[string]any{"num": int64(1), "str": "two", "json": []int{3}, "bool": false},
			map[string]any{"num": nil, "str": nil, "json": nil, "bool": nil},
		}, gframer.FramerOptions{Columns: []gframer.ColumnSelector{{Selector: "num", Type: "number"}, {Selector: "str", Type: "string"}, {Selector: "json", Type: "json"}, {Selector: "bool", Type: "boolean"}}})
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata/jsonfield", strings.ReplaceAll(t.Name(), "TestJsonFieldType/", ""), gotFrame, true)
	})
	t.Run("mixed json content with null values as first", func(t *testing.T) {
		gotFrame, err := gframer.ToDataFrame([]any{
			map[string]any{"num": nil, "str": nil, "json": nil, "bool": nil},
			map[string]any{"num": int64(1), "str": "two", "json": []int{3}, "bool": false},
		}, gframer.FramerOptions{Columns: []gframer.ColumnSelector{{Selector: "num", Type: "number"}, {Selector: "str", Type: "string"}, {Selector: "json", Type: "json"}, {Selector: "bool", Type: "boolean"}}})
		require.Nil(t, err)
		require.NotNil(t, gotFrame)
		experimental.CheckGoldenJSONFrame(t, "testdata/jsonfield", strings.ReplaceAll(t.Name(), "TestJsonFieldType/", ""), gotFrame, true)
	})
}
