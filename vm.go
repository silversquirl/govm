package govm

import (
	"io"
	"go.vktec.org.uk/govm/bytecode"
	"go.vktec.org.uk/govm/opcode"
	"go.vktec.org.uk/govm/types"
)

type VM struct {
	stack types.Stack
	scope *types.Scope
	code  *bytecode.Reader
}

func NewVM() (v VM) {
	v.scope = &types.Scope{}
	return
}

func (v *VM) Load(code []byte) error {
	v.code = bytecode.NewSliceReader(code)
	if err := v.exec(); err != nil {
		return err
	}
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

		switch op {
		case opcode.J:
			off, err := v.code.Int()
			if err != nil {
				return err
			}
			if err := v.Jump(off); err != nil {
				return err
			}

		case opcode.JT:
			off, err := v.code.Int()
			if err != nil {
				return err
			}
			if err := v.JumpTrue(off); err != nil {
				return err
			}

		case opcode.JF:
			off, err := v.code.Int()
			if err != nil {
				return err
			}
			if err := v.JumpFalse(off); err != nil {
				return err
			}

		case opcode.JZ:
			off, err := v.code.Int()
			if err != nil {
				return err
			}
			if err := v.JumpZero(off); err != nil {
				return err
			}

		case opcode.JNz:
			off, err := v.code.Int()
			if err != nil {
				return err
			}
			if err := v.JumpNonzero(off); err != nil {
				return err
			}

		case opcode.Push:
			val, err := v.code.TypedValue()
			if err != nil {
				return err
			}
			v.Push(val)

		case opcode.Pop:
			if _, err := v.Pop(); err != nil {
				return err
			}

		case opcode.Dup:
			if err := v.Dup(); err != nil {
				return err
			}

		case opcode.Swp:
			if err := v.Swap(); err != nil {
				return err
			}

		case opcode.Set:
			s, err := v.code.String()
			if err != nil {
				return err
			}
			v.Set(types.Symbol(s))

		case opcode.Get:
			s, err := v.code.String()
			if err != nil {
				return err
			}
			if err := v.Get(types.Symbol(s)); err != nil {
				return err
			}

		case opcode.Inc:
			if err := v.Inc(); err != nil {
				return err
			}
		case opcode.Dec:
			if err := v.Dec(); err != nil {
				return err
			}
		case opcode.Add:
			if err := v.Add(); err != nil {
				return err
			}
		case opcode.Sub:
			if err := v.Sub(); err != nil {
				return err
			}
		case opcode.Mul:
			if err := v.Mul(); err != nil {
				return err
			}
		case opcode.Div:
			if err := v.Div(); err != nil {
				return err
			}
		case opcode.Mod:
			if err := v.Mod(); err != nil {
				return err
			}

		case opcode.EQ:
			if err := v.EQ(); err != nil {
				return err
			}
		case opcode.NE:
			if err := v.NE(); err != nil {
				return err
			}
		case opcode.LT:
			if err := v.LT(); err != nil {
				return err
			}
		case opcode.GT:
			if err := v.GT(); err != nil {
				return err
			}
		case opcode.LE:
			if err := v.LE(); err != nil {
				return err
			}
		case opcode.GE:
			if err := v.GE(); err != nil {
				return err
			}

		case opcode.And:
			if err := v.And(); err != nil {
				return err
			}
		case opcode.Or:
			if err := v.Or(); err != nil {
				return err
			}
		case opcode.Xor:
			if err := v.Xor(); err != nil {
				return err
			}
		case opcode.Not:
			if err := v.Not(); err != nil {
				return err
			}

		case opcode.BAnd:
			if err := v.BAnd(); err != nil {
				return err
			}
		case opcode.BOr:
			if err := v.BOr(); err != nil {
				return err
			}
		case opcode.BXor:
			if err := v.BXor(); err != nil {
				return err
			}
		case opcode.BNot:
			if err := v.BNot(); err != nil {
				return err
			}
		case opcode.BLS:
			if err := v.BLS(); err != nil {
				return err
			}
		case opcode.BRS:
			if err := v.BRS(); err != nil {
				return err
			}
		case opcode.BSet:
			if err := v.BSet(); err != nil {
				return err
			}
		case opcode.BClr:
			if err := v.BClr(); err != nil {
				return err
			}
		case opcode.BTgl:
			if err := v.BTgl(); err != nil {
				return err
			}
		case opcode.BMtch:
			if err := v.BMtch(); err != nil {
				return err
			}

		case opcode.Call:
			if err := v.Call(); err != nil {
				return err
			}

		case opcode.Ret:
			// Doesn't make sense to have a separate function
			return types.Return

		case opcode.Func:
			sig, err := v.code.TypeSignature()
			code, err := v.code.Bytes()
			if err != nil {
				return err
			}
			v.Func(sig, code)

		default:
			panic("Unknown opcode")
		}
	}
}

func (v *VM) Jump(off int) error {
	_, err := v.code.Seek(int64(off), io.SeekCurrent)
	return err
}

func (v *VM) JumpTrue(off int) error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case bool:
		if val {
			return v.Jump(off)
		}
	default:
		return types.TypeError{types.TypeBool, types.TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpFalse(off int) error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case bool:
		if !val {
			return v.Jump(off)
		}
	default:
		return types.TypeError{types.TypeBool, types.TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpZero(off int) error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case int:
		if val == 0 {
			return v.Jump(off)
		}
	case float64:
		if val == 0.0 {
			return v.Jump(off)
		}
	default:
		return types.TypeError{types.TypeBool, types.TypeOf(val)}
	}
	return nil
}

func (v *VM) JumpNonzero(off int) error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case int:
		if val != 0 {
			return v.Jump(off)
		}
	case float64:
		if val != 0.0 {
			return v.Jump(off)
		}
	default:
		return types.TypeError{types.TypeBool, types.TypeOf(val)}
	}
	return nil
}

func (v *VM) Push(val types.Value) {
	v.stack.Push(val)
}

func (v *VM) Pop() (types.Value, error) {
	return v.stack.Pop()
}

func (v *VM) Dup() error {
	return v.stack.Dup()
}

func (v *VM) Swap() error {
	return v.stack.Swap()
}

func (v *VM) Set(s types.Symbol) error {
	if val, err := v.Pop(); err == nil {
		v.scope.Set(s, val)
	} else {
		return err
	}
	return nil
}

func (v *VM) Get(s types.Symbol) error {
	val, err := v.scope.Get(s)
	if err != nil {
		return err
	}
	v.Push(val)
	return nil
}

func (v *VM) Inc() error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case int:
		v.Push(val + 1)
		return nil
	default:
		return types.TypeError{types.TypeInt, types.TypeOf(val)}
	}
}

func (v *VM) Dec() error {
	val, err := v.Pop()
	if err != nil {
		return err
	}
	switch val := val.(type) {
	case int:
		v.Push(val - 1)
		return nil
	default:
		return types.TypeError{types.TypeInt, types.TypeOf(val)}
	}
}

func (v *VM) Add() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a + b)
		case float64:
			v.Push(float64(a) + b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a + float64(b))
		case float64:
			v.Push(a + b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) Sub() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a - b)
		case float64:
			v.Push(float64(a) - b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a - float64(b))
		case float64:
			v.Push(a - b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) Mul() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a * b)
		case float64:
			v.Push(float64(a) * b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a * float64(b))
		case float64:
			v.Push(a * b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) Div() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a / b)
		case float64:
			v.Push(float64(a) / b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a / float64(b))
		case float64:
			v.Push(a / b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) Mod() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) % b.(int))
	return nil
}

func (v *VM) EQ() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a == b)
		case float64:
			v.Push(float64(a) == b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a == float64(b))
		case float64:
			v.Push(a == b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) NE() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a != b)
		case float64:
			v.Push(float64(a) != b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a != float64(b))
		case float64:
			v.Push(a != b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) LT() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a < b)
		case float64:
			v.Push(float64(a) < b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a < float64(b))
		case float64:
			v.Push(a < b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) GT() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a > b)
		case float64:
			v.Push(float64(a) > b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a > float64(b))
		case float64:
			v.Push(a > b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) LE() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a <= b)
		case float64:
			v.Push(float64(a) <= b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a <= float64(b))
		case float64:
			v.Push(a <= b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) GE() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	switch a := a.(type) {
	case int:
		switch b := b.(type) {
		case int:
			v.Push(a >= b)
		case float64:
			v.Push(float64(a) >= b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	case float64:
		switch b := b.(type) {
		case int:
			v.Push(a >= float64(b))
		case float64:
			v.Push(a >= b)
		default:
			return types.TypeError{types.TypeNum, types.TypeOf(b)}
		}

	default:
		return types.TypeError{types.TypeNum, types.TypeOf(b)}
	}
	return nil
}

func (v *VM) And() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) && b.(bool))
	return nil
}

func (v *VM) Or() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) || b.(bool))
	return nil
}

func (v *VM) Xor() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeBool.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeBool.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(bool) != b.(bool))
	return nil
}

func (v *VM) Not() error {
	a, err := v.Pop()
	if err != nil {
		return err
	}
	if err := types.TypeBool.TypeCheck(a); err != nil {
		return err
	}
	v.Push(!a.(bool))
	return nil
}

func (v *VM) BAnd() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) & b.(int))
	return nil
}

func (v *VM) BOr() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) | b.(int))
	return nil
}

func (v *VM) BXor() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) ^ b.(int))
	return nil
}

func (v *VM) BNot() error {
	a, err := v.Pop()
	if err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	v.Push(^a.(int))
	return nil
}

func (v *VM) BLS() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) << uint(b.(int)))
	return nil
}

func (v *VM) BRS() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) >> uint(b.(int)))
	return nil
}

func (v *VM) BSet() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) | (1 << uint(b.(int))))
	return nil
}

func (v *VM) BClr() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) &^ (1 << uint(b.(int))))
	return nil
}

func (v *VM) BTgl() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int) ^ (1 << uint(b.(int))))
	return nil
}

func (v *VM) BMtch() error {
	b, err := v.Pop()
	if err != nil {
		return err
	}
	a, err := v.Pop()
	if err != nil {
		return err
	}

	if err := types.TypeInt.TypeCheck(a); err != nil {
		return err
	}
	if err := types.TypeInt.TypeCheck(b); err != nil {
		return err
	}
	v.Push(a.(int)&b.(int) != 0)
	return nil
}

func (v *VM) checkTypes(types []types.Type) error {
	for i, t := range types {
		val, err := v.stack.Peek(i)
		if err != nil {
			return err
		}
		if err := t.TypeCheck(val); err != nil {
			return err
		}
	}
	return nil
}

func (v *VM) Call() error {
	f, err := v.Pop()
	if err != nil {
		return err
	}

	switch f := f.(type) {
	case types.Function:
		if err := v.checkTypes(f.Sig.Args); err != nil {
			return err
		}

		v.scope = v.scope.Child()
		code := v.code
		v.code = bytecode.NewSliceReader(f.Code)
		defer func() {
			v.code = code
			v.scope = v.scope.Parent
		}()
		if err := v.exec(); err != types.Return && err != nil {
			return err
		}
		if err := v.checkTypes(f.Sig.Ret); err != nil {
			return err
		}

	case types.Builtin:
		if err := v.checkTypes(f.Sig.Args); err != nil {
			return err
		}
		args, err := v.stack.PopN(len(f.Sig.Args))
		if err != nil {
			return err
		}
		rets := f.F(args...)
		v.stack = append(v.stack, rets...)
		if err := v.checkTypes(f.Sig.Ret); err != nil {
			return err
		}

	default:
		return types.TypeError{types.TypeFunc, types.TypeOf(f)}
	}
	return nil
}

func (v *VM) Func(sig types.TypeSignature, code []byte) {
	v.Push(types.Function{sig, code})
}

func (v *VM) Builtin(sig types.TypeSignature, f func(...types.Value) []types.Value) {
	v.Push(types.Builtin{sig, f})
}
