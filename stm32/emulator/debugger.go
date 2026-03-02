package emulator

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/eiannone/keyboard"
	"github.com/fatih/color"
	"ru.prostoyartemka.mppt/stm32/data"
)

func clearConsole() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func printState(index int, memory data.RegistersState) {
	bold := color.New(color.Bold).PrintFunc()

	fmt.Print("#", index, " ")

	fmt.Print("Executed instruction: ")
	bold(memory.Instruction)
	fmt.Print("{", memory.Suffix, "}")

	for _, arg := range memory.Args {
		argStr := data.ArgToString(arg)

		if argStr == "" {
			continue
		}

		fmt.Printf(" %v ", argStr)
	}

	fmt.Println()
	fmt.Println()

	for regI, reg := range memory.Registers {
		fmt.Printf("Reg #%02v\t\t\t0x%016x\t\t%8d\n", regI, uint(reg.Get()), reg.Get())
	}

	for i := range 5 {
		fmt.Print(data.PSR_NAMES[i], " = ", memory.Registers[16].GetBit(i), "\t   ")
	}

	fmt.Println()
}

func RunDebugger(context data.EmulatorContext) {
	if err := keyboard.Open(); err != nil {
		fmt.Println(err)

		return
	}

	index := 0

	clearConsole()
	printState(index, context.Memory[0])

	for {

		_, key, err := keyboard.GetKey()

		clearConsole()

		if err != nil {
			continue
		}

		if key == keyboard.KeyArrowDown {
			index++
		}

		if key == keyboard.KeyArrowUp {
			index--
		}

		if key == keyboard.KeyEsc {
			break
		}

		if index >= len(context.Memory) {
			index = len(context.Memory) - 1
		}

		if index < 0 {
			index = 0
		}

		memory := context.Memory[index]

		printState(index, memory)

	}

	keyboard.Close()
}
