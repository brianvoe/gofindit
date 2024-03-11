package gofindit

import (
	"fmt"
)

func Example() {
	type Test struct {
		Name string `find:"Name" index:"leve_string"`
		Age  int    `find:"Age"`
	}

	// Create a new index
	index := New()

	// Create a new document
	doc := Test{
		Name: "Test",
		Age:  10,
	}

	// Index the document
	id := "1"
	err := index.Index(id, doc)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Get the document
	docGet, err := index.Get(id)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Printf("%+v", docGet)

	// Output: {Name:Test Age:10}
}
