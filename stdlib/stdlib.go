package stdlib

import (
	"fmt"
	"strconv"
	"strings"
	"../types"
	"../codegen"
)

type FuncDef struct {
	name string
	f    func(...types.Value) []types.Value
}

func (d FuncDef) Name() types.Symbol {
	return types.Symbol(d.name)
}

func values(args ...types.Value) []types.Value {
	return args
}

func (d FuncDef) Builtin() types.Builtin {
	sig := codegen.Sig(strings.SplitN(d.name, ":", 2)[1])
	return types.Builtin{sig, d.f}
}

var Functions = []FuncDef{
	FuncDef{"Println:string", func(a... types.Value) []types.Value {
		s := a[0].(string)
		fmt.Println(s)
		return nil
	}},

	// FIXME: errors are not yet implemented
	//FuncDef{"ToInt:string->int:error", func(a... types.Value) []types.Value {
	//	s := a[0].(string)
	//	return values(strconv.Atoi(s))
	//}},

	FuncDef{"ToString:int->string", func(a... types.Value) []types.Value {
		i := a[0].(int)
		return values(strconv.Itoa(i))
	}},

	// TODO: many more stdlib functions need implemented
}
