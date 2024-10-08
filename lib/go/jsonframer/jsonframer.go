package jsonframer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/gframer"
	"github.com/tidwall/gjson"
	jsonata "github.com/xiatechs/jsonata-go"
)

type FramerType string

const (
	FramerTypeGJSON FramerType = "gjson"
)

type FrameFormat string

const (
	FrameFormatTable      FrameFormat = "table"
	FrameFormatTimeSeries FrameFormat = "timeseries"
	FrameFormatNumeric    FrameFormat = "numeric"
)

type FramerOptions struct {
	FramerType      FramerType // `gjson`
	FrameName       string
	RootSelector    string
	Columns         []ColumnSelector
	OverrideColumns []ColumnSelector
	FrameFormat     FrameFormat
}

type ColumnSelector struct {
	Selector   string
	Alias      string
	Type       string
	TimeFormat string
}

func validateJson(jsonString string) (err error) {
	if strings.TrimSpace(jsonString) == "" {
		return errors.Join(errors.New("empty json received"), ErrInvalidJSONContent)
	}
	if !gjson.Valid(jsonString) {
		return errors.Join(errors.New("invalid json response received"), ErrInvalidJSONContent)
	}
	return err
}

func ToFrames(jsonString string, options FramerOptions) (frames []*data.Frame, err error) {
	err = validateJson(jsonString)
	if err != nil {
		return frames, err
	}
	switch options.FramerType {
	default:
		outString, err := GetRootData(jsonString, options.RootSelector)
		if err != nil {
			return frames, err
		}
		outString, err = getColumnValuesFromResponseString(outString, options.Columns)
		if err != nil {
			return frames, err
		}
		result := gjson.Parse(outString)
		if result.IsArray() {
			nonArrayItemsFound := false
			for _, item := range result.Array() {
				if item.Exists() && !item.IsArray() {
					nonArrayItemsFound = true
				}
			}
			if nonArrayItemsFound {
				frame, err := getFrameFromResponseString(outString, options)
				if err != nil {
					return frames, err
				}
				frames = append(frames, frame)
				return frames, err
			}
			for _, v := range result.Array() {
				frame, err := getFrameFromResponseString(v.Raw, options)
				if err != nil {
					return frames, err
				}
				if frame != nil {
					if options.FrameFormat == FrameFormatTimeSeries && frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
						frame, err = data.LongToWide(frame, nil)
						if err != nil {
							return frames, err
						}
					}
					frames = append(frames, frame)
				}
			}
			if options.FrameFormat == FrameFormatTimeSeries && len(frames) > 1 {
				for k := range frames {
					if frames[k].Meta == nil {
						frames[k].Meta = &data.FrameMeta{}
					}
					frames[k].Meta.Type = data.FrameTypeTimeSeriesMulti
					frames[k].Meta.TypeVersion = data.FrameTypeVersion{0, 1}
				}
			}
			if options.FrameFormat == FrameFormatNumeric && len(frames) > 1 {
				for k := range frames {
					if frames[k].Meta == nil {
						frames[k].Meta = &data.FrameMeta{}
					}
					frames[k].Meta.Type = data.FrameTypeNumericMulti
					frames[k].Meta.TypeVersion = data.FrameTypeVersion{0, 1}
				}
			}
			return frames, err
		}
		frame, err := getFrameFromResponseString(outString, options)
		if err != nil {
			return frames, err
		}
		if frame != nil {
			if options.FrameFormat == FrameFormatTimeSeries && frame.TimeSeriesSchema().Type == data.TimeSeriesTypeLong {
				frame, err = data.LongToWide(frame, nil)
				if err != nil {
					return frames, err
				}
			}
			frames = append(frames, frame)
		}
	}
	return frames, err
}

func ToFrame(jsonString string, options FramerOptions) (frame *data.Frame, err error) {
	err = validateJson(jsonString)
	if err != nil {
		return frame, err
	}
	outString, err := GetRootData(jsonString, options.RootSelector)
	if err != nil {
		return frame, err
	}
	outString, err = getColumnValuesFromResponseString(outString, options.Columns)
	if err != nil {
		return frame, err
	}
	return getFrameFromResponseString(outString, options)
}

func GetRootData(jsonString string, rootSelector string) (string, error) {
	if rootSelector != "" {
		r := gjson.Get(string(jsonString), rootSelector)
		if r.Exists() {
			return r.String(), nil
		}
		expr, err := jsonata.Compile(rootSelector)
		if err != nil {
			return "", errors.Join(ErrInvalidRootSelector, err)
		}
		if expr == nil {
			return "", errors.Join(ErrInvalidRootSelector)
		}
		var data any
		err = json.Unmarshal([]byte(jsonString), &data)
		if err != nil {
			return "", errors.Join(ErrInvalidJSONContent, err)
		}
		res, err := expr.Eval(data)
		if err != nil {
			return "", errors.Join(ErrEvaluatingJSONata, err)
		}
		r2, err := json.Marshal(res)
		if err != nil {
			return "", errors.Join(ErrInvalidJSONContent, err)
		}
		return string(r2), nil
	}
	return jsonString, nil
}

func getColumnValuesFromResponseString(responseString string, columns []ColumnSelector) (string, error) {
	if len(columns) > 0 {
		outString := responseString
		result := gjson.Parse(outString)
		out := []map[string]interface{}{}
		if result.IsArray() {
			result.ForEach(func(key, value gjson.Result) bool {
				oi := map[string]interface{}{}
				for _, col := range columns {
					name := col.Alias
					if name == "" {
						name = col.Selector
					}
					oi[name] = convertFieldValueType(gjson.Get(value.Raw, col.Selector).Value(), col)
				}
				out = append(out, oi)
				return true
			})
		}
		if !result.IsArray() && result.IsObject() {
			oi := map[string]interface{}{}
			for _, col := range columns {
				name := col.Alias
				if name == "" {
					name = col.Selector
				}
				oi[name] = convertFieldValueType(gjson.Get(result.Raw, col.Selector).Value(), col)
			}
			out = append(out, oi)
		}
		a, err := json.Marshal(out)
		if err != nil {
			return "", errors.Join(err, ErrInvalidJSONContent)
		}
		return string(a), nil
	}
	return responseString, nil
}

func getFrameFromResponseString(responseString string, options FramerOptions) (frame *data.Frame, err error) {
	var out interface{}
	err = json.Unmarshal([]byte(responseString), &out)
	if err != nil {
		return frame, errors.Join(fmt.Errorf("error while un-marshaling response. %s", err.Error()), ErrInvalidJSONContent)
	}
	columns := []gframer.ColumnSelector{}
	for _, c := range options.Columns {
		columns = append(columns, gframer.ColumnSelector{
			Alias:      c.Alias,
			Selector:   c.Selector,
			Type:       c.Type,
			TimeFormat: c.TimeFormat,
		})
	}
	overrides := []gframer.ColumnSelector{}
	for _, c := range options.OverrideColumns {
		overrides = append(overrides, gframer.ColumnSelector{
			Alias:      c.Alias,
			Selector:   c.Selector,
			Type:       c.Type,
			TimeFormat: c.TimeFormat,
		})
	}
	frame, err = gframer.ToDataFrame(out, gframer.FramerOptions{
		FrameName:       options.FrameName,
		Columns:         columns,
		OverrideColumns: overrides,
	})
	if frame != nil {
		if frame.Meta == nil {
			frame.Meta = &data.FrameMeta{}
		}
		if options.FrameFormat == FrameFormatTimeSeries {
			frame.Meta.Type = data.FrameTypeTimeSeriesWide
			frame.Meta.TypeVersion = data.FrameTypeVersion{0, 1}
		}
		if options.FrameFormat == FrameFormatNumeric {
			frame.Meta.Type = data.FrameTypeNumericLong
			frame.Meta.TypeVersion = data.FrameTypeVersion{0, 1}
		}
	}
	return frame, err
}

func convertFieldValueType(input interface{}, _ ColumnSelector) interface{} {
	return input
}
