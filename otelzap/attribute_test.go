package otelzap_test

import (
	"errors"
	"math"
	"testing"

	"github.com/spechtlabs/go-otel-utils/otelzap"
	"github.com/stretchr/testify/assert"
	"go.opentelemetry.io/otel/attribute"
)

type stringer struct{}

func (stringer) String() string { return "stringer" }

type (
	flag    bool
	named   string
	counter int32
	ratio   float64
	point   struct{ X, Y int }
	// noJSON has a field encoding/json can't marshal.
	noJSON struct{ C chan int }
)

func TestAttribute(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  attribute.KeyValue
	}{
		{name: "nil", value: nil, want: attribute.String("k", "<nil>")},
		{name: "string", value: "v", want: attribute.String("k", "v")},
		{name: "int", value: 42, want: attribute.Int("k", 42)},
		{name: "int64", value: int64(42), want: attribute.Int64("k", 42)},
		{name: "uint64 in the int64 range", value: uint64(42), want: attribute.Int64("k", 42)},
		{name: "uint64 past the int64 range", value: uint64(math.MaxUint64), want: attribute.String("k", "18446744073709551615")},
		{name: "float64", value: 1.5, want: attribute.Float64("k", 1.5)},
		{name: "bool", value: true, want: attribute.Bool("k", true)},
		{name: "fmt.Stringer", value: stringer{}, want: attribute.String("k", "stringer")},
		{name: "error", value: errors.New("failed"), want: attribute.String("k", "failed")},
		{name: "bool slice", value: []bool{true, false}, want: attribute.BoolSlice("k", []bool{true, false})},
		{name: "int slice", value: []int{1, 2}, want: attribute.IntSlice("k", []int{1, 2})},
		{name: "int64 slice", value: []int64{1, 2}, want: attribute.Int64Slice("k", []int64{1, 2})},
		{name: "float64 slice", value: []float64{1.5}, want: attribute.Float64Slice("k", []float64{1.5})},
		{name: "string slice", value: []string{"a", "b"}, want: attribute.StringSlice("k", []string{"a", "b"})},
		{name: "string array", value: [2]string{"a", "b"}, want: attribute.StringSlice("k", []string{"a", "b"})},
		{name: "int array", value: [2]int{1, 2}, want: attribute.IntSlice("k", []int{1, 2})},
		{name: "slice of named bools", value: []flag{true}, want: attribute.BoolSlice("k", []bool{true})},
		{name: "slice of named strings", value: []named{"a"}, want: attribute.StringSlice("k", []string{"a"})},
		{name: "slice of another type", value: []uint16{1}, want: attribute.KeyValue{Key: "k"}},
		{name: "named bool", value: flag(true), want: attribute.Bool("k", true)},
		{name: "named int", value: counter(7), want: attribute.Int64("k", 7)},
		{name: "named string", value: named("v"), want: attribute.String("k", "v")},
		{name: "struct as JSON", value: point{X: 1, Y: 2}, want: attribute.String("k", `{"X":1,"Y":2}`)},
		{name: "named float64", value: ratio(1.5), want: attribute.Float64("k", 1.5)},
		{name: "fmt.Sprint when JSON fails", value: noJSON{}, want: attribute.String("k", "{<nil>}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, otelzap.Attribute("k", tt.value))
		})
	}
}

func TestLogValue(t *testing.T) {
	tests := []struct {
		name  string
		value any
		want  attribute.Value
	}{
		{name: "nil", value: nil, want: attribute.StringValue("<nil>")},
		{name: "string", value: "v", want: attribute.StringValue("v")},
		{name: "int", value: 42, want: attribute.IntValue(42)},
		{name: "int64", value: int64(42), want: attribute.Int64Value(42)},
		{name: "uint64 in the int64 range", value: uint64(42), want: attribute.Int64Value(42)},
		{name: "uint64 past the int64 range", value: uint64(math.MaxUint64), want: attribute.StringValue("18446744073709551615")},
		{name: "float64", value: 1.5, want: attribute.Float64Value(1.5)},
		{name: "bool", value: true, want: attribute.BoolValue(true)},
		{name: "fmt.Stringer", value: stringer{}, want: attribute.StringValue("stringer")},
		{name: "error", value: errors.New("failed"), want: attribute.StringValue("failed")},
		{
			name:  "slice, element by element",
			value: []any{"a", 1, uint64(math.MaxUint64)},
			want:  attribute.SliceValue(attribute.StringValue("a"), attribute.IntValue(1), attribute.StringValue("18446744073709551615")),
		},
		{name: "array", value: [1]string{"a"}, want: attribute.SliceValue(attribute.StringValue("a"))},
		{name: "slice of named strings", value: []named{"a"}, want: attribute.SliceValue(attribute.StringValue("a"))},
		{name: "named bool", value: flag(true), want: attribute.BoolValue(true)},
		{name: "named int", value: counter(7), want: attribute.Int64Value(7)},
		{name: "named string", value: named("v"), want: attribute.StringValue("v")},
		{name: "struct as JSON", value: point{X: 1, Y: 2}, want: attribute.StringValue(`{"X":1,"Y":2}`)},
		{name: "named float64", value: ratio(1.5), want: attribute.Float64Value(1.5)},
		{name: "fmt.Sprint when JSON fails", value: noJSON{}, want: attribute.StringValue("{<nil>}")},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.want, otelzap.LogValue(tt.value))
		})
	}
}
