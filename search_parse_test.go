package gofindit

import (
	"reflect"
	"testing"
	"time"
)

type TestUser struct {
	ID        string      `find:"id"`
	FirstName string      `find:"first_name"`
	LastName  string      `find:"last_name"`
	Username  string      `find:"username"`
	Password  string      `find:"password"`
	Age       int         `find:"age"`
	Gender    string      `find:"gender"`
	SSN       string      `find:"ssn"`
	Hobby     string      `find:"hobby"`
	Job       string      `find:"job"`
	Address   TestAddress `find:"address"`
	CreatedAt time.Time   `find:"created_at"`
}

type TestAddress struct {
	Street    string  `find:"street"`
	City      string  `find:"city"`
	State     string  `find:"state"`
	Zip       string  `find:"zip"`
	Latitude  float64 `find:"latitude"`
	Longitude float64 `find:"longitude"`
}

func TestMapStringToSearchQuery(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *SearchQuery
		hasError bool
	}{
		// Match
		{
			name:  "match",
			input: "first_name=value",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "match", Value: "value"},
			}},
			hasError: false,
		},
		{
			name:  "match with type set",
			input: "first_name=match:value",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "match", Value: "value"},
			}},
			hasError: false,
		},
		// Partial
		{
			name:  "partial",
			input: "first_name=partial:value",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "partial", Value: "value"},
			}},
			hasError: false,
		},
		{
			name:  "partial multiple",
			input: "first_name=partial:billy&last_name=partial:mister",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "partial", Value: "billy"},
				{Field: "last_name", Type: "partial", Value: "mister"},
			}},
			hasError: false,
		},
		// Range
		{
			name:  "range",
			input: "age=18,38",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "age", Type: "range", Value: []int{18, 38}},
			}},
			hasError: false,
		},
		{
			name:  "range with floats",
			input: "address.latitude=689.555,548.23",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "address.latitude", Type: "range", Value: []float64{689.555, 548.23}},
			}},
		},
		{
			name:  "range with one int, one float",
			input: "age=18,22.5",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "age", Type: "range", Value: []float64{18, 22.5}},
			}},
		},
		// Special cases
		{
			name:  "string using numbers only",
			input: "ssn=123456789", // ssn is a string value
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "ssn", Type: "match", Value: "123456789"},
			}},
		},
		// limit, offset, sort
		{
			name:  "limit",
			input: "limit=10",
			expected: &SearchQuery{
				Limit: 10,
			},
		},
		{
			name:  "skip",
			input: "skip=10",
			expected: &SearchQuery{
				Skip: 10,
			},
		},
		{
			name:  "sort",
			input: "sort=first_name",
			expected: &SearchQuery{
				Sort: "first_name",
			},
		},
		// Mix and match
		{
			name:  "mix and match match and partial",
			input: "first_name=partial:bri&last_name=chatly&age=range:18,38",
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "partial", Value: "bri"},
				{Field: "last_name", Type: "match", Value: "chatly"},
				{Field: "age", Type: "range", Value: []int{18, 38}},
			}},
		},

		// Invalids
		{
			input:    "invalid_input",
			expected: nil,
			hasError: true,
		},
		{
			input:    "firstname=billy&limit=invalid",
			expected: nil,
			hasError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searchQuery, err := StringToSearchQuery(tc.input)
			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(searchQuery, tc.expected) {
					t.Errorf("Expected: \n%+v\n got: \n%+v", tc.expected, searchQuery)
				}
			}
		})
	}
}

func TestMapJsontoSearchQueries(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected *SearchQuery
		hasError bool
	}{
		{
			name:  "single query",
			input: `{"fields":[{"field":"first_name","type":"match","value":"billy"}]}`,
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "match", Value: "billy"},
			}},
			hasError: false,
		},
		{
			name:  "multiple queries",
			input: `{"fields":[{"field":"first_name","type":"match","value":"billy"},{"field":"last_name","type":"match","value":"mister"}]}`,
			expected: &SearchQuery{Fields: []SearchQueryField{
				{Field: "first_name", Type: "match", Value: "billy"},
				{Field: "last_name", Type: "match", Value: "mister"},
			}},
			hasError: false,
		},
		{
			name:  "limit",
			input: `{"limit":10}`,
			expected: &SearchQuery{
				Limit: 10,
			},
			hasError: false,
		},
		{
			name:  "skip",
			input: `{"skip":10}`,
			expected: &SearchQuery{
				Skip: 10,
			},
			hasError: false,
		},
		{
			name:  "sort",
			input: `{"sort":"first_name"}`,
			expected: &SearchQuery{
				Sort: "first_name",
			},
			hasError: false,
		},
		{
			name:  "mix and match",
			input: `{"fields":[{"field":"first_name","type":"match","value":"billy"},{"field":"last_name","type":"match","value":"mister"}],"limit":10,"skip":10,"sort":"first_name"}`,
			expected: &SearchQuery{
				Fields: []SearchQueryField{
					{Field: "first_name", Type: "match", Value: "billy"},
					{Field: "last_name", Type: "match", Value: "mister"},
				},
				Limit: 10,
				Skip:  10,
				Sort:  "first_name",
			},
			hasError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			searchQueries, err := JsontoSearchQueries([]byte(tc.input))
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if tc.hasError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if !reflect.DeepEqual(searchQueries, tc.expected) {
					t.Errorf("Expected: \n%+v\n got: \n%+v", tc.expected, searchQueries)
				}
			}
		})
	}

}
