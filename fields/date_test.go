package fields

import (
	"bytes"
	"encoding/binary"
	"testing"
	"time"
)

func TestDate_Process(t *testing.T) {
	location, _ := time.LoadLocation("UTC") // Ensure consistent use of the UTC time zone

	tests := []struct {
		name        string
		date        time.Time
		granularity string
		wantYear    int
		wantMonth   time.Month
		wantDay     int
	}{
		{
			name:        "Test Day Granularity",
			date:        time.Date(2023, 3, 14, 12, 0, 0, 0, location),
			granularity: "day",
			wantYear:    2023,
			wantMonth:   3,
			wantDay:     14,
		},
		{
			name:        "Test Month Granularity",
			date:        time.Date(2023, 3, 14, 12, 0, 0, 0, location),
			granularity: "month",
			wantYear:    2023,
			wantMonth:   3,
			wantDay:     1,
		},
		{
			name:        "Test Year Granularity",
			date:        time.Date(2023, 3, 14, 12, 0, 0, 0, location),
			granularity: "year",
			wantYear:    2023,
			wantMonth:   1,
			wantDay:     1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			field, _ := NewDate(map[string]any{"granularity": tc.granularity})

			fieldDate := field.(*Date)
			fieldDate.Process(tc.date)

			// Convert the stored bytes back to a time.Time value, explicitly using UTC
			var timestamp int64
			buf := bytes.NewReader(fieldDate.value)
			binary.Read(buf, binary.BigEndian, &timestamp)
			storedDate := time.Unix(timestamp, 0).In(location) // Use the same location for comparison

			if storedDate.Year() != tc.wantYear || storedDate.Month() != tc.wantMonth || storedDate.Day() != tc.wantDay {
				t.Errorf("Process() with %v granularity, got %v, want %v-%v-%v",
					tc.granularity, storedDate, tc.wantYear, tc.wantMonth, tc.wantDay)
			}
		})
	}
}

func TestDate_Search(t *testing.T) {
	// Initial setup
	testDate := time.Date(2023, 3, 14, 0, 0, 0, 0, time.UTC)
	field, _ := NewDate(map[string]any{"granularity": "day"})
	field.Process(testDate)

	// Prepare a matching search value
	searchBytes, _ := dateToSearchBytes(testDate, "day")
	match, err := field.Search(searchBytes)
	if err != nil {
		t.Fatalf("Search() error = %v", err)
	}
	if !match {
		t.Errorf("Expected search to match, but it did not")
	}

	// Prepare a non-matching search value (different day)
	nonMatchDate := time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC)
	nonMatchBytes, _ := dateToSearchBytes(nonMatchDate, "day")
	nonMatch, err := field.Search(nonMatchBytes)
	if err != nil {
		t.Fatalf("Search() error with non-matching date = %v", err)
	}
	if nonMatch {
		t.Errorf("Expected search with non-matching date to not match, but it did")
	}
}

func TestDate_SearchRange(t *testing.T) {
	location, _ := time.LoadLocation("UTC") // Use UTC for consistent testing

	// Create a Date instance with day granularity for testing
	config := map[string]any{"granularity": "day"}
	date, err := NewDate(config)
	if err != nil {
		t.Fatalf("NewDate() failed with error: %v", err)
	}

	// Process a specific date to be searched
	testDate := time.Date(2023, 3, 15, 12, 0, 0, 0, location)
	if err := date.Process(testDate); err != nil {
		t.Fatalf("Process() failed with error: %v", err)
	}

	// Define various ranges to test
	tests := []struct {
		name          string
		start         time.Time
		end           time.Time
		expectInRange bool
	}{
		{
			name:          "Date exactly at start of range",
			start:         testDate,
			end:           testDate.Add(24 * time.Hour), // One day later
			expectInRange: true,
		},
		{
			name:          "Date exactly at end of range",
			start:         testDate.Add(-24 * time.Hour), // One day earlier
			end:           testDate,
			expectInRange: true,
		},
		{
			name:          "Date within range",
			start:         testDate.Add(-24 * time.Hour),
			end:           testDate.Add(24 * time.Hour),
			expectInRange: true,
		},
		{
			name:          "Date before range",
			start:         testDate.Add(24 * time.Hour),
			end:           testDate.Add(48 * time.Hour),
			expectInRange: false,
		},
		{
			name:          "Date after range",
			start:         testDate.Add(-48 * time.Hour),
			end:           testDate.Add(-24 * time.Hour),
			expectInRange: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Convert start and end to bytes
			startBytes, err := dateToSearchBytes(tc.start, config["granularity"].(string))
			if err != nil {
				t.Fatalf("Failed to convert start date to bytes: %v", err)
			}
			endBytes, err := dateToSearchBytes(tc.end, config["granularity"].(string))
			if err != nil {
				t.Fatalf("Failed to convert end date to bytes: %v", err)
			}

			// Perform the range search
			inRange, err := date.SearchRange(startBytes, endBytes)
			if err != nil {
				t.Fatalf("SearchRange() failed with error: %v", err)
			}
			if inRange != tc.expectInRange {
				t.Errorf("Expected inRange to be %v, got %v", tc.expectInRange, inRange)
			}
		})
	}
}
