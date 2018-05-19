package codegen

import (
	"testing"
	"../opcode"
	"../types"
)

func TestCodegen(t *testing.T) {
	// Input
	g := New()
	mainEnd := new(int)
	g.Func(Sig(":"), mainEnd)
	g.Push("Hello, world!")
	g.Get("Println:string")
	g.Call()
	g.Label(mainEnd)
	g.Set("Main")
	generated, err := g.Generate()
	if err != nil {
		t.Fatal(err)
	}

	// Expected result
	code := []byte{
		byte(opcode.Func),
		0x00, 0x00, 0x00, 0x00, // 0 args
		0x00, 0x00, 0x00, 0x00, // 0 return values
		0x00, 0x00, 0x00, 0x27, // 39 bytes
		byte(opcode.Push), byte(types.String),
		0x00, 0x00, 0x00, 0x0d, // 13 bytes
	}
	code = append(code, []byte("Hello, world!")...)
	code = append(code,
		byte(opcode.Get),
		0x00, 0x00, 0x00, 0x0e, // 14 bytes
	)
	code = append(code, []byte("Println:string")...)
	code = append(code,
		byte(opcode.Call),
		byte(opcode.Set),
		0x00, 0x00, 0x00, 0x04, // 4 bytes
	)
	code = append(code, []byte("Main")...)

	// Compare the two
	t.Log("Expected:  ", code)
	t.Log("Generated: ", generated)

	if len(code) != len(generated) {
		t.Fatal("Generated code is the wrong length")
	}
	for i := range code {
		if code[i] != generated[i] {
			t.Fatal("Generated code differs at byte ", i)
		}
	}
}
