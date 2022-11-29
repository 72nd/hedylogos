package enum

import "fmt"

// Used to simplify the usage of the enums in the cli and notification
// portion of the library.
type Enum[T ~string] interface {
	// A list of all possible values.
	Instances() []T
}

// Returns an Enum instance by it's value.
func EnumByValue[T ~string](e Enum[T], value T) (*T, error) {
	for _, entry := range e.Instances() {
		if entry == value {
			return &entry, nil
		}
	}
	return nil, fmt.Errorf("no %T enum found for value %s", e, value)
}

// Returns all values of an Enum type.
func EnumValues[T ~string](e Enum[T]) []string {
	var rsl []string
	for _, item := range e.Instances() {
		rsl = append(rsl, string(item))
	}
	return rsl
}
