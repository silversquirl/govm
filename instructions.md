# govm Instruction Set

The syntax used in this document is as follows:

    opcode:...:stack2ndtype:stacktoptype->...:newstack2ndtype:newstacktoptype (operand1type:operand2type:...)

Or a variant such as (but not limited to) one of:

    opcode:...:stack2ndtype:stacktoptype->...:newstack2ndtype:newstacktoptype
    opcode:...:stack2ndtype:stacktoptype
    opcode->...:newstack2ndtype:newstacktoptype
    opcode->newstack2ndtype:newstacktoptype (operand1type:operand2type:...)

## Gotos

The `label` type is not an actual type. In govm IR, this will be the name
of a local label. In govm bytecode, this will be a byte offset from the
current position. A positive offset of `n` will jump forward by `n` bytes,
whereas a negative offset of `n` will jump backward by `-n` bytes.

- `j (label)` Jump
- `jt:bool (label)` Jump true
- `jf:bool (label)` Jump false
- `jz:int (label)` Jump zero
- `jnz:int (label)` Jump nonzero

## Stack operations

T and T1 are stand-ins for any type. Multiple occurrences of T or T1 refer
to the same concrete type.

- `push->T (T)` In govm bytecode, a type for T is placed before the operand value
- `pop:T`
- `dup:T->T:T`
- `swp:T:T1->T1:T`
- `set:T (symbol)`
- `get->T (symbol)`

## Arithmetic

- `inc:int->int`
- `dec:int->int`

`N1`, `N2` and `N3` are `int` or `float`. If either `N1` or `N2` is
`float`, `N3` is also float, otherwise it is `int`.

- `add:N1:N2->N3`
- `sub:N1:N2->N3`
- `mul:N1:N2->N3`
- `div:N1:N2->N3`
- `mod:int:int->int`

## Logic

### Comparison

`T` is a stand in for a comparable type (`int`, `float` or `string`). Two
occurrences of `T` refer to the same type. There are also variants for
comparing `int` and `float` values (but beware of floating point
inaccuracies when testing equality).

- `eq:T:T->bool`
- `ne:T:T->bool`
- `lt:T:T->bool`
- `gt:T:T->bool`
- `le:T:T->bool`
- `ge:T:T->bool`

### Boolean operations

- `and:bool:bool->bool`
- `or:bool:bool->bool`
- `xor:bool:bool->bool`
- `not:bool->bool`

## Bitwise operators

Though there's often little use for these in scripting languages, we have
them anyway for completeness.

- `band:int:int->int`
- `bor:int:int->int`
- `bxor:int:int->int`
- `bnot:int->int`
- `bls:int:int->int`
- `brs:int:int->int`

We even have some unique operators for bit fiddling that many low-level
languages don't even support.

- `bset:int:int->int`
- `bclr:int:int->int`
- `btgl:int:int->int`

And an operator for checking if a bitmask matches an int:

- `bmtch:int:int->bool`

## Functions

- `call:func:<args>-><rets>`
- `ret`

This instruction is used for creating functions. It takes a type signature
as its first operand and the number of bytes until the end of the function
code as its second.

In govm bytecode, the type signature takes the form
`nargs arg1type arg2type... nret ret1type ret2type...` where `nargs` and
`nret` are ints.

In govm IR, the type signature takes the form
`:arg1type:arg2type:...->ret1type:ret2type:...`. The arg list or return
list may be omitted, but to omit both a single colon (`:`) must be used as
the signature.

In govm IR, the number of bytes is omitted and the end of the function is
specified using `end func`.

- `func->func (type signature:int:byte...)`
