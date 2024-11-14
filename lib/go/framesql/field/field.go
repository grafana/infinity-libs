package field

import (
	"github.com/grafana/grafana-plugin-sdk-go/data"
)

func Count(field *data.Field) (float64, error) {
	return float64(field.Len()), nil
}

func First(field *data.Field) (any, error) {
	if field.Len() == 0 {
		return nil, nil
	}
	if field.Nullable() {
		if _, ok := field.ConcreteAt(0); !ok {
			return nil, nil
		}
	}
	return field.At(0), nil
}

func Last(field *data.Field) (any, error) {
	if field.Len() == 0 {
		return nil, nil
	}
	if field.Nullable() {
		if _, ok := field.ConcreteAt(field.Len() - 1); !ok {
			return nil, nil
		}
	}
	return field.At(field.Len() - 1), nil
}

func Sum(field *data.Field) (float64, error) {
	sum := float64(0)
	for i := 0; i < field.Len(); i++ {
		if v, err := field.NullableFloatAt(i); v != nil && err == nil {
			sum += *v
		}
	}
	return sum, nil
}

func Min(field *data.Field) (float64, error) {
	var min *float64
	for i := 0; i < field.Len(); i++ {
		if v, err := field.NullableFloatAt(i); v != nil && err == nil {
			if min == nil || *v < *min {
				min = v
			}
		}
	}
	return *min, nil
}

func Max(field *data.Field) (float64, error) {
	var max *float64
	for i := 0; i < field.Len(); i++ {
		if v, err := field.NullableFloatAt(i); v != nil && err == nil {
			if max == nil || *v > *max {
				max = v
			}
		}
	}
	return *max, nil
}

func Mean(field *data.Field) (float64, error) {
	sum := float64(0)
	for i := 0; i < field.Len(); i++ {
		if v, err := field.NullableFloatAt(i); v != nil && err == nil {
			sum += *v
		}
	}
	if field.Len() < 1 {
		return 0, nil
	}
	return sum / float64(field.Len()), nil
}
