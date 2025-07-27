package main

import rl "github.com/gen2brain/raylib-go/raylib"

func updateKeys(c *Chip8) {
	if rl.IsKeyDown(rl.KeyOne) {
		c.keys[1] = 1
	} else {
		c.keys[1] = 0
	}

	if rl.IsKeyDown(rl.KeyTwo) {
		c.keys[2] = 1
	} else {
		c.keys[2] = 0
	}

	if rl.IsKeyDown(rl.KeyThree) {
		c.keys[3] = 1
	} else {
		c.keys[3] = 0
	}

	if rl.IsKeyDown(rl.KeyFour) {
		c.keys[0xC] = 1
	} else {
		c.keys[0xC] = 0
	}

	if rl.IsKeyDown(rl.KeyQ) {
		c.keys[4] = 1
	} else {
		c.keys[4] = 0
	}

	if rl.IsKeyDown(rl.KeyW) {
		c.keys[5] = 1
	} else {
		c.keys[5] = 0
	}

	if rl.IsKeyDown(rl.KeyE) {
		c.keys[6] = 1
	} else {
		c.keys[6] = 0
	}

	if rl.IsKeyDown(rl.KeyR) {
		c.keys[0xD] = 1
	} else {
		c.keys[0xD] = 0
	}

	if rl.IsKeyDown(rl.KeyA) {
		c.keys[7] = 1
	} else {
		c.keys[7] = 0
	}

	if rl.IsKeyDown(rl.KeyS) {
		c.keys[8] = 1
	} else {
		c.keys[8] = 0
	}

	if rl.IsKeyDown(rl.KeyD) {
		c.keys[9] = 1
	} else {
		c.keys[9] = 0
	}

	if rl.IsKeyDown(rl.KeyF) {
		c.keys[0xE] = 1
	} else {
		c.keys[0xE] = 0
	}

	if rl.IsKeyDown(rl.KeyZ) {
		c.keys[0xA] = 1
	} else {
		c.keys[0xA] = 0
	}

	if rl.IsKeyDown(rl.KeyX) {
		c.keys[0] = 1
	} else {
		c.keys[0] = 0
	}

	if rl.IsKeyDown(rl.KeyC) {
		c.keys[0xB] = 1
	} else {
		c.keys[0xB] = 0
	}

	if rl.IsKeyDown(rl.KeyV) {
		c.keys[0xF] = 1
	} else {
		c.keys[0xF] = 0
	}
}

func renderDisplay(c *Chip8, pixelSize int) {
	backgroundColour := rl.Black
	pixelColour := rl.White

	xPos := 0
	yPos := 0
	rl.ClearBackground(backgroundColour)
	for y := range c.displayHeight {
		for x := range c.displayWidth {
			pixel := rl.Rectangle{X: float32(xPos), Y: float32(yPos), Width: float32(pixelSize), Height: float32(pixelSize)}
			rl.BeginDrawing()
			if c.display[getDisplayIndex(uint16(x), uint16(y), 64, 32)] == 1 {
				rl.DrawRectangleRec(pixel, pixelColour)
			}
			xPos += pixelSize
		}
		yPos += pixelSize
		xPos = 0
	}
	rl.EndDrawing()
}

func (c *Chip8) raylibFrontendLoop(ROMName string, screenWidth int) {
	screenHeight := screenWidth / 2
	pixelSize := screenWidth / int(c.displayWidth)

	rl.InitWindow(int32(screenWidth), int32(screenHeight), "Chip8 Emulator - running "+ROMName)
	rl.SetTargetFPS(60)

	for !rl.WindowShouldClose() {
		updateKeys(c)
		c.FDECycle()
		renderDisplay(c, pixelSize)
	}
	rl.CloseWindow()

}
