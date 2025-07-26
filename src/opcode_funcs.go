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
	}
}

func (c *Chip8) i8XY5(nibbles [4]uint16) {
	c.registers[nibbles[1]] = c.registers[nibbles[1]] - c.registers[nibbles[2]]
	if c.registers[nibbles[1]] > c.registers[nibbles[2]] {
		c.registers[0xF] = 1
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
	c.registers[instruction&0x0F00] = uint8(instruction&0x00FF) & randNum
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
			displayIndex := getDisplayIndex(uint16(currX + uint8(i)), uint16(currY), uint16(c.displayWidth), uint16(c.displayHeight))
			c.display[displayIndex] = bit ^ c.display[displayIndex]

			// Flag VF if collision
			if bit == 1 && c.display[displayIndex] == 0 {
				c.registers[0xF] = 1
			}
		}
		currY++
	}

}
