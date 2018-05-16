package main

import (
	"./govm"
	"fmt"
)

func main() {
	// See hello.gvm for a more readable version
	code := []byte{
		byte(govm.Func),
		0x00, 0x00, 0x00, 0x00, // 0 args
		0x00, 0x00, 0x00, 0x00, // 0 return values
		0x00, 0x00, 0x00, 0x27, // 39 bytes
		byte(govm.Push), byte(govm.String),
		0x00, 0x00, 0x00, 0x0d, // 13 bytes
	}
	code = append(code, []byte("Hello, world!")...)
	code = append(code,
		byte(govm.Get),
		0x00, 0x00, 0x00, 0x0e, // 14 bytes
	)
	code = append(code, []byte("Println:string")...)
	code = append(code,
		byte(govm.Call),
		byte(govm.Set),
		0x00, 0x00, 0x00, 0x04, // 4 bytes
	)
	code = append(code, []byte("Main")...)

	v := govm.NewVM()
	v.Load(code)

	v.Builtin(govm.TypeSignature{[]govm.Type{govm.TypeString}, nil}, func (a ...govm.Value) []govm.Value {
		s := a[0].(string)
		fmt.Println(s)
		return nil
	})
	v.Set(govm.Symbol("Println:string"))

	v.Get(govm.Symbol("Main"))
	v.Call()
}
