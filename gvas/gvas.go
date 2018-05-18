package main

import (
	"flag"
	"fmt"
	"io"
	"bufio"
	"os"
	"strings"
	"strconv"
	"go.vktec.org.uk/govm/codegen"
	"go.vktec.org.uk/govm/types"
)

type Converter struct {
	in *bufio.Scanner
	gen codegen.Generator
	labels map[string]*int
}

type InvalidOpcodeError struct { opcode string }

func (e InvalidOpcodeError) Error() string {
	return "Invalid opcode: " + e.opcode
}

type UnknownTokenError struct { tok string }

func (e UnknownTokenError) Error() string {
	return "Unknown token: '" + e.tok + "'"
}

func sep(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n'
}

func scanToken(data []byte, atEOF bool) (advance int, token []byte, err error) {
	for ; advance < len(data) && sep(data[advance]); advance++ {}
	if advance >= len(data) {
		return advance, nil, nil
	}
	start := advance

	if data[advance] == '"' {
stringLoop:
		for advance < len(data) - 1 {
			advance++
			switch data[advance] {
			case '"':
				advance++
				break stringLoop
			case '\\':
				advance++
			}
		}
		token = data[start:advance]
		return
	}

	for ; advance < len(data) && !sep(data[advance]); advance++ {}
	token = data[start:advance]
	if len(token) == 0 {
		token = nil
	}
	return
}

func readOperands(in *bufio.Scanner, n int) ([]string, error) {
	values := make([]string, n)
	for i := 0; i < n; i++ {
		if !in.Scan() {
			if err := in.Err(); err == nil {
				return values, io.ErrUnexpectedEOF
			} else {
				return values, err
			}
		}
		values[i] = in.Text()
	}
	return values, nil
}

func readOperand(in *bufio.Scanner) (string, error) {
	values, err := readOperands(in, 1)
	return values[0], err
}

func readValue(in *bufio.Scanner) (types.Value, error) {
	val, err := readOperand(in)
	if err != nil {
		return nil, err
	}
	if len(val) == 0 {
		return nil, UnknownTokenError{val}
	}

	if val[0] == '"' {
		return val[1:len(val)-1], nil
	} else if i, err := strconv.Atoi(val); err == nil {
		return i, nil
	} else if f, err := strconv.ParseFloat(val, 64); err == nil {
		return f, nil
	}
	return nil, UnknownTokenError{val}
}

func readSym(in *bufio.Scanner) (string, error) {
	val, err := readOperand(in)
	if err != nil {
		return "", err
	}
	if len(val) == 0 || val[0] != '@' {
		return "", UnknownTokenError{val}
	}
	return val[1:], nil
}

func (c *Converter) convertInstruction(opcode string) error {
	opcode = strings.ToLower(opcode)
	switch opcode {
	case "j":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		if c.labels[lbl] == nil {
			c.labels[lbl] = new(int)
		}
		c.gen.J(c.labels[lbl])

	case "jt":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		if c.labels[lbl] == nil {
			c.labels[lbl] = new(int)
		}
		c.gen.JT(c.labels[lbl])

	case "jf":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		if c.labels[lbl] == nil {
			c.labels[lbl] = new(int)
		}
		c.gen.JF(c.labels[lbl])

	case "jz":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		if c.labels[lbl] == nil {
			c.labels[lbl] = new(int)
		}
		c.gen.JZ(c.labels[lbl])

	case "jnz":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		if c.labels[lbl] == nil {
			c.labels[lbl] = new(int)
		}
		c.gen.JNz(c.labels[lbl])

	case "push":
		value, err := readValue(c.in)
		if err != nil {
			return err
		}
		c.gen.Push(value)

	case "pop":
		c.gen.Pop()
	case "dup":
		c.gen.Dup()
	case "swp":
		c.gen.Swp()
	case "set":
		value, err := readSym(c.in)
		if err != nil {
			return err
		}
		c.gen.Set(value)
	case "get":
		value, err := readSym(c.in)
		if err != nil {
			return err
		}
		c.gen.Get(value)

	case "inc":
		c.gen.Inc()
	case "dec":
		c.gen.Dec()
	case "add":
		c.gen.Add()
	case "sub":
		c.gen.Sub()
	case "mul":
		c.gen.Mul()
	case "div":
		c.gen.Div()
	case "mod":
		c.gen.Mod()

	case "eq":
		c.gen.EQ()
	case "ne":
		c.gen.NE()
	case "lt":
		c.gen.LT()
	case "gt":
		c.gen.GT()
	case "le":
		c.gen.LE()
	case "ge":
		c.gen.GE()

	case "and":
		c.gen.And()
	case "or":
		c.gen.Or()
	case "xor":
		c.gen.Xor()
	case "not":
		c.gen.Not()

	case "band":
		c.gen.BAnd()
	case "bor":
		c.gen.BOr()
	case "bxor":
		c.gen.BXor()
	case "bnot":
		c.gen.BNot()
	case "bls":
		c.gen.BLS()
	case "brs":
		c.gen.BRS()

	case "bset":
		c.gen.BSet()
	case "bclr":
		c.gen.BClr()
	case "btgl":
		c.gen.BTgl()
	case "bmtch":
		c.gen.BMtch()

	case "call":
		c.gen.Call()
	case "ret":
		c.gen.Ret()

	// Special cases
	case "func":
		sig, err := readOperand(c.in)
		if err != nil {
			return err
		}
		c.parseFunction(codegen.Sig(sig))

	case "//":
		_, err := readOperand(c.in)
		return err

	case ".":
		lbl, err := readOperand(c.in)
		if err != nil {
			return err
		}
		c.labels[lbl] = c.gen.Label(c.labels[lbl])

	default:
		return InvalidOpcodeError{opcode}
	}
	return nil
}

func (c *Converter) parseToplevel() error {
	for c.in.Scan() {
		opcode := c.in.Text()
		if err := c.convertInstruction(opcode); err != nil {
			return err
		}
	}
	return c.in.Err()
}

func (c *Converter) parseFunction(sig types.TypeSignature) error {
	endLbl := new(int)
	c.gen.Func(sig, endLbl)
	for c.in.Scan() {
		opcode := c.in.Text()
		if opcode == "endfunc" {
			c.gen.Label(endLbl)
			return nil
		}
		if err := c.convertInstruction(opcode); err != nil {
			return err
		}
	}
	if err := c.in.Err(); err == nil {
		return io.ErrUnexpectedEOF
	} else {
		return err
	}
}

func Main() int {
	var input, output string
	flag.StringVar(&output, "o", "", "Output filename")
	flag.Parse()

	var in io.Reader
	var out io.Writer
	var err error

	if flag.NArg() > 0 {
		input = flag.Arg(0)
		in, err = os.Open(input)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}

		if output == "" {
			output = input[:strings.LastIndexByte(input, '.')] + ".gvb"
		}
		out, err = os.Create(output)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			return 1
		}
	} else {
		in = os.Stdin
		out = os.Stdout
	}

	inScan := bufio.NewScanner(in)
	inScan.Split(scanToken)
	c := Converter{inScan, codegen.New(), make(map[string]*int)}
	if err := c.parseToplevel(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		return 1
	}
	c.gen.GenerateTo(out)

	return 0
}

func main() {
	os.Exit(Main())
}
