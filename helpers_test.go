package gofindit

import (
	"fmt"
	"math"
	"reflect"
	"testing"
)

func TestStringToAny(t *testing.T) {
	examples := []struct {
		name       string
		value      string
		result     any
		resultType string // string, int, float, bool, []int, []float64, []bool, []string
		err        bool
	}{
		{
			name:       "int",
			value:      "1",
			result:     1,
			resultType: "int",
			err:        false,
		},
		{
			name:       "float",
			value:      "1.1",
			result:     1.1,
			resultType: "float",
			err:        false,
		},
		{
			name:       "bool",
			value:      "true",
			result:     true,
			resultType: "bool",
			err:        false,
		},
		{
			name:       "int array",
			value:      "1,2,3",
			result:     []int{1, 2, 3},
			resultType: "[]int",
			err:        false,
		},
		{
			name:       "float array",
			value:      "1.1,2.2,3.3",
			result:     []float64{1.1, 2.2, 3.3},
			resultType: "[]float64",
			err:        false,
		},
		{
			name:       "bool array",
			value:      "true,false,true",
			result:     []bool{true, false, true},
			resultType: "[]bool",
			err:        false,
		},
		{
			name:       "string array",
			value:      "one,two,three",
			result:     []string{"one", "two", "three"},
			resultType: "[]string",
			err:        false,
		},
		// {
		// 	name:       "any array",
		// 	value:      "1,2.2,true,one",
		// 	result:     []any{1, 2.2, true, "one"},
		// 	resultType: "[]any",
		// 	err:        false,
		// },
	}

	for _, example := range examples {
		t.Run(example.name, func(t *testing.T) {
			got, err := stringToAny(example.value)
			if (err != nil) != example.err {
				t.Errorf("stringToAny(%q) error = %v, wantErr %v", example.value, err, example.err)
				return
			}

			if !reflect.DeepEqual(got, example.result) {
				// output result and type
				fmt.Printf("result: %v, type: %T\n", got, got)
				fmt.Printf("expected result: %v, type: %s\n\n", example.result, example.resultType)

				t.Errorf("stringToAny(%q) = %v, want %v", example.value, got, example.result)
			}
		})
	}
}

func TestToFloat64(t *testing.T) {
	tests := []struct {
		name   string
		value  any
		result float64
		err    bool
	}{
		{
			name:   "int",
			value:  1,
			result: 1,
			err:    false,
		},
		{
			name:   "int8",
			value:  int8(1),
			result: 1,
			err:    false,
		},
		{
			name:   "int16",
			value:  int16(1),
			result: 1,
			err:    false,
		},
		{
			name:   "int32",
			value:  int32(1),
			result: 1,
			err:    false,
		},
		{
			name:   "int64",
			value:  int64(1),
			result: 1,
			err:    false,
		},
		{
			name:   "uint",
			value:  uint(1),
			result: 1,
			err:    false,
		},
		{
			name:   "uint8",
			value:  uint8(1),
			result: 1,
			err:    false,
		},
		{
			name:   "uint16",
			value:  uint16(1),
			result: 1,
			err:    false,
		},
		{
			name:   "uint32",
			value:  uint32(1),
			result: 1,
			err:    false,
		},
		{
			name:   "uint64",
			value:  uint64(1),
			result: 1,
			err:    false,
		},
		{
			name:   "float32",
			value:  float32(1.1),
			result: 1.1,
			err:    false,
		},
		{
			name:   "float64",
			value:  float64(1.1),
			result: 1.1,
			err:    false,
		},
		{
			name:   "string",
			value:  "1",
			result: 0,
			err:    true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := toFloat64(test.value)
			if test.err && err == nil {
				t.Errorf("expected error, got nil")
			}
			if !test.err && err != nil {
				t.Errorf("expected nil, got %s", err)
			}
			if math.Round(result*10000)/10000 != math.Round(test.result*10000)/10000 {
				t.Errorf("expected %f, got %f", test.result, result)
			}
		})
	}
}

func TestIsZero(t *testing.T) {
	tests := []struct {
		name   string
		value  interface{}
		result bool
	}{
		{
			name:   "int",
			value:  1,
			result: false,
		},
		{
			name:   "int8",
			value:  int8(1),
			result: false,
		},
		{
			name:   "int16",
			value:  int16(1),
			result: false,
		},
		{
			name:   "int32",
			value:  int32(1),
			result: false,
		},
		{
			name:   "int64",
			value:  int64(1),
			result: false,
		},
		{
			name:   "uint",
			value:  uint(1),
			result: false,
		},
		{
			name:   "uint8",
			value:  uint8(1),
			result: false,
		},
		{
			name:   "uint16",
			value:  uint16(1),
			result: false,
		},
		{
			name:   "uint32",
			value:  uint32(1),
			result: false,
		},
		{
			name:   "uint64",
			value:  uint64(1),
			result: false,
		},
		{
			name:   "float32",
			value:  float32(1.1),
			result: false,
		},
		{
			name:   "float64",
			value:  float64(1.1),
			result: false,
		},
		{
			name:   "string",
			value:  "1",
			result: false,
		},
		{
			name:   "empty string",
			value:  "",
			result: true,
		},
		{
			name:   "bool true",
			value:  true,
			result: false,
		},
		{
			name:   "bool false",
			value:  false,
			result: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := isZero(test.value)
			if result != test.result {
				t.Errorf("expected %t, got %t", test.result, result)
			}
		})
	}
}
