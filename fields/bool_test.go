package fields

import (
	"reflect"
	"testing"
)

func TestBool_Process(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{"true to bytes", true, []byte{1}, false},
		{"false to bytes", false, []byte{0}, false},
		{"invalid type", "not a bool", nil, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &Bool{}
			err := b.Process(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(b.value, tt.want) {
				t.Errorf("Process() got = %v, want %v", b.value, tt.want)
			}
		})
	}
}

func TestBool_Search(t *testing.T) {
	b := &Bool{}
	_ = b.Process(true) // Assuming Process works as expected here

	tests := []struct {
		name    string
		search  []byte
		want    bool
		wantErr bool
	}{
		{"search true", []byte{1}, true, false},
		{"search false", []byte{0}, false, false},
		{"invalid length", []byte{1, 1}, false, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := b.Search(tt.search)
			if (err != nil) != tt.wantErr {
				t.Errorf("Search() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Search() got = %v, want %v", got, tt.want)
			}
		})
	}
}
