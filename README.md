![Gofindit](https://raw.githubusercontent.com/brianvoe/gofindit/master/logo.png)

# Gofindit [![Go Report Card](https://goreportcard.com/badge/github.com/brianvoe/gofindit)](https://goreportcard.com/report/github.com/brianvoe/gofindit) ![Test](https://github.com/brianvoe/gofindit/workflows/Test/badge.svg?branch=master) [![GoDoc](https://godoc.org/github.com/brianvoe/gofindit?status.svg)](https://godoc.org/github.com/brianvoe/gofindit) [![license](http://img.shields.io/badge/license-MIT-green.svg?style=flat)](https://raw.githubusercontent.com/brianvoe/gofindit/master/LICENSE.txt)

Struct index and searching

[![ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/G2G0R5EJT)

<a href="https://www.buymeacoffee.com/brianvoe" target="_blank"><img src="https://www.buymeacoffee.com/assets/img/custom_images/orange_img.png" alt="Buy Me A Coffee" style="height: auto !important;width: auto !important;" ></a>

## Features

- Simplicity

## Installation

```go
go get github.com/brianvoe/gofindit
```

## Simple Usage

```go
import "github.com/brianvoe/gofindit"

type Test struct {
    Name string `find:"name"` // Tag with find
    Age  int                  // or field name Age is used
}

// Create a new index
index := gofindit.New()

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
```

## Search Usage

```go
// Create a new index
index := gofindit.New()

// Add all your documents
for i := 0; i < 1000; i++ {
    // Add to index
    // ... Code here
}

// Create a search query
search := SearchQuery{
    Limit:  10, // default 10
	Skip:   0,  // default 0
	Sort:   "", // "", asc or desc
	SortBy: "", // field name

    // Search fields
    Fields: []SearchQueryField{
        {
            Field: "name",    // find tag
            Type:  "partial", // match, partial or range
            Value: "billy",   // Case insensitive
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

// Output: [{Name:Billy Age:10}]
```