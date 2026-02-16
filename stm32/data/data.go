package data

import (
	"encoding/binary"
)

const (
	TYPE_REGISTER = iota
	TYPE_CONSTANT
	TYPE_REGISTER_SHIFT
	TYPE_MEM_LOAD
	TYPE_LABEL
)

const (
	SUFFIX_EQ = iota
	SUFFIX_NE
	SUFFIX_CS
	SUFFIX_HS
	SUFFIX_CC
	SUFFIX_LO
	SUFFIX_MI
	SUFFIX_PL
	SUFFIX_VS
	SUFFIX_VC
	SUFFIX_HI
	SUFFIX_LS
	SUFFIX_GE
	SUFFIX_LT
	SUFFIX_GT
	SUFFIX_LE
	SUFFIX_AL

	SUFFIX_S
)

const OPERAND_SIZE = 5 // info in description.txt
const LABEL_SIZE = 5
const REGISTER_SIZE = 2 // byte of register index

type EmulatorContext struct {
	Registers []Register

	Err byte
}

func (e *EmulatorContext) GetRegister(index byte) (*Register, bool) {

	if int(index) >= len(e.Registers) {
		return &NIL_REGISTER, true
	}

	return &e.Registers[index], false

}

type Argument struct {
	Size byte
}

type Operand interface {
	GetOperandType() byte

	GetValue(context *EmulatorContext) int32
}

// == OPERANDS ==

type RegisterOperand struct {
	Index byte
}

func (o RegisterOperand) GetOperandType() byte {
	return TYPE_REGISTER
}

func (o RegisterOperand) GetValue(context *EmulatorContext) int32 {

	register, _ := context.GetRegister(o.Index)

	return register.Get()

}

type ConstantOperand struct {
	Value uint32
}

func (o ConstantOperand) GetOperandType() byte {
	return TYPE_CONSTANT
}

func (o ConstantOperand) GetValue(context *EmulatorContext) int32 {

	return int32(o.Value)

}

func ParseInt32(bytes []byte, index int) uint32 {
	return binary.BigEndian.Uint32([]byte{bytes[0+index], bytes[1+index], bytes[2+index], bytes[3+index]})
}

func ParseOperand(bytes []byte) Operand {
	operand_type := bytes[0]

	switch operand_type {

	case TYPE_REGISTER:
		return RegisterOperand{
			bytes[1],
		}

	case TYPE_CONSTANT:

		return ConstantOperand{
			ParseInt32(bytes, 1),
		}
	}

	return nil
}

func Ternary[T any](cond bool, vtrue T, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}
