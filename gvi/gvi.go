package main

import (
	"flag"
	"fmt"
	"go.vktec.org.uk/govm"
	"go.vktec.org.uk/govm/codegen"
	"go.vktec.org.uk/govm/types"
	"io"
	"os"
)

func Main() int {
	flag.Parse()

	var input io.ReadSeeker
	if len(flag.Args()) > 0 {
		var err error
		input, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stdout, err)
			return 1
		}
	} else {
		input = os.Stdin
	}

	vm := govm.New()

	vm.Builtin(codegen.Sig(":string"), func(a ...types.Value) []types.Value {
		s := a[0].(string)
		fmt.Println(s)
		return nil
	})
	vm.Set(types.Symbol("Println:string"))

	vm.Builtin(codegen.Sig(":int->string"), func(a ...types.Value) []types.Value {
		i := a[0].(int)
		return []types.Value{fmt.Sprintf("%d", i)}
	})
	vm.Set(types.Symbol("ToString:int->string"))

	vm.LoadFrom(input)
	if err := vm.Get("Main:"); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	if err := vm.Call(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}

	return 0
}

func main() {
	os.Exit(Main())
}
