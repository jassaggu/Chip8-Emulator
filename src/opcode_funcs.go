package main

import (
	"math/rand/v2"
)

func (c *Chip8) i00E0() {
	for i := range c.display {
		c.display[i] = 0
	}
}

func (c *Chip8) i1NNN(instruction uint16) {
	var address uint16 = instruction & 0x0FFF
	c.programCounter = address
}

func (c *Chip8) i3XNN(instruction uint16) {
	registerIndex := (instruction >> 8) & 0x0F
	value := instruction & 0x00FF

	if c.registers[registerIndex] == uint8(value) {
		c.programCounter += 2
	}
}

func (c *Chip8) i4XNN(instruction uint16) {
	registerIndex := (instruction >> 8) & 0x0F
	value := instruction & 0x00FF

	if c.registers[registerIndex] != uint8(value) {
		c.programCounter += 2
	}
}

func (c *Chip8) i5XY0(instruction uint16) {
	registerIndexX := (instruction >> 8) & 0x0F
	registerIndexY := (instruction >> 4) & 0x00F

	if c.registers[registerIndexX] == c.registers[registerIndexY] {
		c.programCounter += 2
	}
}

func (c *Chip8) i6XNN(instruction uint16) {
	registerIndex := (instruction >> 8) & 0x0F
	value := instruction & 0x00FF

	c.registers[registerIndex] = uint8(value)
}

func (c *Chip8) i7XNN(instruction uint16) {
	registerIndex := (instruction >> 8) & 0x0F
	value := instruction & 0x00FF

	c.registers[registerIndex] += uint8(value)
}

func (c *Chip8) i8XY0(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[2]]
}

func (c *Chip8) i8XY1(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] | c.registers[nibbles[2]]
}

func (c *Chip8) i8XY2(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] & c.registers[nibbles[2]]
}

func (c *Chip8) i8XY3(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] ^ c.registers[nibbles[2]]
}

func (c *Chip8) i8XY4(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] + c.registers[nibbles[2]]
	if uint16(c.registers[nibbles[1]])+uint16(c.registers[nibbles[2]]) > 255 {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}
}

func (c *Chip8) i8XY5(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] - c.registers[nibbles[2]]
	if c.registers[nibbles[1]] > c.registers[nibbles[2]] {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}
}

func (c *Chip8) i8XY6(nibbles [4]uint16) {
	// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
	if c.registers[nibbles[1]]&0x0F == 1 {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}

	c.registers[nibbles[1]] >>= 1
}

func (c *Chip8) i8XY7(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[2]] - c.registers[nibbles[1]]
	if c.registers[nibbles[2]] > c.registers[nibbles[1]] {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}
}

func (c *Chip8) i8XYE(nibbles [4]uint16) {
	// If the least-significant bit of Vx is 1, then VF is set to 1, otherwise 0. Then Vx is divided by 2.
	if c.registers[nibbles[1]]&0x0F == 1 {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}

	c.registers[nibbles[1]] <<= 1
}

func (c *Chip8) i9XY0(instruction uint16) {
	registerIndexX := (instruction >> 8) & 0x0F
	registerIndexY := (instruction >> 4) & 0x00F

	if c.registers[registerIndexX] != c.registers[registerIndexY] {
		c.programCounter += 2
	}
}

func (c *Chip8) iANNN(instruction uint16) {
	var address uint16 = instruction & 0x0FFF
	c.indexRegister = address
}

func (c *Chip8) iBNNN(instruction uint16) {
	var address uint16 = instruction & 0x0FFF
	c.programCounter = address + uint16(c.registers[0])
}

func (c *Chip8) iCXNN(instruction uint16) {
	var randNum uint8 = uint8(rand.IntN(255))
	c.registers[(instruction&0x0F00) >> 8] = uint8(instruction&0x00FF) & randNum
}

func (c *Chip8) iDXYN(nibbles [4]uint16) {
	// Write a sprite to the display, N pixels tall and 1 byte wide

	c.registers[0xF] = 0
	startAddress := c.indexRegister
	height := nibbles[3]
	x, y := c.registers[nibbles[1]], c.registers[nibbles[2]]
	currX, currY := x, y

	// Each row of sprite data
	for currAddress := startAddress; currAddress < startAddress+height; currAddress++ {
		spriteRowBits := getByteBits(c.memory[currAddress])

		for i, bit := range spriteRowBits {
			displayIndex := getDisplayIndex(uint16(currX+uint8(i)), uint16(currY), uint16(c.displayWidth), uint16(c.displayHeight))
			c.display[displayIndex] = bit ^ c.display[displayIndex]

			// Flag VF if collision
			if bit == 1 && c.display[displayIndex] == 0 {
				c.registers[0xF] = 1
			}
		}
		currY++
	}

}

func (c *Chip8) iEX9E(nibbles [4]uint16) {
	key := c.keys[c.registers[nibbles[1]]]
	if key == 1 {
		c.programCounter += 2
	}
}

func (c *Chip8) iEXA1(nibbles [4]uint16) {
	key := c.keys[c.registers[nibbles[1]]]
	if key == 0 {
		c.programCounter += 2
	}
}

func (c *Chip8) iFX07(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.delayTimer
}

func (c *Chip8) iFX15(nibbles [4]uint16) {
	c.delayTimer = c.registers[nibbles[1]]
}

func (c *Chip8) iFX18(nibbles [4]uint16) {
	c.soundTimer = c.registers[nibbles[1]]
}

func (c *Chip8) iFX1E(nibbles [4]uint16) {
	orig := c.indexRegister
	c.indexRegister += nibbles[1]

	// Overflow
	if orig > c.indexRegister {
		c.registers[0xF] = 1
	} else {
		c.registers[0xF] = 0
	}
}

func (c *Chip8) iFX0A(nibbles [4]uint16) {
	pressed := false
	for i, key := range c.keys {
		if key == 1 {
			c.registers[nibbles[1]] = uint8(i)
			pressed = true
		}
	}

	if !pressed {
		c.programCounter -= 2
	}
}

func (c *Chip8) iFX29(nibbles [4]uint16) {
	// The index register I is set to the address of the hexadecimal character in VX
	fontMemoryMap := map[uint16]uint16{
		0x0: 0x50,
		0x1: 0x55,
		0x2: 0x5A,
		0x3: 0x5F,
		0x4: 0x64,
		0x5: 0x69,
		0x6: 0x6E,
		0x7: 0x73,
		0x8: 0x78,
		0x9: 0x7D,
		0xA: 0x82,
		0xB: 0x87,
		0xC: 0x8C,
		0xD: 0x91,
		0xE: 0x96,
		0xF: 0x9B,
	}

	c.indexRegister = fontMemoryMap[nibbles[1]]
}

func (c *Chip8) iFX33(nibbles [4]uint16) {
	num := c.registers[nibbles[1]]
	var digits [3]uint8

	digits[0] = num / 100
	digits[1] = (num / 10) % 10
	digits[2] = num % 10

	c.memory[c.indexRegister] = digits[0]
	c.memory[c.indexRegister+1] = digits[1]
	c.memory[c.indexRegister+2] = digits[2]
}

func (c *Chip8) iFX55(nibbles [4]uint16) {
	for i := range nibbles[1] + 1 {
		c.memory[c.indexRegister+i] = c.registers[i]
	}
}

func (c *Chip8) iFX65(nibbles [4]uint16) {
	for i := range nibbles[1] + 1 {
		c.registers[i] = c.memory[c.indexRegister+i]
	}
}
