package main

import "fmt"

func main() {
	chip8 := NewChip8("Airplane.ch8")

	//chip8.FDELoop()
	chip8.indexRegister = 0x50
	chip8.registers[2] = 31
	chip8.iDXYN([4]uint16{0xD, 2, 2, 5})

	for y := range chip8.displayHeight {
		for x := range chip8.displayWidth {
			fmt.Print(chip8.display[getDisplayIndex(uint16(x), uint16(y), uint16(chip8.displayWidth), uint16(chip8.displayHeight))])
		}
		fmt.Println()
	}

}
