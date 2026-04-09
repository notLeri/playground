package main

import "fmt"

type Packet struct {
	Header byte      // 1 байт
	Data   [256]byte // 256 байт
}

const (
	FlagA byte = 1 << 0 // 00000001
	FlagB byte = 1 << 1 // 00000010
)

func main() {
	var status byte

	fmt.Printf("Начальный status: %08b (%d)\n", status, status) // 00000000 (0)

	status |= FlagA
	fmt.Printf("После |= FlagA:  %08b (%d)\n", status, status) // 00000001 (1)

	status &^= FlagA
	fmt.Printf("После &^= FlagA: %08b (%d)\n", status, status) // 00000000 (0)

	status |= FlagB
	fmt.Printf("После |= FlagB:  %08b (%d)\n", status, status) // 00000010 (2)

	status &^= FlagB
	fmt.Printf("После &^= FlagB: %08b (%d)\n", status, status) // 00000000 (0)

	if status&FlagA != 0 {
		fmt.Println("FlagA включён")
	}
}
