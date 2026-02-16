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
		fmt.Println("\t./compiler <path_to_file>")

		return
	}

	path := args[1]

	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	context := emulator.LoadEmulator(data)

	for index, reg := range context.Registers {

		fmt.Println("Reg #", index, "Value: ", reg.Value)

	}

}
