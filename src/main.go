package main

import (
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
	display        [64 * 32]bool

	fontSet [80]uint8
}

func (c *Chip8) pushStack(address uint16) {
	c.stack[c.stackPointer] = address
	c.stackPointer += 1
}

func (c *Chip8) popStack() uint16 {
	c.stackPointer -= 1
	return c.stack[c.stackPointer+1]
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

func (c *Chip8) FDELoop() {
	// ~~ Fetch ~~
	var instruction uint16 = (uint16(c.memory[c.programCounter]) << 8) | uint16(c.memory[c.programCounter+1])
	c.programCounter += 2

	// Use masks to break the instruction into nibbles
	var nibbles [4]uint16
	nibbles[0] = instruction >> 12
	nibbles[1] = (instruction >> 8) & 0b00001111
	nibbles[2] = (instruction >> 4) & 0b000000001111
	nibbles[3] = instruction & 0b0000000000001111

	// ~~ Decode ~~
	switch nibbles[0] {
	case 0x0:
		if instruction == 0x00E0 {
			// Clear screen
			c.i00E0()
		} else if instruction == 0x00EE {
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

	case 0x4:

	case 0x5:

	case 0x6:

	case 0x7:

	case 0x8:

	case 0x9:

	case 0xA:

	case 0xB:

	case 0xC:

	case 0xD:

	case 0xE:

	case 0xF:

	default:
		// optional: handle unknown or invalid nibble
	}
}

func (c *Chip8) i00E0() {
	for i := range c.display {
		c.display[i] = false
	}
}

func (c *Chip8) i1NNN(instruction uint16) {
	var address uint16 = instruction & 0x0FFF
	c.programCounter = address
}

func main() {
	chip8 := NewChip8("Airplane.ch8")

	chip8.FDELoop()

}
