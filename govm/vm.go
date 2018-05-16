package govm

import (
	"bytes"
	"encoding/binary"
	"io"
)

type VM struct {
	stack Stack
	scope *Scope
	code  *bytes.Reader // The bit of code being currently executed
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
			if err := v.Jump(off); err != nil {
				return err
			}

		case JT:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			if err := v.JumpTrue(off); err != nil {
				return err
			}

		case JF:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			if err := v.JumpFalse(off); err != nil {
				return err
			}

		case JZ:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			if err := v.JumpZero(off); err != nil {
				return err
			}

		case JNz:
			off, err := v.readI()
			if err != nil {
				return nil
			}
			if err := v.JumpNonzero(off); err != nil {
				return err
			}

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
			if err := v.Inc(); err != nil {
				return err
			}
		case Dec:
			if err := v.Dec(); err != nil {
				return err
			}
		case Add:
			if err := v.Add(); err != nil {
				return err
			}
		case Sub:
			if err := v.Sub(); err != nil {
				return err
			}
		case Mul:
			if err := v.Mul(); err != nil {
				return err
			}
		case Div:
			if err := v.Div(); err != nil {
				return err
			}

		case EQ:
			if err := v.EQ(); err != nil {
				return err
			}
		case NE:
			if err := v.NE(); err != nil {
				return err
			}
		case LT:
			if err := v.LT(); err != nil {
				return err
			}
		case GT:
			if err := v.GT(); err != nil {
				return err
			}
		case LE:
			if err := v.LE(); err != nil {
				return err
			}
		case GE:
			if err := v.GE(); err != nil {
				return err
			}

		case And:
			if err := v.And(); err != nil {
				return err
			}
		case Or:
			if err := v.Or(); err != nil {
				return err
			}
		case Xor:
			if err := v.Xor(); err != nil {
				return err
			}
		case Not:
			if err := v.Not(); err != nil {
				return err
			}

		case BAnd:
			if err := v.BAnd(); err != nil {
				return err
			}
		case BOr:
			if err := v.BOr(); err != nil {
				return err
			}
		case BXor:
			if err := v.BXor(); err != nil {
				return err
			}
		case BNot:
			if err := v.BNot(); err != nil {
				return err
			}
		case BLS:
			if err := v.BLS(); err != nil {
				return err
			}
		case BRS:
			if err := v.BRS(); err != nil {
				return err
			}
		case BSet:
			if err := v.BSet(); err != nil {
				return err
			}
		case BClr:
			if err := v.BClr(); err != nil {
				return err
			}
		case BTgl:
			if err := v.BTgl(); err != nil {
				return err
			}
		case BMtch:
			if err := v.BMtch(); err != nil {
				return err
			}

		case Call:
			if err := v.Call(); err != nil {
				return err
			}

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

func (v *VM) Jump(off int) error {
	_, err := v.code.Seek(int64(off), io.SeekCurrent)
	return err
}

func (v *VM) JumpTrue(off int) error {
	switch val := v.Pop().(type) {
	case bool:
		if val {
			return v.Jump(off)
		}
	default:
		return TypeError{TypeBool, TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpFalse(off int) error {
	switch val := v.Pop().(type) {
	case bool:
		if !val {
			return v.Jump(off)
		}
	default:
		return TypeError{TypeBool, TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpZero(off int) error {
	switch val := v.Pop().(type) {
	case int:
		if val == 0 {
			return v.Jump(off)
		}
	case float64:
		if val == 0.0 {
			return v.Jump(off)
		}
	default:
		return TypeError{TypeBool, TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpNonzero(off int) error {
	switch val := v.Pop().(type) {
	case int:
		if val != 0 {
			return v.Jump(off)
		}
	case float64:
		if val != 0.0 {
			return v.Jump(off)
		}
	default:
		return TypeError{TypeBool, TypeOf(val)}
	}
	return nil
}

func (v *VM) Push(val Value) {
	v.stack.Push(val)
}

func (v *VM) Pop() Value {
	// TODO: errors
	return v.stack.Pop()
}

func (v *VM) Dup() {
	// TODO: errors
	v.stack.Dup()
}

func (v *VM) Swap() {
	// TODO: errors
	v.stack.Swap()
}

func (v *VM) Set(s Symbol) {
	v.scope.Set(s, v.Pop())
}

func (v *VM) Get(s Symbol) error {
	val, err := v.scope.Get(s)
	if err != nil {
		return err
	}
	v.Push(val)
	return nil
}

func (v *VM) Inc() error {
	switch val := v.Pop().(type) {
	case int:
		v.Push(val + 1)
		return nil
	default:
		return TypeError{TypeInt, TypeOf(val)}
	}
}

func (v *VM) Dec() error {
	switch val := v.Pop().(type) {
	case int:
		v.Push(val - 1)
		return nil
	default:
		return TypeError{TypeInt, TypeOf(val)}
	}
}

func (v *VM) Add() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a + b)
		case float64:
			v.Push(float64(a) + b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a + float64(b))
		case float64:
			v.Push(a + b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) Sub() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a - b)
		case float64:
			v.Push(float64(a) - b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a - float64(b))
		case float64:
			v.Push(a - b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) Mul() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a * b)
		case float64:
			v.Push(float64(a) * b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a * float64(b))
		case float64:
			v.Push(a * b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) Div() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a / b)
		case float64:
			v.Push(float64(a) / b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a / float64(b))
		case float64:
			v.Push(a / b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) EQ() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a == b)
		case float64:
			v.Push(float64(a) == b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a == float64(b))
		case float64:
			v.Push(a == b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) NE() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a != b)
		case float64:
			v.Push(float64(a) != b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a != float64(b))
		case float64:
			v.Push(a != b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) LT() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a < b)
		case float64:
			v.Push(float64(a) < b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a < float64(b))
		case float64:
			v.Push(a < b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) GT() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a > b)
		case float64:
			v.Push(float64(a) > b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a > float64(b))
		case float64:
			v.Push(a > b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) LE() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a <= b)
		case float64:
			v.Push(float64(a) <= b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a <= float64(b))
		case float64:
			v.Push(a <= b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) GE() error {
	b, a := v.Pop(), v.Pop()
	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a >= b)
		case float64:
			v.Push(float64(a) >= b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a >= float64(b))
		case float64:
			v.Push(a >= b)
		default:
			return TypeError{TypeNum, TypeOf(b)}
		}

	default:
		return TypeError{TypeNum, TypeOf(b)}
	}
	return nil
}

func (v *VM) And() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) && b.(bool))
	return nil
}

func (v *VM) Or() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) || b.(bool))
	return nil
}

func (v *VM) Xor() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) != b.(bool))
	return nil
}

func (v *VM) Not() error {
	a := v.Pop()
	if err := TypeBool.TypeCheck(a); err != nil {
		return err
	}
	v.Push(!a.(bool))
	return nil
}

func (v *VM) BAnd() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) & b.(int))
	return nil
}

func (v *VM) BOr() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) | b.(int))
	return nil
}

func (v *VM) BXor() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) ^ b.(int))
	return nil
}

func (v *VM) BNot() error {
	a := v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	v.Push(^a.(int))
	return nil
}

func (v *VM) BLS() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) << uint(b.(int)))
	return nil
}

func (v *VM) BRS() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) >> uint(b.(int)))
	return nil
}

func (v *VM) BSet() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) | (1 << uint(b.(int))))
	return nil
}

func (v *VM) BClr() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) &^ (1 << uint(b.(int))))
	return nil
}

func (v *VM) BTgl() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) ^ (1 << uint(b.(int))))
	return nil
}

func (v *VM) BMtch() error {
	b, a := v.Pop(), v.Pop()
	if err := TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int)&b.(int) != 0)
	return nil
}

func (v *VM) checkTypes(types []Type) error {
	for i, t := range types {
		if err := t.TypeCheck(v.stack.Peek(i)); err != nil {
			return err
		}
	}
	return nil
}

func (v *VM) Call() error {
	switch f := v.Pop().(type) {
	case Function:
		if err := v.checkTypes(f.Sig.Args); err != nil {
			return err
		}

		v.scope = v.scope.Child()
		code := v.code
		v.code = bytes.NewReader(f.Code)
		defer func() {
			v.code = code
			v.scope = v.scope.parent
		}()
		v.exec()
		if err := v.checkTypes(f.Sig.Ret); err != nil {
			return err
		}

	case Builtin:
		if err := v.checkTypes(f.Sig.Args); err != nil {
			return err
		}
		rets := f.F(v.stack.PopN(len(f.Sig.Args))...)
		v.stack = append(v.stack, rets...)
		if err := v.checkTypes(f.Sig.Ret); err != nil {
			return err
		}

	default:
		return TypeError{TypeFunc, TypeOf(f)}
	}
	return nil
}

func (v *VM) Func(sig TypeSignature, code []byte) {
	v.Push(Function{sig, code})
}

func (v *VM) Builtin(sig TypeSignature, f func(...Value) []Value) {
	v.Push(Builtin{sig, f})
}
