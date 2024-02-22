package gofindit

// Supported types
var supportedTypes = []string{"string", "int", "uint", "float", "bool", "time.Time", "[]string", "[]int"}

func isSupportedType(t string) bool {
	for _, supportedType := range supportedTypes {
		if t == supportedType {
			return true
		}
	}

	return false
}
