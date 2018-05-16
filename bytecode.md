# govm Bytecode

govm's bytecode is very simple. It is made up of single byte instructions
followed by operands. This document focusses on the encoding of the
operands.

## Ints

Stored as big-endian 32-bit signed integer values. Hopefully nobody tries
to make a >2GiB function.

## Floats

Stored as big-endian 64-bit floats.

## Bools

Stored as one byte with 0x01 being true and 0x00 being false.

## Strings and Symbols

Strings and symbols are encoded identically. The only difference is the
context in which they are used (symbols cannot be placed on the stack).

```
len string
```

`len` is an `int` representing the number of bytes, `string` is a sequence
of bytes representing a UTF8 string.

## Structs

Constructed at runtime, therefore not stored in bytecode.

## Types

These are only used for struct or function definitions.

- Int: `0x01`
- Float: `0x02`
- Bool: `0x04`
- String: `0x08`
- Func: `0x10 sig` where `sig` is a type signature as specified in `instructions.md`
- Struct: `0x20 i` where `i` is an `int` index in the struct table
