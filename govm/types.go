package govm

type Value interface{}

type Symbol string

type Kind byte

const (
	Int Kind = iota
	Float
	Bool
	String
	FuncT
	Struct
)

type Type struct {
	Kind Kind
	Sig TypeSignature // Only used when Kind == Func
	I int // Index in struct table. Only used when Kind == Struct
}

var (
	TypeInt Type = Type{Int, TypeSignature{}, 0}
	TypeFloat Type = Type{Float, TypeSignature{}, 0}
	TypeBool Type = Type{Bool, TypeSignature{}, 0}
	TypeString Type = Type{String, TypeSignature{}, 0}
)

func TypeOf(v Value) (t Type) {
	switch v := v.(type) {
	case int:
		return TypeInt
	case float64:
		return TypeFloat
	case bool:
		return TypeBool
	case string:
		return TypeString
	case Function:
		t.Kind = FuncT
		t.Sig = v.Sig
		return
	case Builtin:
		t.Kind = FuncT
		t.Sig = v.Sig
		return
	default:
		panic("Unknown type")
	}
}

func (t Type) TypeCheck(v Value) bool {
	switch TypeOf(v).Kind {
	case Struct:
		panic("Structs not yet implemented")
	case t.Kind:
		return true
	default:
		return false
	}
}

type TypeSignature struct {
	Args, Ret []Type
}

type Function struct {
	Sig TypeSignature
	Code []byte
}

type Builtin struct {
	Sig TypeSignature
	F func(...Value) []Value
}

type Stack []Value

func (s *Stack) Push(v Value) {
	*s = append(*s, v)
}

func (s *Stack) Pop() Value {
	l := len(*s) - 1
	v := (*s)[l]
	*s = (*s)[:l]
	return v
}

func (s *Stack) PopN(n int) []Value {
	vals := (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return vals
}

func (s *Stack) Peek(n int) Value {
	l := len(*s) - 1
	return (*s)[l-n]
}

func (s *Stack) Dup() {
	*s = append(*s, (*s)[len(*s)-1])
}

func (s *Stack) Swap() {
	l := len(*s) - 1
	(*s)[l-1], (*s)[l] = (*s)[l], (*s)[l-1]
}

type Scope struct {
	parent *Scope
	m map[Symbol]Value
}

func (s *Scope) Child() *Scope {
	return &Scope{s, nil}
}

func (s *Scope) Set(k Symbol, v Value) {
	if s.m == nil {
		s.m = make(map[Symbol]Value)
	}
	s.m[k] = v
}

func (s *Scope) Get(k Symbol) Value {
	if s.m != nil {
		if v, ok := s.m[k]; ok {
			return v
		}
	}
	if s.parent == nil {
		// TODO: handle panics
		panic(NameError{k})
	}
	return s.parent.Get(k)
}
