package gframer

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/grafana/grafana-plugin-sdk-go/data"
	"github.com/grafana/infinity-libs/lib/go/utils"
)

func anyToNullableString(input []any, fieldName string, labels data.Labels, o []any) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableString, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch cvt := currentValue.(type) {
		case string:
			field.Set(i, pointer(currentValue.(string)))
		case float64, float32, int, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			field.Set(i, pointer(fmt.Sprintf("%v", currentValue)))
		case bool:
			field.Set(i, pointer(fmt.Sprintf("%v", currentValue.(bool))))
		default:
			noOperation(cvt)
			field.Set(i, nil)
		}
	}
	return field
}

func anyToNullableBool(input []any, fieldName string, labels data.Labels, o []interface{}) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableBool, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch cvt := currentValue.(type) {
		case bool:
			if val, ok := (currentValue.(bool)); ok {
				field.Set(i, pointer(val))
			}
		case string:
			val, ok := (currentValue.(string))
			if ok && strings.ToLower(val) == "true" {
				field.Set(i, pointer(true))
			}
		default:
			noOperation(cvt)
			field.Set(i, nil)
		}
	}
	return field
}

func anyToNullableNumber(input []any, fieldName string, labels data.Labels, o []interface{}) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableFloat64, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch cvt := currentValue.(type) {
		case string:
			if item, err := strconv.ParseFloat(currentValue.(string), 64); err == nil {
				field.Set(i, pointer(item))
			}
		case float64:
			field.Set(i, pointer(currentValue.(float64)))
		case float32:
			field.Set(i, pointer(float64(currentValue.(float32))))
		case int64:
			field.Set(i, pointer(float64(currentValue.(int64))))
		case int32:
			field.Set(i, pointer(float64(currentValue.(int32))))
		case int16:
			field.Set(i, pointer(float64(currentValue.(int16))))
		case int8:
			field.Set(i, pointer(float64(currentValue.(int8))))
		case int:
			field.Set(i, pointer(float64(currentValue.(int))))
		default:
			noOperation(cvt)
			field.Set(i, nil)
		}
	}
	return field
}

func anyToNullableTimestamp(input []any, fieldName string, labels data.Labels, o []interface{}, timeFormat string) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch a := currentValue.(type) {
		case float64:
			if v := fmt.Sprintf("%.0f", currentValue); v != "" {
				format := "2006"
				if timeFormat != "" {
					format = timeFormat
				}
				if t, err := time.Parse(format, v); err == nil {
					field.Set(i, pointer(t))
				}
			}
		case string:
			if currentValue.(string) != "" {
				field.Set(i, utils.GetTimeFromString(currentValue.(string), timeFormat))
			}
		default:
			noOperation(a)
			field.Set(i, nil)
		}
	}
	return field
}

func anyToNullableTimestampEpoch(input []any, fieldName string, labels data.Labels, o []interface{}) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch cvt := currentValue.(type) {
		case string:
			if item, err := strconv.ParseInt(currentValue.(string), 10, 64); err == nil && currentValue.(string) != "" {
				field.Set(i, pointer(time.UnixMilli(item)))
			}
		case float64:
			field.Set(i, pointer(time.UnixMilli(int64(currentValue.(float64)))))
		default:
			noOperation(cvt)
			field.Set(i, nil)
		}
	}
	return field
}

func anyToNullableTimestampEpochSecond(input []any, fieldName string, labels data.Labels, o []interface{}) *data.Field {
	field := data.NewFieldFromFieldType(data.FieldTypeNullableTime, len(input))
	field.Name = fieldName
	field.Labels = labels
	for i := 0; i < len(input); i++ {
		currentValue := o[i]
		switch cvt := currentValue.(type) {
		case string:
			if item, err := strconv.ParseInt(currentValue.(string), 10, 64); err == nil && currentValue.(string) != "" {
				field.Set(i, pointer(time.Unix(item, 0)))
			}
		case float64:
			field.Set(i, pointer(time.Unix(int64(currentValue.(float64)), 0)))
		default:
			noOperation(cvt)
			field.Set(i, nil)
		}
	}
	return field
}
