package govm

import "fmt"

type TypeError struct{ expected, actual Type }

func (e TypeError) Error() string {
	return fmt.Sprintf("Type error: expected %s, got %s", e.expected, e.actual)
}

func (t Type) String() string {
	switch t.Kind {
	case Int:
		return "int"
	case Float:
		return "float"
	case Bool:
		return "bool"
	case String:
		return "string"
	case Struct:
		panic("Structs not implemented")
	default:
		panic("Unknown type")
	}
}

type NameError struct{ name Symbol }

func (e NameError) Error() string {
	return fmt.Sprintf("Name error: could not find variable named %s", e.name)
}
