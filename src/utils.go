package main

func getByteBits(b byte) []uint8 {
	var bits []uint8
	for i := 7; i >= 0; i-- {
		bits = append(bits, (b>>i)&1)
	}
	return bits
}

func getDisplayIndex(x uint16, y uint16, width uint16, height uint16) uint16 {
	return uint16((y%height)*width + (x % width))
}
