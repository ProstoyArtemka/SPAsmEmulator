package data

type Instruction interface {
	GetArgs() []Argument
	GetName() string

	Execute(args [][]byte, context *EmulatorContext)
}

type InstructionSet struct {
	Instructions map[byte]Instruction
}
