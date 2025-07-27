package main

func main() {
	ROMName := "Breakout.ch8"

	chip8 := NewChip8(ROMName)
	chip8.raylibFrontendLoop(ROMName, 512)

}
