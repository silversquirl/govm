# govm

govm is a semi-statically typed bytecode VM written in Go. It utilises
Go's garbage collector rather than implementing its own, simplifying the
implementation and reducing the size of the library when linking to Go
programs.

## Source code layout

govm contains a number of different packages with different purposes. Here
is a summary of how the source code is layed out:

- `bytecode/` A package for reading and writing bytecode
- `codegen/` A package for generating GVB code
- `doc/` Documentation of the VM's internals
	- `doc/bytecode.md` Documentation of GVB, govm's bytecode representation
	- `doc/instructions.md` Documentation of the VM's instruction set
- `examples/` Example programs written in GVA, govm's assembly-like IR
- `gvas/` The govm assembler. Converts from GVA to GVB
- `gvi/` A CLI for the VM. Allows running GVB files from the command line
- `opcode/` A package containing constants for each opcode byte
- `stdlib/` The standard library
- `types/` Types used in many places throughout the project
- `./vm.go` The core VM package. Interprets GVB
- `./*_test.go` Tests for the VM
