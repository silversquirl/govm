package types

import "fmt"

type TypeError struct{ Expected, Actual Type }

func (e TypeError) Error() string {
	return fmt.Sprintf("Type error: expected %s, got %s", e.Expected, e.Actual)
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

type NameError struct{ Name Symbol }

func (e NameError) Error() string {
	return fmt.Sprintf("Name error: could not find variable named %s", e.Name)
}

type StackUnderflow struct{}

func (e StackUnderflow) Error() string {
	return "Stack underflow"
}
