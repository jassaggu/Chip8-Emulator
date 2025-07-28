package main

import (
	"fmt"
	"os"
)

type Chip8 struct {
	memory         [4096]byte
	registers      [16]uint8
	indexRegister  uint16
	programCounter uint16
	stack          [16]uint16
	stackPointer   uint16
	delayTimer     uint8
	soundTimer     uint8
	display        [64 * 32]uint8
	displayWidth   uint8
	displayHeight  uint8
	keys           [16]uint8

	fontSet [80]uint8
}

func (c *Chip8) pushStack(address uint16) {
	c.stack[c.stackPointer] = address
	c.stackPointer += 1
}

func (c *Chip8) popStack() uint16 {
	c.stackPointer -= 1
	return c.stack[c.stackPointer]
}

func (c *Chip8) printDisplay() {
	for y := range c.displayHeight {
		for x := range c.displayWidth {
			fmt.Print(c.display[getDisplayIndex(uint16(x), uint16(y), uint16(c.displayWidth), uint16(c.displayHeight))])
		}
		fmt.Println()
	}
}

func NewChip8(ROMPath string) Chip8 {
	// NewChip8: sets PC, loads fonts, loads program

	chip8 := Chip8{}
	chip8.programCounter = 0x200
	const fontSetStartAddress uint16 = 0x50
	chip8.fontSet = [...]uint8{
		0xF0, 0x90, 0x90, 0x90, 0xF0, // 0
		0x20, 0x60, 0x20, 0x20, 0x70, // 1
		0xF0, 0x10, 0xF0, 0x80, 0xF0, // 2
		0xF0, 0x10, 0xF0, 0x10, 0xF0, // 3
		0x90, 0x90, 0xF0, 0x10, 0x10, // 4
		0xF0, 0x80, 0xF0, 0x10, 0xF0, // 5
		0xF0, 0x80, 0xF0, 0x90, 0xF0, // 6
		0xF0, 0x10, 0x20, 0x40, 0x40, // 7
		0xF0, 0x90, 0xF0, 0x90, 0xF0, // 8
		0xF0, 0x90, 0xF0, 0x10, 0xF0, // 9
		0xF0, 0x90, 0xF0, 0x90, 0x90, // A
		0xE0, 0x90, 0xE0, 0x90, 0xE0, // B
		0xF0, 0x80, 0x80, 0x80, 0xF0, // C
		0xE0, 0x90, 0x90, 0x90, 0xE0, // D
		0xF0, 0x80, 0xF0, 0x80, 0xF0, // E
		0xF0, 0x80, 0xF0, 0x80, 0x80, // F
	}

	// Load font to memory
	for i, fontByte := range chip8.fontSet {
		chip8.memory[fontSetStartAddress+uint16(i)] = fontByte
	}

	chip8.displayWidth = 64
	chip8.displayHeight = 32

	// Load ROM/program into memory
	const ROMStartAddress uint16 = 0x200
	ROMDump, err := os.ReadFile(ROMPath)
	if err != nil {
		panic(err)
	}
	for i, ROMByte := range ROMDump {
		chip8.memory[ROMStartAddress+uint16(i)] = ROMByte
	}

	return chip8
}

func (c *Chip8) FDECycle() {
	// ~~ Fetch ~~
	var instruction uint16 = (uint16(c.memory[c.programCounter]) << 8) | uint16(c.memory[c.programCounter+1])
	c.programCounter += 2

	// Use masks to break the instruction into nibbles
	var nibbles [4]uint16
	nibbles[0] = instruction >> 12
	nibbles[1] = (instruction >> 8) & 0x0F
	nibbles[2] = (instruction >> 4) & 0x00F
	nibbles[3] = instruction & 0x000F

	// ~~ Decode ~~
	switch nibbles[0] {
	case 0x0:
		switch instruction {
		case 0x00E0:
			// Clear screen
			c.i00E0()
		case 0x00EE:
			// Return from subroutine
			c.programCounter = c.popStack()
		}

	case 0x1:
		// Jump
		c.i1NNN(instruction)

	case 0x2:
		// Call subroutine
		// Push current PC to stack, set the PC to NNN
		c.pushStack(c.programCounter)
		c.i1NNN(instruction)

	case 0x3:
		// Skip next instruction if Vx = NN
		c.i3XNN(instruction)

	case 0x4:
		// Skip next instruction if Vx != NN
		c.i4XNN(instruction)

	case 0x5:
		// Skip next instruction if Vx = Vy
		c.i5XY0(instruction)

	case 0x6:
		// Set the register VX to the value NN
		c.i6XNN(instruction)

	case 0x7:
		// Add the value NN to VX
		c.i7XNN(instruction)

	case 0x8:
		switch nibbles[3] {
		case 0x0:
			// VX is set to the value of VYdd
			c.i8XY0(nibbles)

		case 0x1:
			// VX is set to VX OR VY
			c.i8XY1(nibbles)

		case 0x2:
			// VX is set to VX AND VY
			c.i8XY2(nibbles)

		case 0x3:
			// VX is set to VX XOR VY
			c.i8XY3(nibbles)

		case 0x4:
			// Add: VX is set to the value of VX plus the value of VY
			c.i8XY4(nibbles)

		case 0x5:
			// 8XY5 sets VX to the result of VX - VY
			c.i8XY5(nibbles)

		case 0x6:
			// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0, then Vx is divided by 2
			// TODO this may be incorrect
			c.i8XY6(nibbles)

		case 0x7:
			// 8XY7 sets VX to the result of VY - VX
			c.i8XY7(nibbles)

		case 0xE:
			// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0, then Vx is multiplied by 2
			// TODO this may be incorrect
			c.i8XYE(nibbles)
		}

	case 0x9:
		// Skip next instruction if Vx != Vy
		c.i9XY0(instruction)

	case 0xA:
		// Set index register to NNN
		c.iANNN(instruction)

	case 0xB:
		// Jump with offset
		c.iBNNN(instruction)

	case 0xC:
		// Random number, AND with NN, store in VX
		c.iCXNN(instruction)

	case 0xD:
		// Display n-byte sprite starting at memory location I at (Vx, Vy), set VF = collision
		c.iDXYN(nibbles)

	case 0xE:
		switch instruction & 0xFF {
		case 0x9E:
			// If the key in VX is pressed, skip next instruction
			c.iEX9E(nibbles)

		case 0xA1:
			// If the key in VX is not pressed, skip next instruction
			c.iEXA1(nibbles)

		}

	case 0xF:
		// All F instructions start with FX
		switch instruction & 0xFF {
		case 0x07:
			// Set VX to value in the delay timer
			c.iFX07(nibbles)

		case 0x15:
			// Set delay timer to value in VX
			c.iFX07(nibbles)

		case 0x18:
			// Set sound timer to value in VX
			c.iFX18(nibbles)

		case 0x1E:
			// Adds VF to index register
			c.iFX1E(nibbles)

		case 0x0A:
			// Pauses until key is pressed and stored key value in VX
			c.iFX0A(nibbles)

		case 0x29:
			// The index register I is set to the address of the hexadecimal character in VX
			c.iFX29(nibbles)

		case 0x33:
			// Binary-coded decimal conversion (digits stored at index register + ...)
			c.iFX33(nibbles)

		case 0x55:
			// Write values in registers 0-X in memory starting from index register
			c.iFX55(nibbles)

		case 0x65:
			// Write values from memory starting from index register into registers 0-X
			c.iFX65(nibbles)
		}

	default:
		fmt.Print(instruction)
		panic(": invalid instruction")
	}
}
