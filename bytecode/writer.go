package bytecode

import (
	"encoding/binary"
	"io"
	"go.vktec.org.uk/govm/types"
)

func SizeOfType(t types.Type) int {
	switch t.Kind {
	case types.Int, types.Float, types.Bool, types.String:
		return 1 // Single-byte representation
	case types.FuncT:
		return 1 /* kind */ + SizeOf(t.Sig) + 4 /* int for length of body */
	case types.Struct:
		panic("Structs are not implemented")
	default:
		panic("Unknown kind")
	}
}

func SizeOf(val types.Value) int {
	switch val := val.(type) {
	case int, *int: // *int is a codegen label
		return 4
	case float64:
		return 8
	case bool:
		return 1
	case string:
		n := len(val)
		return SizeOf(n) + n
	case types.Type:
		return SizeOfType(val)
	case types.TypeSignature:
		s := 8 // 2 ints
		for _, t := range val.Args {
			s += SizeOfType(t)
		}
		for _, t := range val.Ret {
			s += SizeOfType(t)
		}
		return s
	default:
		panic("Unknown type")
	}
}

type Writer struct {
	w io.Writer
	off int
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w, 0}
}

func (w *Writer) Write(p []byte) (int, error) {
	off, err := w.w.Write(p)
	w.off += off
	return off, err
}

func (w *Writer) WriteByte(b byte) error {
	_, err := w.Write([]byte{b})
	return err
}

func (w *Writer) Int(i int) error {
	return binary.Write(w, binary.BigEndian, int32(i))
}

func (w *Writer) Float(f float64) error {
	return binary.Write(w, binary.BigEndian, f)
}

func (w *Writer) Bool(b bool) error {
	return binary.Write(w, binary.BigEndian, b)
}

func (w *Writer) Bytes(b []byte) error {
	if err := w.Int(len(b)); err != nil {
		return err
	}
	_, err := w.Write(b)
	return err
}

func (w *Writer) String(s string) error {
	return w.Bytes([]byte(s))
}

func (w *Writer) Type(t types.Type) error {
	switch t.Kind {
	case types.Int, types.Float, types.Bool, types.String:
	case types.Struct:
		panic("Cannot write struct type")
	case types.FuncT:
		panic("Cannot write function type")
	}
	return w.WriteByte(byte(t.Kind))
}

func (w *Writer) Value(val types.Value) error {
	switch val := val.(type) {
	case int:
		return w.Int(val)
	case float64:
		return w.Float(val)
	case bool:
		return w.Bool(val)
	case string:
		return w.String(val)

	// These two are mainly use from the codegen package
	case *int: // This is a label
		return w.Int(*val - w.off - 4 /* -4 because we jump from the end of the int not the beginning of it */)
	case types.TypeSignature:
		return w.TypeSignature(val)

	default:
		panic("Unknown type")
	}
}

func (w *Writer) TypedValue(val types.Value) error {
	if err := w.Type(types.TypeOf(val)); err != nil {
		return err
	}
	return w.Value(val)
}

func (w *Writer) TypeSignature(ts types.TypeSignature) error {
	if err := w.Int(len(ts.Args)); err != nil {
		return err
	}
	for _, t := range ts.Args {
		if err := w.Type(t); err != nil {
			return err
		}
	}
	if err := w.Int(len(ts.Ret)); err != nil {
		return err
	}
	for _, t := range ts.Ret {
		if err := w.Type(t); err != nil {
			return err
		}
	}
	return nil
}
