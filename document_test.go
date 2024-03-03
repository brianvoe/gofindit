package gofindit

import "fmt"

func ExampleNewDoc() {
	type Test struct {
		Name string `find:"Name"`
		Age  int    `find:"Age"`
	}

	// Create a new document
	doc := Test{
		Name: "Test",
		Age:  10,
	}

	// Create a new document
	document, err := NewDoc(doc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(document.Values)

	// Output: map[Age:{int 10 []} Name:{string Test [test]}]
}

func ExampleDocument_GetFieldValue() {
	type Test struct {
		Name string `find:"name"`
		Age  int    `find:"age"`
	}

	// Create a new document
	doc := Test{
		Name: "Test",
		Age:  10,
	}

	// Create a new document
	document, _ := NewDoc(doc)

	value, _ := document.GetFieldValue("name")
	fmt.Println(value)

	value, _ = document.GetFieldValue("age")
	fmt.Println(value)

	// Output: test true
	// 10 true
}
