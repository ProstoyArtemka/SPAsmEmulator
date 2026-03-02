package emulator

import (
	"fmt"
	"slices"

	"ru.prostoyartemka.mppt/stm32/data"
)

const MAX_PROGRAM_SIZE = 2048

var instructionSet = data.STM32

func loadRegisters(context *data.EmulatorContext) {

	// Load empty registers

	for range 13 {
		context.Registers = append(context.Registers, data.Register{})
	}

	// Stack pointer - WIP
	context.Registers = append(context.Registers, data.NIL_REGISTER)

	// Link pointer
	context.Registers = append(context.Registers, data.Register{})

	// Program Counter
	context.Registers = append(context.Registers, data.Register{})

	// Program Status Register
	context.Registers = append(context.Registers, data.Register{})

}

func loadInstructions(instructions []byte, context data.EmulatorContext) {
	// WIP
}

func getProgramCounter(context *data.EmulatorContext) *data.Register {
	return &context.Registers[data.PC]
}

func appendMemory(context *data.EmulatorContext, instruction data.Instruction, suffix byte, args [][]byte) {
	registers := slices.Clone(context.Registers)

	context.Memory = append(context.Memory, data.RegistersState{
		Registers:   registers,
		Instruction: instruction.GetName(),
		Suffix:      data.SUFFIXES_NAMES[suffix],
		Args:        args,
	})
}

func LoadEmulator(instructions []byte, debug bool) data.EmulatorContext {
	var context = data.EmulatorContext{}

	loadRegisters(&context)
	loadInstructions(instructions, context)

	programSize := min(MAX_PROGRAM_SIZE, len(instructions))

	for _ = 0; getProgramCounter(&context).Get() < int32(programSize); {
		currentByte := instructions[getProgramCounter(&context).Get()]

		instruction := instructionSet.Instructions[currentByte]

		if instruction == nil {
			context.Err = 1

			break
		}

		getProgramCounter(&context).Increment()

		var bytesArgs [][]byte = [][]byte{}

		for _, arg := range instruction.GetArgs() {

			var bytedArgument []byte = []byte{}

			for i := 0; i < int(arg.Size); i++ {
				newByte := instructions[int32(i)+getProgramCounter(&context).Get()]

				bytedArgument = append(bytedArgument, newByte)
			}

			getProgramCounter(&context).Add(int32(arg.Size))

			bytesArgs = append(bytesArgs, bytedArgument)
		}

		suffix := bytesArgs[0][0]
		status, _ := context.GetRegister(data.PSR)

		if data.ExecuteSuffix(suffix, status) {
			instruction.Execute(bytesArgs, &context)

			if debug {
				appendMemory(&context, instruction, suffix, bytesArgs)
			}
		}

		if context.Err != 0 {
			fmt.Println("Exit with error", context.Err)

			return context
		}
	}

	return context
}
