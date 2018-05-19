package codegen

import (
	"../types"
	"strings"
)

func Typ(s string) (t types.Type) {
	switch s {
	case "int":
		return types.TypeInt
	case "float":
		return types.TypeFloat
	case "bool":
		return types.TypeBool
	case "string":
		return types.TypeString
	}
	if strings.HasPrefix(s, "func(") && strings.HasSuffix(s, ")") {
		t := types.TypeFunc
		t.Sig = Sig(strings.TrimPrefix(strings.TrimSuffix(s, ")"), "func("))
		return t
	}
	panic("Unknown type")
}

func Sig(s string) (ts types.TypeSignature) {
	s = strings.TrimPrefix(s, ":")
	if s == "" {
		return
	}
	sections := strings.SplitN(s, "->", 2)
	for _, a := range strings.Split(sections[0], ":") {
		ts.Args = append(ts.Args, Typ(a))
	}
	if len(sections) > 1 {
		for _, r := range strings.Split(sections[1], ":") {
			ts.Ret = append(ts.Ret, Typ(r))
		}
	}
	return
}
