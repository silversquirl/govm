package govm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type VM struct {
	stack Stack
	scope *Scope
	code *bytes.Reader // The bit of code being currently executed
}

func NewVM() (v VM) {
	v.scope = &Scope{}
	return
}

func (v *VM) Load(code []byte) error {
	v.code = bytes.NewReader(code)
	if err := v.exec(); err != nil {
		return err
	}
	v.code = nil
	return nil
}

func (v *VM) exec() error {
	for {
		op, err := v.code.ReadByte()
		if err == io.EOF {
			return nil
		} else if err != nil {
			return err
		}

		switch Instruction(op) {
		case J:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			v.Jump(off)

		case JT:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			v.JumpTrue(off)

		case JF:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			v.JumpFalse(off)

		case JZ:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			v.JumpZero(off)

		case JNz:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			v.JumpNonzero(off)

		case Push:
			val, err := v.readTypedValue()
			if err != nil {
				return nil
			}
			v.Push(val)

		case Pop:
			v.Pop()

		case Dup:
			v.Dup()

		case Swp:
			v.Swap()

		case Set:
			s, err := v.readS()
			if err != nil {
				return err
			}
			v.Set(Symbol(s))

		case Get:
			s, err := v.readS()
			if err != nil {
				return err
			}
			v.Get(Symbol(s))

		case Inc:
			v.Inc()
		case Dec:
			v.Dec()
		case Add:
			v.Add()
		case Sub:
			v.Sub()
		case Mul:
			v.Mul()
		case Div:
			v.Div()

		case EQ:
			v.EQ()
		case NE:
			v.NE()
		case LT:
			v.LT()
		case GT:
			v.GT()
		case LE:
			v.LE()
		case GE:
			v.GE()

		case And:
			v.And()
		case Or:
			v.Or()
		case Xor:
			v.Xor()
		case Not:
			v.Not()

		case BAnd:
			v.BAnd()
		case BOr:
			v.BOr()
		case BXor:
			v.BXor()
		case BNot:
			v.BNot()
		case BLS:
			v.BLS()
		case BRS:
			v.BRS()
		case BSet:
			v.BSet()
		case BClr:
			v.BClr()
		case BTgl:
			v.BTgl()
		case BMtch:
			v.BMtch()

		case Call:
			v.Call()

		case Func:
			sig, err := v.readTS()
			code, err := v.readBytes()
			if err != nil {
				return err
			}
			v.Func(sig, code)

		default:
			panic("Unknown opcode")
		}
	}
}

func (v *VM) readI() (int, error) {
	var i int32
	err := binary.Read(v.code, binary.BigEndian, &i)
	return int(i), err
}

func (v *VM) readF() (f float64, err error) {
	err = binary.Read(v.code, binary.BigEndian, &f)
	return
}

func (v *VM) readBool() (b bool, err error) {
	err = binary.Read(v.code, binary.BigEndian, &b)
	return
}

func (v *VM) readBytes() ([]byte, error) {
	l, err := v.readI()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, l)
	if _, err := io.ReadFull(v.code, buf); err != nil {
		return nil, err
	}
	if len(buf) < l {
		return nil, io.ErrUnexpectedEOF
	}
	return buf, nil
}

func (v *VM) readS() (string, error) {
	s, err := v.readBytes()
	return string(s), err
}

func (v *VM) readType() (Type, error) {
	b, err := v.code.ReadByte()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return Type{}, err
	}

	k := Kind(b)
	switch k {
	case Int, Float, Bool, String:
		return Type{k, TypeSignature{}, 0}, nil
	case Struct:
		if i, err := v.readI(); err != nil {
			return Type{}, err
		} else {
			return Type{k, TypeSignature{}, i}, nil
		}
	default:
		panic("Unknown kind")
	}
}

func (v *VM) readTypedValue() (Value, error) {
	t, err := v.readType()
	if err != nil {
		return nil, err
	}

	switch t.Kind {
	case Int:
		return v.readI()
	case Float:
		return v.readF()
	case Bool:
		return v.readBool()
	case String:
		return v.readS()
	case Struct:
		panic("structs not implemented")
	default:
		panic("Unknown type")
	}
}

func (v *VM) readTS() (TypeSignature, error) {
	var ts TypeSignature
	nargs, err := v.readI()
	if err != nil {
		return TypeSignature{}, err
	}
	ts.Args = make([]Type, nargs)
	for i := 0; i < nargs; i++ {
		if t, err := v.readType(); err != nil {
			return TypeSignature{}, err
		} else {
			ts.Args[i] = t
		}
	}

	nret, err := v.readI()
	if err != nil {
		return TypeSignature{}, err
	}
	ts.Ret = make([]Type, nret)
	for i := 0; i < nret; i++ {
		if t, err := v.readType(); err != nil {
			return TypeSignature{}, err
		} else {
			ts.Ret[i] = t
		}
	}

	return ts, nil
}

func (v *VM) Jump(off int) {
	v.code.Seek(int64(off), io.SeekCurrent)
}

func (v *VM) JumpTrue(off int) {
	// TODO: don't panic
	if v.Pop().(bool) {
		v.Jump(off)
	}
}

func (v *VM) JumpFalse(off int) {
	// TODO: don't panic
	if !v.Pop().(bool) {
		v.Jump(off)
	}
}

func (v *VM) JumpZero(off int) {
	switch n := v.Pop().(type) {
	case int:
		if n == 0 {
			v.Jump(off)
		}
	case float64:
		if n == 0.0 {
			v.Jump(off)
		}
	default:
		// TODO: don't panic
		panic("Type of top of stack not float or int for jz instruction")
	}
}

func (v *VM) JumpNonzero(off int) {
	switch n := v.Pop().(type) {
	case int:
		if n != 0 {
			v.Jump(off)
		}
	case float64:
		if n != 0.0 {
			v.Jump(off)
		}
	default:
		// TODO: don't panic
		panic("Type of top of stack not float or int for jnz instruction")
	}
}

func (v *VM) Push(val Value) {
	v.stack.Push(val)
}

func (v *VM) Pop() Value {
	return v.stack.Pop()
}

func (v *VM) Dup() {
	v.stack.Dup()
}

func (v *VM) Swap() {
	v.stack.Swap()
}

func (v *VM) Set(s Symbol) {
	v.scope.Set(s, v.Pop())
}

func (v *VM) Get(s Symbol) {
	v.Push(v.scope.Get(s))
}

func (v *VM) Inc() {
	// TODO: don't panic
	v.Push(v.Pop().(int)+1)
}

func (v *VM) Dec() {
	// TODO: don't panic
	v.Push(v.Pop().(int)-1)
}

func (v *VM) Add() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a+b)
		case float64:
			v.Push(float64(a)+b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a+float64(b))
		case float64:
			v.Push(a+b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in add instr")
}

func (v *VM) Sub() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a-b)
		case float64:
			v.Push(float64(a)-b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a-float64(b))
		case float64:
			v.Push(a-b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in sub instr")
}

func (v *VM) Mul() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a*b)
		case float64:
			v.Push(float64(a)*b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a*float64(b))
		case float64:
			v.Push(a*b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in mul instr")
}

func (v *VM) Div() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a/b)
		case float64:
			v.Push(float64(a)/b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a/float64(b))
		case float64:
			v.Push(a/b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in div instr")
}

func (v *VM) EQ() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a == b)
		case float64:
			v.Push(float64(a) == b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a == float64(b))
		case float64:
			v.Push(a == b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in eq instr")
}

func (v *VM) NE() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a != b)
		case float64:
			v.Push(float64(a) != b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a != float64(b))
		case float64:
			v.Push(a != b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in ne instr")
}

func (v *VM) LT() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a < b)
		case float64:
			v.Push(float64(a) < b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a < float64(b))
		case float64:
			v.Push(a < b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in lt instr")
}

func (v *VM) GT() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a > b)
		case float64:
			v.Push(float64(a) > b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a > float64(b))
		case float64:
			v.Push(a > b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in gt instr")
}

func (v *VM) LE() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a <= b)
		case float64:
			v.Push(float64(a) <= b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a <= float64(b))
		case float64:
			v.Push(a <= b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in le instr")
}

func (v *VM) GE() {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a >= b)
		case float64:
			v.Push(float64(a) >= b)
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a >= float64(b))
		case float64:
			v.Push(a >= b)
		}
	}
	// TODO: don't panic
	panic("Bad type for a or b in ge instr")
}

func (v *VM) And() {
	b, a := v.Pop(), v.Pop()
	// TODO: don't panic
	v.Push(a.(bool) && b.(bool))
}

func (v *VM) Or() {
	b, a := v.Pop(), v.Pop()
	// TODO: don't panic
	v.Push(a.(bool) || b.(bool))
}

func (v *VM) Xor() {
	b, a := v.Pop(), v.Pop()
	// TODO: don't panic
	v.Push(a.(bool) != b.(bool))
}

func (v *VM) Not() {
	// TODO: don't panic
	v.Push(!v.Pop().(bool))
}

func (v *VM) BAnd() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a&b)
}

func (v *VM) BOr() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a|b)
}

func (v *VM) BXor() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a^b)
}

func (v *VM) BNot() {
	// TODO: don't panic
	a := v.Pop().(int)
	v.Push(^a)
}

func (v *VM) BLS() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a<<uint(b))
}

func (v *VM) BRS() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a>>uint(b))
}

func (v *VM) BSet() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a|(1<<uint(b)))
}

func (v *VM) BClr() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a&^(1<<uint(b)))
}

func (v *VM) BTgl() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a^(1<<uint(b)))
}

func (v *VM) BMtch() {
	// TODO: don't panic
	b, a := v.Pop().(int), v.Pop().(int)
	v.Push(a&b != 0)
}

func (v *VM) checkTypes(types []Type) {
	for i, t := range types {
		if !t.TypeCheck(v.stack.Peek(i)) {
			// TODO: handle panics
			panic(TypeError{t, TypeOf(v)})
		}
	}
}

func (v *VM) Call() {
	// TODO: don't panic
	switch f := v.Pop().(type) {
	case Function:
		v.checkTypes(f.Sig.Args)

		v.scope = v.scope.Child()
		code := v.code
		v.code = bytes.NewReader(f.Code)
		defer func() {
			v.code = code
			v.scope = v.scope.parent
		}()
		v.exec()
		v.checkTypes(f.Sig.Ret)

	case Builtin:
		v.checkTypes(f.Sig.Args)
		rets := f.F(v.stack.PopN(len(f.Sig.Args))...)
		v.stack = append(v.stack, rets...)
		v.checkTypes(f.Sig.Ret)
	}
}

func (v *VM) Func(sig TypeSignature, code []byte) {
	v.Push(Function{sig, code})
}

func (v *VM) Builtin(sig TypeSignature, f func (...Value) []Value) {
	v.Push(Builtin{sig, f})
}
