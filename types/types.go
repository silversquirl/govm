package types

type Value interface{}

type Symbol string

type Kind byte

const (
	Int Kind = 1 << iota
	Float
	Bool
	String
	FuncT
	Struct
)

type Type struct {
	Kind Kind
	Sig  TypeSignature // Only used when Kind is Func
	I    int           // Index in struct table. Only used when Kind is Struct
}

var (
	TypeInt    Type = Type{Int, TypeSignature{}, 0}
	TypeFloat  Type = Type{Float, TypeSignature{}, 0}
	TypeNum    Type = Type{Int | Float, TypeSignature{}, 0}
	TypeBool   Type = Type{Bool, TypeSignature{}, 0}
	TypeString Type = Type{String, TypeSignature{}, 0}
	TypeFunc   Type = Type{FuncT, TypeSignature{}, 0}
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

func (t Type) TypeCheck(val Value) error {
	t2 := TypeOf(val)
	if t.Equal(t2) {
		return nil
	} else {
		return TypeError{t, TypeOf(val)}
	}
}

func (t Type) Equal(t2 Type) bool {
	if t2.Kind == Struct {
		panic("Structs not yet implemented")
	} else if t.Kind&t2.Kind != 0 {
		if t2.Kind == FuncT {
			if t.Sig.Equal(t2.Sig) {
				return true
			}
		} else {
			return true
		}
	}
	return false
}

type TypeSignature struct {
	Args, Ret []Type
}

func (ts TypeSignature) Equal(ts2 TypeSignature) bool {
	if len(ts.Args) != len(ts2.Args) || len(ts.Ret) != len(ts2.Ret) {
		return false
	}
	for i := range ts.Args {
		if ts.Args[i].Equal(ts2.Args[i]) {
			return false
		}
	}
	for i := range ts.Ret {
		if ts.Ret[i].Equal(ts2.Ret[i]) {
			return false
		}
	}
	return true
}

type Function struct {
	Sig  TypeSignature
	Code []byte
}

type Builtin struct {
	Sig TypeSignature
	F   func(...Value) []Value
}

type Stack []Value

func (s *Stack) Push(v Value) {
	*s = append(*s, v)
}

func (s *Stack) Pop() (Value, error) {
	if len(*s) < 1 {
		return nil, StackUnderflow{}
	}
	l := len(*s) - 1
	v := (*s)[l]
	*s = (*s)[:l]
	return v, nil
}

func (s *Stack) PopN(n int) ([]Value, error) {
	if len(*s) < n {
		return nil, StackUnderflow{}
	}
	vals := (*s)[len(*s)-n:]
	*s = (*s)[:len(*s)-n]
	return vals, nil
}

func (s *Stack) Peek(n int) (Value, error) {
	if len(*s) < n + 1 {
		return nil, StackUnderflow{}
	}
	l := len(*s) - 1
	return (*s)[l-n], nil
}

func (s *Stack) Dup() error {
	if len(*s) < 1 {
		return StackUnderflow{}
	}
	*s = append(*s, (*s)[len(*s)-1])
	return nil
}

func (s *Stack) Swap() error {
	if len(*s) < 2 {
		return StackUnderflow{}
	}
	l := len(*s) - 1
	(*s)[l-1], (*s)[l] = (*s)[l], (*s)[l-1]
	return nil
}

type Scope struct {
	Parent *Scope
	m      map[Symbol]Value
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

func (s *Scope) Get(k Symbol) (Value, error) {
	if s.m != nil {
		if v, ok := s.m[k]; ok {
			return v, nil
		}
	}
	if s.Parent == nil {
		return nil, NameError{k}
	}
	return s.Parent.Get(k)
}
