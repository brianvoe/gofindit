package gofindit

import "fmt"

type Test struct {
	Name string `find:"Name"`
	Age  int    `find:"Age"`
}

func Example() {
	// Create a new index
	index, err := New("test", Test{})
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a new document
	doc := Test{
		Name: "Test",
		Age:  10,
	}

	// Index the document
	id := "1"
	err = index.Index(id, doc)
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

	fmt.Println(docGet)

	// Search the index
	results, err := index.Search(SearchQuery{
		Limit: 10,
		Skip:  0,
		Sort:  "asc",
		Fields: []DocumentQuery{
			{
				Name:  "Name",
				Value: "Test",
			},
		},
	})
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(results)

	// Output: [1]
}
