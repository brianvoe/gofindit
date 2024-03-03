package gofindit

import "testing"

func TestLowercaseFilter(t *testing.T) {
	tests := []struct {
		name   string
		tokens []string
		result []string
	}{
		{
			name:   "empty",
			tokens: []string{},
			result: []string{},
		},
		{
			name:   "single",
			tokens: []string{"Hello"},
			result: []string{"hello"},
		},
		{
			name:   "multiple",
			tokens: []string{"Hello", "World"},
			result: []string{"hello", "world"},
		},
		{
			name:   "mixed case",
			tokens: []string{"Hello", "WORLD"},
			result: []string{"hello", "world"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := FilterLowercase(test.tokens)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}
			if len(result) != len(test.result) {
				t.Errorf("Expected %v, got %v", test.result, result)
			}
			for i, r := range result {
				if r != test.result[i] {
					t.Errorf("Expected %v, got %v", test.result, result)
				}
			}
		})
	}
}
