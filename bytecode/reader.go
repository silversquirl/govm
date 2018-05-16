package bytecode

import (
	"bytes"
	"encoding/binary"
	"io"
	"go.vktec.org.uk/govm/types"
)

type Reader struct { R io.ReadSeeker }

func NewSliceReader(code []byte) Reader {
	return Reader{bytes.NewReader(code)}
}

func (r Reader) ReadByte() (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r.R, buf)
	return buf[0], err
}

func (r Reader) Seek(off int64, whence int) (int64, error) {
	return r.R.Seek(off, whence)
}

func (r Reader) Int() (int, error) {
	var i int32
	err := binary.Read(r.R, binary.BigEndian, &i)
	return int(i), err
}

func (r Reader) Float() (f float64, err error) {
	err = binary.Read(r.R, binary.BigEndian, &f)
	return
}

func (r Reader) Bool() (b bool, err error) {
	err = binary.Read(r.R, binary.BigEndian, &b)
	return
}

func (r Reader) Bytes() ([]byte, error) {
	l, err := r.Int()
	if err != nil {
		return nil, err
	}
	buf := make([]byte, l)
	if _, err := io.ReadFull(r.R, buf); err != nil {
		return nil, err
	}
	if len(buf) < l {
		return nil, io.ErrUnexpectedEOF
	}
	return buf, nil
}

func (r Reader) String() (string, error) {
	s, err := r.Bytes()
	return string(s), err
}

func (r Reader) Type() (types.Type, error) {
	b, err := r.ReadByte()
	if err == io.EOF {
		err = io.ErrUnexpectedEOF
	}
	if err != nil {
		return types.Type{}, err
	}

	k := types.Kind(b)
	switch k {
	case types.Int, types.Float, types.Bool, types.String:
		return types.Type{k, types.TypeSignature{}, 0}, nil
	case types.Struct:
		if i, err := r.Int(); err != nil {
			return types.Type{}, err
		} else {
			return types.Type{k, types.TypeSignature{}, i}, nil
		}
	default:
		panic("Unknown kind")
	}
}

func (r Reader) TypedValue() (types.Value, error) {
	t, err := r.Type()
	if err != nil {
		return nil, err
	}

	switch t.Kind {
	case types.Int:
		return r.Int()
	case types.Float:
		return r.Float()
	case types.Bool:
		return r.Bool()
	case types.String:
		return r.String()
	case types.Struct:
		panic("structs not implemented")
	default:
		panic("Unknown type")
	}
}

func (r Reader) TypeSignature() (types.TypeSignature, error) {
	var ts types.TypeSignature
	nargs, err := r.Int()
	if err != nil {
		return types.TypeSignature{}, err
	}
	ts.Args = make([]types.Type, nargs)
	for i := 0; i < nargs; i++ {
		if t, err := r.Type(); err != nil {
			return types.TypeSignature{}, err
		} else {
			ts.Args[i] = t
		}
	}

	nret, err := r.Int()
	if err != nil {
		return types.TypeSignature{}, err
	}
	ts.Ret = make([]types.Type, nret)
	for i := 0; i < nret; i++ {
		if t, err := r.Type(); err != nil {
			return types.TypeSignature{}, err
		} else {
			ts.Ret[i] = t
		}
	}

	return ts, nil
}
