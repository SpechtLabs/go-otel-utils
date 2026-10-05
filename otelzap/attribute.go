package otelzap

import (
	"encoding/json"
	"fmt"
	"math"
	"reflect"
	"strconv"

	"go.opentelemetry.io/otel/attribute"
)

// Attribute converts a value of any type into a span attribute under key:
// scalars and slices of scalars keep their types, a fmt.Stringer or error
// becomes its string, a uint64 past the int64 range becomes its decimal
// string, and anything else its JSON or, failing that, its fmt.Sprint form.
func Attribute(key string, value any) attribute.KeyValue {
	switch value := value.(type) {
	case nil:
		return attribute.String(key, "<nil>")
	case string:
		return attribute.String(key, value)
	case int:
		return attribute.Int(key, value)
	case int64:
		return attribute.Int64(key, value)
	case uint64:
		if value > math.MaxInt64 {
			return attribute.String(key, strconv.FormatUint(value, 10))
		}
		return attribute.Int64(key, int64(value))
	case float64:
		return attribute.Float64(key, value)
	case bool:
		return attribute.Bool(key, value)
	case fmt.Stringer:
		return attribute.String(key, value.String())
	case error:
		return attribute.String(key, value.Error())
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	// Arrays and slices whose elements are named types (type flag bool) are
	// read element by element: an array value can't be sliced, and a []flag
	// isn't a []bool.
	case reflect.Array, reflect.Slice:
		switch rv.Type().Elem().Kind() {
		case reflect.Bool:
			return attribute.BoolSlice(key, elements(rv, reflect.Value.Bool))
		case reflect.Int:
			return attribute.IntSlice(key, elements(rv, func(v reflect.Value) int { return int(v.Int()) }))
		case reflect.Int64:
			return attribute.Int64Slice(key, elements(rv, reflect.Value.Int))
		case reflect.Float64:
			return attribute.Float64Slice(key, elements(rv, reflect.Value.Float))
		case reflect.String:
			return attribute.StringSlice(key, elements(rv, reflect.Value.String))
		default:
			return attribute.KeyValue{Key: attribute.Key(key)}
		}
	case reflect.Bool:
		return attribute.Bool(key, rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return attribute.Int64(key, rv.Int())
	case reflect.Float64:
		return attribute.Float64(key, rv.Float())
	case reflect.String:
		return attribute.String(key, rv.String())
	}
	if b, err := json.Marshal(value); b != nil && err == nil {
		return attribute.String(key, string(b))
	}
	return attribute.String(key, fmt.Sprint(value))
}

// LogValue converts a value of any type into a log record value, like
// Attribute converts it into a span attribute. Slices of any element type
// become a slice value, converted element by element.
func LogValue(value any) attribute.Value {
	switch value := value.(type) {
	case nil:
		return attribute.StringValue("<nil>")
	case string:
		return attribute.StringValue(value)
	case int:
		return attribute.IntValue(value)
	case int64:
		return attribute.Int64Value(value)
	case uint64:
		if value > math.MaxInt64 {
			return attribute.StringValue(strconv.FormatUint(value, 10))
		}
		return attribute.Int64Value(int64(value))
	case float64:
		return attribute.Float64Value(value)
	case bool:
		return attribute.BoolValue(value)
	case fmt.Stringer:
		return attribute.StringValue(value.String())
	case error:
		return attribute.StringValue(value.Error())
	}

	rv := reflect.ValueOf(value)

	switch rv.Kind() {
	case reflect.Array, reflect.Slice:
		values := make([]attribute.Value, rv.Len())
		for i := range values {
			values[i] = LogValue(rv.Index(i).Interface())
		}
		return attribute.SliceValue(values...)
	case reflect.Bool:
		return attribute.BoolValue(rv.Bool())
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return attribute.Int64Value(rv.Int())
	case reflect.Float64:
		return attribute.Float64Value(rv.Float())
	case reflect.String:
		return attribute.StringValue(rv.String())
	}
	if b, err := json.Marshal(value); err == nil {
		return attribute.StringValue(string(b))
	}
	return attribute.StringValue(fmt.Sprint(value))
}

// elements reads every element of the array or slice rv with get.
func elements[T any](rv reflect.Value, get func(reflect.Value) T) []T {
	out := make([]T, rv.Len())
	for i := range out {
		out[i] = get(rv.Index(i))
	}
	return out
}
