package main

import (
	"fmt"
	"log"
	"os"

	"ru.prostoyartemka.mppt/stm32/emulator"
)

const INPUT_FILE = "code/out.bin"

func main() {

	data, err := os.ReadFile(INPUT_FILE)
	if err != nil {
		log.Fatal(err)
	}

	context := emulator.LoadEmulator(data)

	for index, reg := range context.Registers {

		fmt.Println("Reg #", index, "Value: ", reg.Value)

	}

}
