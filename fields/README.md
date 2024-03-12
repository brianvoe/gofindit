# Fields

Welcome to the `fields` package—your toolkit for dynamic field definition in indexing systems. It simplifies data management and enhances search capabilities with ease.

## Features

- Manage dynamic data types
- Execute exact and partial searches
- Extend with custom field types
- [Customize through configurations](#custom-fields)

## Types

- Text - string
- Number - all int, uint and floats
- Boolean - bool
- Date - time.Time

## Fields

The following fields are currently registered and available for use:

- Text (`text`) - Default, exact match
- Partial (`partial`) - Partial match
- Num (`num`) - All number types, exact match and range search
- Bool (`bool`) - Exact match
- Date (`date`) - Exact match and range search

## Usage

```go
import "github.com/brianvoe/gofindit/fields"

func main() {
    // Get a field you want from the store
    textField, _ := fields.GetField("text", nil)

    // Process a string value
    _ = textField.Process("Hello")

    // Use textField to run ToSearchBytes so we can use the field 
    // structure to generate the bytes required for the search
    searchBytes := textField.ToSearchBytes("Hello")

    // Search for a matching value
    match, _ := textField.Search(serchBytes)

    fmt.Printf("Search match: %v\n", match)
}
```

## Custom Fields

To create a custom field, you need to implement the `Field` interface.
See test for examples

The `Field` interface has the following methods:

```go
import "github.com/brianvoe/gofindit/fields"

type Field interface {
    // Their are 4 types of fields
    // fields.TextType, fields.NumberType, fields.BooleanType, fields.DateType
    Type() string
    Value() []byte

    // Process will take in an any value and
    // use it to fill out the struct fields
    Process(val any) error

    // To use the struct values to calculate how search
    // bytes should be passed to the search function
    ToSearchBytes(val any) ([]byte, error)

    Search(val []byte) (bool, error)
    SearchRange(min, max []byte) (bool, error)
}
```
