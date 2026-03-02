package main

import (
	"fmt"
	"log"
	"os"

	"ru.prostoyartemka.mppt/stm32/emulator"
)

func main() {

	args := os.Args

	if len(args) == 1 {

		fmt.Println("Usage: ")
		fmt.Println("\t./compiler emulate <path_to_file>")
		fmt.Println("\t./compiler debug <path_to_file>")

		return
	}

	executeMode := args[1]
	path := args[2]

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	context := emulator.LoadEmulator(data, executeMode == "debug")

	if executeMode == "emulate" {
		for index, reg := range context.Registers {

			fmt.Printf("Reg #%02v\t\t\t0x%016x\t\t%8d\n", index, uint(reg.Get()), reg.Get())

		}

		return
	}

	if executeMode == "debug" {
		emulator.RunDebugger(context)
	}

}
