package jsonframer_test

import (
	"encoding/json"
	"os"
	"strings"
	"testing"

	"github.com/grafana/grafana-plugin-sdk-go/experimental"
	"github.com/grafana/infinity-libs/lib/go/jsonframer"
	"github.com/stretchr/testify/require"
)

func TestToFrame(t *testing.T) {
	t.Run("jq", func(t *testing.T) {
		t.Run(("downstream tests"), func(t *testing.T) {
			// All these test cases are validated against https://jqlang.org/tutorial/
			fileContent, err := os.ReadFile("./testdata/jq/github_jqlang_jq_commits.json")
			require.Nil(t, err)
			t.Run("no root selector", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: "",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("dot as root selector", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: ".",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("first commit", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: ".[0]",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("first commit with selected fields", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: ".[0] | {message: .commit.message, name: .commit.committer.name}",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("all commits with selected fields", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: ".[] | {message: .commit.message, name: .commit.committer.name}",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("all commits with selected fields into an array", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: "[.[] | {message: .commit.message, name: .commit.committer.name}]",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("parent commits", func(t *testing.T) {
				options := jsonframer.FramerOptions{
					FramerType:   jsonframer.FramerTypeJQ,
					RootSelector: "[.[] | {message: .commit.message, name: .commit.committer.name, parents: [.parents[].html_url]}]",
				}
				var out interface{}
				err = json.Unmarshal(fileContent, &out)
				require.Nil(t, err)
				gotFrame, err := jsonframer.ToFrame(string(fileContent), options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
		})
		t.Run("basic tests", func(t *testing.T) {
			t.Run("simple object", func(t *testing.T) {
				jsonContent := `{ "foo" : "foo1", "bar" : "bar1" }`
				options := jsonframer.FramerOptions{FramerType: jsonframer.FramerTypeJQ, RootSelector: "."}
				gotFrame, err := jsonframer.ToFrame(jsonContent, options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("array of object", func(t *testing.T) {
				jsonContent := `[{ "foo" : "foo1", "bar" : "bar1" }]`
				options := jsonframer.FramerOptions{FramerType: jsonframer.FramerTypeJQ, RootSelector: "."}
				gotFrame, err := jsonframer.ToFrame(jsonContent, options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
			t.Run("array of objects", func(t *testing.T) {
				jsonContent := `[{ "foo" : "foo1", "bar" : "bar1" }, { "foo" : "foo2", "bar" : "bar2" }]`
				options := jsonframer.FramerOptions{FramerType: jsonframer.FramerTypeJQ, RootSelector: "."}
				gotFrame, err := jsonframer.ToFrame(jsonContent, options)
				require.Nil(t, err)
				require.NotNil(t, gotFrame)
				experimental.CheckGoldenJSONFrame(t, "testdata/jq", strings.ReplaceAll(strings.ReplaceAll(t.Name(), "TestToFrame/jq/", ""), " ", ""), gotFrame, true)
			})
		})
	})
}
