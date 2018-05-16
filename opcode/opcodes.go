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

package opcode

const (
	J   byte = 0x00
	JT  byte = 0x01
	JF  byte = 0x02
	JZ  byte = 0x03
	JNz byte = 0x04

	Push byte = 0x10
	Pop  byte = 0x11
	Dup  byte = 0x12
	Swp  byte = 0x13
	Set  byte = 0x14
	Get  byte = 0x15

	Inc byte = 0x20
	Dec byte = 0x21
	Add byte = 0x22
	Sub byte = 0x23
	Mul byte = 0x24
	Div byte = 0x25

	EQ byte = 0x30
	NE byte = 0x31
	LT byte = 0x32
	GT byte = 0x33
	LE byte = 0x34
	GE byte = 0x35

	And byte = 0x36
	Or  byte = 0x37
	Xor byte = 0x38
	Not byte = 0x39

	BAnd byte = 0x40
	BOr  byte = 0x41
	BXor byte = 0x42
	BNot byte = 0x43
	BLS  byte = 0x44
	BRS  byte = 0x45

	BSet  byte = 0x46
	BClr  byte = 0x47
	BTgl  byte = 0x48
	BMtch byte = 0x49

	Call byte = 0x50
	Func byte = 0x51
)
