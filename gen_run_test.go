package govm

import (
	"fmt"
	"go.vktec.org.uk/govm/codegen"
	"go.vktec.org.uk/govm/types"
	"testing"
)

func TestGenRun(t *testing.T) {
	// This will be eaiser to read when compared line-by-line to fizzbuzz.gvm

	g := codegen.New()
	fbEnd := new(int)
	g.Func(codegen.Sig(":int->string"), fbEnd)

	fizz := new(int)
	buzz := new(int)
	else_ := new(int)
	endIf := new(int)

	g.Dup()
	g.Push(15)
	g.Mod()
	g.JNz(fizz)
	g.Push("FizzBuzz")
	g.J(endIf)

	g.Label(fizz)
	g.Dup()
	g.Push(3)
	g.Mod()
	g.JNz(buzz)
	g.Push("Fizz")
	g.J(endIf)

	g.Label(buzz)
	g.Dup()
	g.Push(5)
	g.Mod()
	g.JNz(else_)
	g.Push("Buzz")
	g.J(endIf)

	g.Label(else_)
	g.Dup()
	g.Get("ToString:int->string")
	g.Call()
	g.J(endIf)

	g.Label(endIf)
	g.Label(fbEnd)
	g.Set("fizzbuzz:int->string")

	mainEnd := new(int)
	g.Func(codegen.Sig(":"), mainEnd)
	g.Push(1)
	startLoop := g.Label(nil)
	endLoop := new(int)
	g.Dup()
	g.Push(100)
	g.LT()
	g.JF(endLoop)
	g.Dup()
	g.Get("fizzbuzz:int->string")
	g.Call()
	g.Get("Println:string")
	g.Call()
	g.Inc()
	g.J(startLoop)
	g.Label(endLoop)
	g.Label(mainEnd)
	g.Set("Main:")

	code, err := g.Generate()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(code)

	v := NewVM()
	v.Load(code)

	v.Builtin(codegen.Sig(":string"), func(a ...types.Value) []types.Value {
		s := a[0].(string)
		fmt.Println(s)
		return nil
	})
	v.Set(types.Symbol("Println:string"))

	v.Builtin(codegen.Sig(":int->string"), func(a ...types.Value) []types.Value {
		i := a[0].(int)
		return []types.Value{fmt.Sprintf("%d", i)}
	})
	v.Set(types.Symbol("ToString:int->string"))

	if err := v.Get(types.Symbol("Main:")); err != nil {
		t.Fatal("Get:", err)
	}
	if err := v.Call(); err != nil {
		t.Fatal("Call:", err)
	}
}
