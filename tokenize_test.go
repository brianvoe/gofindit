package gofindit

import (
	"reflect"
	"testing"
)

func TestTokenize(t *testing.T) {
	tests := []struct {
		name    string
		str     string
		filters []FilterFunc
		want    []string
		wantErr bool
	}{
		{
			name:    "empty",
			str:     "",
			filters: []FilterFunc{},
			want:    []string{},
			wantErr: false,
		},
		{
			name:    "single",
			str:     "Hello, World!",
			filters: []FilterFunc{FilterLowercase},
			want:    []string{"hello", "world"},
			wantErr: false,
		},
		{
			name:    "multiple",
			str:     "Hello, World and welcome!",
			filters: []FilterFunc{FilterLowercase, FilterOutStopwords},
			want:    []string{"hello", "world", "welcome"},
			wantErr: false,
		},
		{
			name: "apostrophe",
			str:  "I'm a string with an apostrophe",
			filters: []FilterFunc{
				FilterLowercase,
				FilterOutStopwords,
			},
			want:    []string{"string", "apostrophe"},
			wantErr: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := Tokenize(test.str, test.filters)
			if (err != nil) != test.wantErr {
				t.Errorf("Tokenize() error = %v, wantErr %v", err, test.wantErr)
				return
			}
			if !reflect.DeepEqual(got, test.want) {
				t.Errorf("Tokenize() = %v, want %v", got, test.want)
			}
		})
	}
}
