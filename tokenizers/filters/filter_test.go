package filters

import "testing"

func TestFuncsID(t *testing.T) {
	tests := []struct {
		name   string
		funcs  []Func
		result string
	}{
		{
			name:   "empty",
			funcs:  []Func{},
			result: "",
		},
		{
			name:   "single",
			funcs:  []Func{Lowercase},
			result: "ecf35cbbd166",
		},
		{
			name:   "multiple",
			funcs:  []Func{Lowercase, RemoveStopwords},
			result: "76040804a901",
		},
		{
			name:   "multiple_reverse",
			funcs:  []Func{RemoveStopwords, Lowercase},
			result: "007bdd218ee4", // Should be different from multiple
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := FuncsID(test.funcs...)
			if result != test.result {
				t.Errorf("Expected %v, got %v", test.result, result)
			}
		})
	}
}

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
			result, err := Lowercase(test.tokens)
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
