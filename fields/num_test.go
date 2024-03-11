package fields

import "testing"

func TestNumField_Process(t *testing.T) {
	tests := []struct {
		name      string
		input     any
		expectErr bool
	}{
		{"Int", 42, false},
		{"Int8", int8(42), false},
		{"Int16", int16(42), false},
		{"Int32", int32(42), false},
		{"Int64", int64(42), false},
		{"Uint", uint(42), false},
		{"Uint8", uint8(42), false},
		{"Uint16", uint16(42), false},
		{"Uint32", uint32(42), false},
		{"Uint64", uint64(42), false},
		{"Float32", float32(42.42), false},
		{"Float64", float64(42.42), false},
		{"String", "not a number", true}, // Invalid type
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			nf := &NumField{}
			err := nf.Process(tt.input)

			if (err != nil) != tt.expectErr {
				t.Errorf("Process() error = %v, expectErr %v", err, tt.expectErr)
			}
		})
	}
}

func TestNumField_Search(t *testing.T) {
	// Process a value first
	nf := &NumField{}
	inputValue := float64(42.42)
	if err := nf.Process(inputValue); err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	searchValue, _ := numToSearchBytes(float64(42.42))    // Matching value
	nonMatchingValue, _ := numToSearchBytes(float64(100)) // Non-matching value

	tests := []struct {
		name      string
		searchVal []byte
		wantMatch bool
	}{
		{"Matching value", searchValue, true},
		{"Non-matching value", nonMatchingValue, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMatch, err := nf.Search(tt.searchVal)
			if err != nil {
				t.Errorf("Search() error = %v", err)
				return
			}

			if gotMatch != tt.wantMatch {
				t.Errorf("Search() gotMatch = %v, want %v", gotMatch, tt.wantMatch)
			}
		})
	}
}

func TestNumField_SearchRange(t *testing.T) {
	// Process a middle value for range testing
	nf := &NumField{}
	if err := nf.Process(50); err != nil {
		t.Fatalf("Process() error = %v", err)
	}

	minBytes, _ := numToSearchBytes(40)
	maxBytes, _ := numToSearchBytes(60)
	outOfRangeLow, _ := numToSearchBytes(30)
	outOfRangeHigh, _ := numToSearchBytes(70)

	tests := []struct {
		name          string
		min           []byte
		max           []byte
		expectInRange bool
	}{
		{"Value in range", minBytes, maxBytes, true},
		{"Value below range", outOfRangeLow, minBytes, false},
		{"Value above range", maxBytes, outOfRangeHigh, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			inRange, err := nf.SearchRange(tt.min, tt.max)
			if err != nil {
				t.Errorf("SearchRange() error = %v", err)
				return
			}

			if inRange != tt.expectInRange {
				t.Errorf("SearchRange() inRange = %v, expectInRange %v", inRange, tt.expectInRange)
			}
		})
	}
}
