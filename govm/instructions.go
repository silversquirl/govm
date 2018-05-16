/*
	First nybble: section
	Second nybble: instruction

	Sections:

	0 - Jump
	1 - Stack
	2 - Arithmetic
	3 - Logic
	4 - Bitwise
	5 - Functions
	6 -
	7 -
	8 -
	9 -
	a -
	b -
	c -
	d -
	e -
	f -
*/

package govm

type Instruction byte

const (
	J   Instruction = 0x00
	JT  Instruction = 0x01
	JF  Instruction = 0x02
	JZ  Instruction = 0x03
	JNz Instruction = 0x04

	Push Instruction = 0x10
	Pop  Instruction = 0x11
	Dup  Instruction = 0x12
	Swp  Instruction = 0x13
	Set  Instruction = 0x14
	Get  Instruction = 0x15

	Inc Instruction = 0x20
	Dec Instruction = 0x21
	Add Instruction = 0x22
	Sub Instruction = 0x23
	Mul Instruction = 0x24
	Div Instruction = 0x25

	EQ Instruction = 0x30
	NE Instruction = 0x31
	LT Instruction = 0x32
	GT Instruction = 0x33
	LE Instruction = 0x34
	GE Instruction = 0x35

	And Instruction = 0x36
	Or  Instruction = 0x37
	Xor Instruction = 0x38
	Not Instruction = 0x39

	BAnd Instruction = 0x40
	BOr  Instruction = 0x41
	BXor Instruction = 0x42
	BNot Instruction = 0x43
	BLS  Instruction = 0x44
	BRS  Instruction = 0x45

	BSet  Instruction = 0x46
	BClr  Instruction = 0x47
	BTgl  Instruction = 0x48
	BMtch Instruction = 0x49

	Call Instruction = 0x50
	Func Instruction = 0x51
)
