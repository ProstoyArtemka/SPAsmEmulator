package data

type Instruction interface {
	GetArgs() []Argument

	Execute(args [][]byte, context *EmulatorContext)
}

type InstructionSet struct {
	Instructions map[byte]Instruction
}
