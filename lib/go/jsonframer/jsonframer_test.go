package jsonframer_test

import (
	"testing"

	"github.com/grafana/infinity-libs/lib/go/jsonframer"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var sampleData = map[string]string{
	"user": `{	
		"name":"foo", 	
		"age": 30,
		"address" : {
			"line1" 	: "123, ABC street",
			"line2" 	: "Foo apartment",
			"country" 	: "Bar",
			"postcode" 	: "ABC123"
		}	
	}`,
	"users": `[ 
		{	
			"name":"foo", 	
			"age": 30	
		}, 
		{	
			"name":"bar", 	
			"age": 14	
		} 
	]`,
	"nested": `{
		"meta": {
			"foo" : "bar"
		},
		"data": [ 
			{	
				"name":"foo", 	
				"age": 30	
			}, 
			{	
				"name":"bar", 	
				"age": 14	
			} 
		]
	}`,
}

func TestGetRootData(t *testing.T) {
	t.Run("jq", func(t *testing.T) {
		tests := []struct {
			name         string
			jsonString   string
			rootSelector string
			want         string
			wantErr      error
		}{
			{
				name:         "should parse json object",
				jsonString:   sampleData["user"],
				rootSelector: `.`,
				want:         "[{\"address\":{\"country\":\"Bar\",\"line1\":\"123, ABC street\",\"line2\":\"Foo apartment\",\"postcode\":\"ABC123\"},\"age\":30,\"name\":\"foo\"}]",
			},
			{
				name:         "should parse json object into array",
				jsonString:   sampleData["user"],
				rootSelector: `[.]`,
				want:         "[{\"address\":{\"country\":\"Bar\",\"line1\":\"123, ABC street\",\"line2\":\"Foo apartment\",\"postcode\":\"ABC123\"},\"age\":30,\"name\":\"foo\"}]",
			},
			{
				name:         "should parse json object and extract field using dot syntax",
				jsonString:   sampleData["user"],
				rootSelector: `.address.postcode`,
				want:         `["ABC123"]`,
			},
			{
				name:         "should parse json object and extract field using pipe syntax",
				jsonString:   sampleData["user"],
				rootSelector: `.address | .postcode`,
				want:         `["ABC123"]`,
			},
			{
				name:         "should parse json array and extract field",
				jsonString:   sampleData["users"],
				rootSelector: `.[] |  .name`,
				want:         `["foo","bar"]`,
			},
			{
				name:         "should parse json array and manipulate items",
				jsonString:   sampleData["users"],
				rootSelector: `.[] |  { "username" : .name , "age_after_30y" : .age + 30 }`,
				want:         `[{"age_after_30y":60,"username":"foo"},{"age_after_30y":44,"username":"bar"}]`,
			},
			{
				name:         "should parse json array and manipulate items into array",
				jsonString:   sampleData["users"],
				rootSelector: `[.[] |  { "username" : .name , "age_after_30y" : .age + 30 }]`,
				want:         `[{"age_after_30y":60,"username":"foo"},{"age_after_30y":44,"username":"bar"}]`,
			},
			{
				name:         "should parse nested json",
				jsonString:   sampleData["nested"],
				rootSelector: `.data`,
				want:         `[{"age":30,"name":"foo"},{"age":14,"name":"bar"}]`,
			},
			{
				name:         "should parse nested json into array",
				jsonString:   sampleData["nested"],
				rootSelector: `.data[]`,
				want:         `[{"age":30,"name":"foo"},{"age":14,"name":"bar"}]`,
			},
			{
				name:         "should parse nested json with conditional statement",
				jsonString:   sampleData["nested"],
				rootSelector: `.data[] | { "name" : .name, "can_vote" : (if .age > 18 then "yes" else "no" end) }`,
				want:         `[{"can_vote":"yes","name":"foo"},{"can_vote":"no","name":"bar"}]`,
			},
		}
		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got, err := jsonframer.GetRootData(tt.jsonString, tt.rootSelector, jsonframer.FramerTypeJQ)
				if tt.wantErr != nil {
					require.NotNil(t, err)
					assert.Equal(t, tt.wantErr, err)
					return
				}
				require.Nil(t, err)
				require.NotNil(t, got)
				assert.Equal(t, tt.want, got)
			})
		}
	})
}
