package codegen

import (
	"bytes"
	"go.vktec.org.uk/govm/bytecode"
	"go.vktec.org.uk/govm/opcode"
	"go.vktec.org.uk/govm/types"
)

type Instruction struct {
	Opcode byte
	Operands []types.Value
}

type Generator struct {
	i []Instruction
	size int // Length of bytecode so far
}

func New() Generator {
	return Generator{}
}

func (g Generator) Generate() ([]byte, error) {
	buf := bytes.Buffer{}
	bw := bytecode.NewWriter(&buf)
	for _, i := range g.i {
		if err := bw.WriteByte(i.Opcode); err != nil {
			return nil, err
		}
		switch i.Opcode {
		case opcode.Push: // This instruction's operand is typed
			for _, v := range i.Operands {
				if err := bw.TypedValue(v); err != nil {
					return nil, err
				}
			}
		default:
			for _, v := range i.Operands {
				if err := bw.Value(v); err != nil {
					return nil, err
				}
			}
		}
	}
	return buf.Bytes(), nil
}

func (g *Generator) Instr(code byte, operands... types.Value) {
	g.i = append(g.i, Instruction{code, operands})
	g.size++
	for _, val := range operands {
		g.size += bytecode.SizeOf(val)
		if code == opcode.Push {
			g.size += bytecode.SizeOfType(types.TypeOf(val)) // Typed operand
		}
	}
}

func (g *Generator) Label(lbl *int) *int {
	if lbl == nil {
		lbl = new(int)
	}
	*lbl = g.size
	return lbl
}

func (g *Generator) J(lbl *int) {
	g.Instr(opcode.J, lbl)
}

func (g *Generator) JT(lbl *int) {
	g.Instr(opcode.JT, lbl)
}

func (g *Generator) JF(lbl *int) {
	g.Instr(opcode.JF, lbl)
}

func (g *Generator) JZ(lbl *int) {
	g.Instr(opcode.JZ, lbl)
}

func (g *Generator) JNz(lbl *int) {
	g.Instr(opcode.JNz, lbl)
}

func (g *Generator) Push(val types.Value) {
	g.Instr(opcode.Push, val)
}

func (g *Generator) Pop() {
	g.Instr(opcode.Pop)
}

func (g *Generator) Dup() {
	g.Instr(opcode.Dup)
}

func (g *Generator) Swp() {
	g.Instr(opcode.Swp)
}

func (g *Generator) Set(s string) {
	g.Instr(opcode.Set, s)
}

func (g *Generator) Get(s string) {
	g.Instr(opcode.Get, s)
}

func (g *Generator) Inc() {
	g.Instr(opcode.Inc)
}

func (g *Generator) Dec() {
	g.Instr(opcode.Dec)
}

func (g *Generator) Add() {
	g.Instr(opcode.Add)
}

func (g *Generator) Sub() {
	g.Instr(opcode.Sub)
}

func (g *Generator) Mul() {
	g.Instr(opcode.Mul)
}

func (g *Generator) Div() {
	g.Instr(opcode.Div)
}

func (g *Generator) Mod() {
	g.Instr(opcode.Mod)
}

func (g *Generator) EQ() {
	g.Instr(opcode.EQ)
}

func (g *Generator) NE() {
	g.Instr(opcode.NE)
}

func (g *Generator) LT() {
	g.Instr(opcode.LT)
}

func (g *Generator) GT() {
	g.Instr(opcode.GT)
}

func (g *Generator) LE() {
	g.Instr(opcode.LE)
}

func (g *Generator) GE() {
	g.Instr(opcode.GE)
}

func (g *Generator) And() {
	g.Instr(opcode.And)
}

func (g *Generator) Or () {
	g.Instr(opcode.Or )
}

func (g *Generator) Xor() {
	g.Instr(opcode.Xor)
}

func (g *Generator) Not() {
	g.Instr(opcode.Not)
}

func (g *Generator) BAnd() {
	g.Instr(opcode.BAnd)
}

func (g *Generator) BOr () {
	g.Instr(opcode.BOr )
}

func (g *Generator) BXor() {
	g.Instr(opcode.BXor)
}

func (g *Generator) BNot() {
	g.Instr(opcode.BNot)
}

func (g *Generator) BLS () {
	g.Instr(opcode.BLS)
}

func (g *Generator) BRS () {
	g.Instr(opcode.BRS)
}

func (g *Generator) BSet () {
	g.Instr(opcode.BSet)
}

func (g *Generator) BClr () {
	g.Instr(opcode.BClr)
}

func (g *Generator) BTgl () {
	g.Instr(opcode.BTgl)
}

func (g *Generator) BMtch() {
	g.Instr(opcode.BMtch)
}

func (g *Generator) Call() {
	g.Instr(opcode.Call)
}

func (g *Generator) Ret() {
	g.Instr(opcode.Ret)
}

func (g *Generator) Func(ts types.TypeSignature, lbl *int) {
	g.Instr(opcode.Func, ts, lbl)
}
