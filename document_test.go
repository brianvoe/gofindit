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

	fmt.Println(document.Fields)

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

	nameField, _ := document.GetField("name")
	fmt.Println(nameField.Value())

	valueField, _ := document.GetField("age")
	fmt.Println(valueField.Value())

	// Output: test true
	// 10 true
}
