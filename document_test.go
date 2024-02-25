package gofindit

import "fmt"

func ExampleNewDocument() {
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
	document, err := NewDocument(doc)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(document.FieldTypes)
	fmt.Println(document.FieldValues)

	// Output: map[Age:int Name:string]
	// map[Age:10 Name:Test]
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
	document, _ := NewDocument(doc)

	value, valueType, found := document.GetFieldValue("name")
	fmt.Println(value, valueType, found)

	value, valueType, found = document.GetFieldValue("age")
	fmt.Println(value, valueType, found)

	// Output: Test string true
	// 10 int true
}
