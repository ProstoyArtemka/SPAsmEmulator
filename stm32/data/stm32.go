package data

import (
	"math/bits"
)

const (
	OPERATION_ADD = iota
	OPERATION_SUB
)

func executeSuffix(suffix byte, status *Register) bool {

	if suffix == SUFFIX_EQ {
		return status.GetBit(Z) == 1
	}

	if suffix == SUFFIX_NE {
		return status.GetBit(Z) == 0
	}

	if suffix == SUFFIX_CS || suffix == SUFFIX_HS {
		return status.GetBit(C) == 1
	}

	if suffix == SUFFIX_CC || suffix == SUFFIX_LO {
		return status.GetBit(C) == 0
	}

	if suffix == SUFFIX_MI {
		return status.GetBit(N) == 1
	}

	if suffix == SUFFIX_PL {
		return status.GetBit(N) == 0
	}

	if suffix == SUFFIX_VS {
		return status.GetBit(V) == 1
	}

	if suffix == SUFFIX_VC {
		return status.GetBit(V) == 0
	}

	if suffix == SUFFIX_HI {
		return status.GetBit(C) == 1 && status.GetBit(Z) == 0
	}

	if suffix == SUFFIX_LS {
		return status.GetBit(C) == 0 && status.GetBit(Z) == 1
	}

	if suffix == SUFFIX_GE {
		return status.GetBit(N) == status.GetBit(V)
	}

	if suffix == SUFFIX_LT {
		return status.GetBit(N) != status.GetBit(V)
	}

	if suffix == SUFFIX_GT {
		return status.GetBit(Z) == 0 && status.GetBit(N) == status.GetBit(V)
	}

	if suffix == SUFFIX_LE {
		return status.GetBit(Z) == 1 && status.GetBit(N) != status.GetBit(V)
	}

	return true
}

func testOverflow(a int32, b int32, result int32, operation int) bool {

	if operation == OPERATION_ADD {

		if a > 0 && b > 0 && result < 0 {
			return true
		}

		if a < 0 && b < 0 && result >= 0 {
			return true
		}

	}

	if operation == OPERATION_SUB {

		if a >= 0 && b < 0 && result >= 0 {
			return true
		}

		if a < 0 && b >= 0 && result < 0 {
			return true
		}

	}

	return false
}

func updateStatusFull(result int32, carry int32, overflow bool, status *Register) {

	status.SetBit(N, Ternary(result < 0, 1, 0))
	status.SetBit(Z, Ternary(result == 0, 1, 0))
	status.SetBit(C, Ternary(carry != 0, 1, 0))
	status.SetBit(V, Ternary(overflow, 1, 0))

}

func updateStatusNZ(result int32, status *Register) {

	status.SetBit(N, Ternary(result < 0, 1, 0))
	status.SetBit(Z, Ternary(result == 0, 1, 0))

}

var STM32 InstructionSet = InstructionSet{
	Instructions: map[byte]Instruction{
		0x01: AddInstruction{},

		0x02: SubInstruction{},

		0x03: MulInstruction{},

		0x04: BranchInstruction{},

		0x05: MoveInstruction{},

		0x06: SignedDivideInstruction{},
	},
}

type AddInstruction struct{}

func (instruction AddInstruction) GetArgs() []Argument {

	return []Argument{
		{Size: 1},             // Suffix
		{Size: 1},             // Has destination?
		{Size: REGISTER_SIZE}, // Destination register
		{Size: REGISTER_SIZE}, // Register
		{Size: OPERAND_SIZE},  // Operand
	}

}

func (instruction AddInstruction) Execute(args [][]byte, context *EmulatorContext) {

	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	hasDestination := false
	if args[1][0] != 0 {
		hasDestination = true
	}

	destination, err := context.GetRegister(args[2][1])
	if err {
		context.Err = 1

		return
	}

	register, err := context.GetRegister(args[3][1])
	if err {
		context.Err = 1

		return
	}

	operand := ParseOperand(args[4])

	if operand == nil {
		context.Err = 1

		return
	}

	a := register.Get()
	b := operand.GetValue(context)

	uResult, uCarry := bits.Add32(uint32(a), uint32(b), 0)
	result := int32(uResult)

	if suffix == SUFFIX_S {

		updateStatusFull(result, int32(uCarry), testOverflow(a, b, result, OPERATION_ADD), status)

	}

	if hasDestination {
		destination.Set(result)

		return
	}

	register.Set(result)
}

type SubInstruction struct{}

func (instruction SubInstruction) GetArgs() []Argument {

	return []Argument{
		{Size: 1},             // Suffix
		{Size: 1},             // Destination Exists?
		{Size: REGISTER_SIZE}, // Destination register
		{Size: REGISTER_SIZE}, // Register
		{Size: OPERAND_SIZE},  // Operand
	}

}

func (instruction SubInstruction) Execute(args [][]byte, context *EmulatorContext) {

	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	hasDestination := false
	if args[1][0] != 0 {
		hasDestination = true
	}

	destination, err := context.GetRegister(args[2][1])
	if err {
		context.Err = 1

		return
	}

	register, err := context.GetRegister(args[3][1])

	if err {
		context.Err = 1

		return
	}

	operand := ParseOperand(args[4])

	if operand == nil {
		context.Err = 1

		return
	}

	a := register.Get()
	b := operand.GetValue(context)

	uResult, uCarry := bits.Sub32(uint32(a), uint32(b), 0)
	result := int32(uResult)

	if suffix == SUFFIX_S {

		updateStatusFull(result, int32(uCarry), testOverflow(a, b, result, OPERATION_ADD), status)

	}

	if hasDestination {
		destination.Set(int32(result))

		return
	}

	register.Set(int32(result))
}

type MulInstruction struct{}

func (instruction MulInstruction) GetArgs() []Argument {

	return []Argument{
		{Size: 1},             // Suffix
		{Size: 1},             // Destination Exists?
		{Size: REGISTER_SIZE}, // Destination register
		{Size: REGISTER_SIZE}, // Register
		{Size: REGISTER_SIZE}, // Operand
	}

}

func (instruction MulInstruction) Execute(args [][]byte, context *EmulatorContext) {

	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	hasDestination := false
	if args[1][0] != 0 {
		hasDestination = true
	}

	destination, err := context.GetRegister(args[2][1])
	if err {
		context.Err = 1

		return
	}

	register, err := context.GetRegister(args[3][1])

	if err {
		context.Err = 1

		return
	}

	operand := ParseOperand(args[4])

	if operand == nil {
		context.Err = 1

		return
	}

	_, uLow := bits.Mul32(uint32(register.Get()), uint32(operand.GetValue(context)))
	result := int32(uLow)

	if suffix == SUFFIX_S {
		updateStatusNZ(result, status)
	}

	if hasDestination {
		destination.Set(result)

		return
	}

	register.Set(result)
}

type BranchInstruction struct{}

func (instruction BranchInstruction) GetArgs() []Argument {

	return []Argument{
		{Size: 1}, // Suffix
		{Size: LABEL_SIZE},
	}

}

func (instruction BranchInstruction) Execute(args [][]byte, context *EmulatorContext) {

	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	label := ParseInt32(args[1], 1)

	programCounter, _ := context.GetRegister(PC)

	programCounter.Set(int32(label))
}

type MoveInstruction struct{}

func (instruction MoveInstruction) GetArgs() []Argument {
	return []Argument{
		{Size: 1},             // Suffix
		{Size: REGISTER_SIZE}, // Destination register
		{Size: OPERAND_SIZE},  // Const or regster to move
	}
}

func (instruction MoveInstruction) Execute(args [][]byte, context *EmulatorContext) {
	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	registerIndex := args[1][1]
	register, err := context.GetRegister(registerIndex)

	if err {
		context.Err = 1

		return
	}

	operand := ParseOperand(args[2])
	result := operand.GetValue(context)

	register.Set(result)

	if suffix == SUFFIX_S {
		updateStatusNZ(result, status)

		// TODO: Update C when register shifted and gives carry
	}

}

type SignedDivideInstruction struct{}

func (instruction SignedDivideInstruction) GetArgs() []Argument {
	return []Argument{
		{Size: 1},             // Suffix
		{Size: 1},             // Has destination
		{Size: REGISTER_SIZE}, // Destination register
		{Size: REGISTER_SIZE}, // Divisible
		{Size: REGISTER_SIZE}, // Divisor
	}
}

func (instruction SignedDivideInstruction) Execute(args [][]byte, context *EmulatorContext) {
	suffix := args[0][0]
	status, _ := context.GetRegister(PSR)

	if !executeSuffix(suffix, status) {
		return
	}

	hasDestination := false
	if args[1][0] != 0 {
		hasDestination = true
	}

	destination, err := context.GetRegister(args[2][1])
	if err {
		context.Err = 1

		return
	}

	divisible, err := context.GetRegister(args[3][1])
	if err {
		context.Err = 1

		return
	}

	divisor, err := context.GetRegister(args[4][1])
	if err {
		context.Err = 1

		return
	}

	uResult := (divisible.Get()) / (divisor.Get())
	result := int32(uResult)

	if hasDestination {
		destination.Set(result)

		return
	}

	divisible.Set(result)
}
