package gofindit

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func ExampleIndex_Search() {
	type Test struct {
		Name string `find:"Name"`
		Age  int    `find:"Age"`
	}

	// Create a new index
	index := New()

	// Create a new document
	doc := Test{
		Name: "Billy is my friend",
		Age:  10,
	}

	// Index the document
	id := "1"
	err := index.Index(id, doc)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "Name",
				Type:  "partial",
				Value: "friend",
			},
		},
	}

	// Search for the document
	results, err := index.Search(search)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v", results)

	// Output: [{Name:Billy is my friend Age:10}]
}

func TestIndex_Search_match(t *testing.T) {
	// Grab a random TestIndex document
	_, doc := TestIndex.Random()
	docData := doc.(TestData)

	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "match",
				Value: docData.Name,
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the name matches the value
	for _, result := range results {
		doc := result.(TestData)
		if doc.Name != docData.Name {
			t.Errorf("expected name to match, got %s", doc.Name)
		}
	}
}

func TestIndex_Search_partial(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "partial",
				Value: "is", // Example Christina
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the name contains the value
	for _, result := range results {
		doc := result.(TestData)
		if !strings.Contains(doc.Name, "is") {
			t.Errorf("expected name to contain is, got %s", doc.Name)
		}
	}
}

func TestIndex_Search_rangeInt(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "age",
				Type:  "range",
				Value: []int{30, 60},
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the age is greater than 30
	for _, result := range results {
		doc := result.(TestData)
		if doc.Age <= 30 && doc.Age >= 60 {
			t.Errorf("expected age to be greater than 30, got %d", doc.Age)
		}
	}
}

func TestIndex_Search_rangeFloat(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "age",
				Type:  "range",
				Value: []float64{30, 60},
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the age is greater than 30
	for _, result := range results {
		doc := result.(TestData)
		if doc.Age <= 30 && doc.Age >= 60 {
			t.Errorf("expected age to be greater than 30, got %d", doc.Age)
		}
	}
}

func TestIndex_Search_rangeTime(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "birthday",
				Type:  "range",
				Value: []time.Time{time.Now().AddDate(-60, 0, 0), time.Now().AddDate(-20, 0, 0)},
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the age is greater than 30
	for _, result := range results {
		doc := result.(TestData)
		if !doc.Birthday.After(time.Now().AddDate(-60, 0, 0)) || !doc.Birthday.Before(time.Now().AddDate(-20, 0, 0)) {
			t.Errorf("expected birthday to be between 60 and 20 years ago, got %s", doc.Birthday)
		}
	}
}

func TestIndex_Search_sort(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Sort:   "asc",
		SortBy: "name",
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "partial",
				Value: "i",
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the name is in order
	for i := 1; i < len(results); i++ {
		if results[i].(TestData).Name < results[i-1].(TestData).Name {
			t.Errorf("expected name to be in order, got %s and %s", results[i-1].(TestData).Name, results[i].(TestData).Name)
		}
	}
}

func TestIndex_Search_sort_desc(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Sort:   "desc",
		SortBy: "name",
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "partial",
				Value: "i",
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}

	// Loop through results and check if the name is in order
	for i := 1; i < len(results); i++ {
		if results[i].(TestData).Name > results[i-1].(TestData).Name {
			t.Errorf("expected name to be in order, got %s and %s", results[i-1].(TestData).Name, results[i].(TestData).Name)
		}
	}
}

func TestIndex_Search_limit(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "partial",
				Value: "is", // Example Christina
			},
		},
		Limit: 5,
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) != 5 {
		t.Errorf("expected results to be 5, got %d", len(results))
		return
	}
}

func TestIndex_Search_offset(t *testing.T) {
	// Create a search query
	search := SearchQuery{
		Limit: 10,
		Fields: []SearchQueryField{
			{
				Field: "name",
				Type:  "partial",
				Value: "i",
			},
		},
	}

	// Search for the document
	results, err := TestIndex.Search(search)
	if err != nil {
		t.Error(err)
		return
	}

	if len(results) == 0 {
		t.Errorf("expected results to be greater than 0, got %d", len(results))
		return
	}
}
