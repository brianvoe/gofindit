package fields

import (
	"reflect"
	"testing"
)

func TestStringToBytes(t *testing.T) {
	tests := []struct {
		name    string
		input   any
		want    []byte
		wantErr bool
	}{
		{"Valid string", "hello", []byte("hello"), false},
		{"Empty string", "", []byte(""), false},
		{"Non-string", 123, nil, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := stringToBytes(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("stringToBytes() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("stringToBytes() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPartial_Process(t *testing.T) {
	p := &Partial{}

	if err := p.Process("test"); err != nil {
		t.Errorf("Process() failed with error: %v", err)
	}

	want := []byte("test")
	if !reflect.DeepEqual(p.value, want) {
		t.Errorf("Process() stored value = %v, want %v", p.value, want)
	}
}

func TestPartial_ToSearchBytes(t *testing.T) {
	p := &Partial{}

	searchBytes, err := p.ToSearchBytes("search")
	if err != nil {
		t.Errorf("ToSearchBytes() failed with error: %v", err)
	}

	want := []byte("search")
	if !reflect.DeepEqual(searchBytes, want) {
		t.Errorf("ToSearchBytes() = %v, want %v", searchBytes, want)
	}
}

func TestPartial_Search(t *testing.T) {
	p := &Partial{}
	p.Process("hello world")

	tests := []struct {
		searchValue string
		want        bool
	}{
		{"hello", true},
		{"world", true},
		{"hello world", true},
		{"goodbye", false},
	}

	for _, tt := range tests {
		searchBytes, _ := p.ToSearchBytes(tt.searchValue)
		got, _ := p.Search(searchBytes)
		if got != tt.want {
			t.Errorf("Search(%s) = %v, want %v", tt.searchValue, got, tt.want)
		}
	}
}
