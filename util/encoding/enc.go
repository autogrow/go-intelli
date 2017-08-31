package encoding

import (
	"fmt"
)

// ParseBoolToInt parses a boolean to an integer (true=1, false=0)
func ParseBoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

// ComposeByte sets a collection of bools as bits in a byte
func ComposeByte(pos0 bool, pos1 bool, pos2 bool, pos3 bool, pos4 bool, pos5 bool, pos6 bool, pos7 bool) byte {
	var result byte
	result = SetBit(result, 0, pos0)
	result = SetBit(result, 1, pos1)
	result = SetBit(result, 2, pos2)
	result = SetBit(result, 3, pos3)
	result = SetBit(result, 4, pos4)
	result = SetBit(result, 5, pos5)
	result = SetBit(result, 6, pos6)
	result = SetBit(result, 7, pos7)
	return result
}

// SetBit sets the bit at pos in the integer n.
func SetBit(n byte, pos uint, value bool) byte {
	if value {
		n |= (1 << pos)
		return n
	}
	return n &^ (1 << pos)
}

// ByteToBitString converts a bye to a bit string
func ByteToBitString(b byte) string {
	return fmt.Sprintf("%08b", b)
}

// GetBoolBitValue will return a boolean from the specified bosition in a byte
func GetBoolBitValue(b byte, position int) bool {
	bitString := ByteToBitString(b)
	if bitString[7-position] == '0' {
		return false
	}
	return true
}

// UnsignedIntToBytes converts an unsigned integer to a pair of bytes
func UnsignedIntToBytes(input int) (byte, byte) {
	return byte(input), byte(input >> 8)
}

// SignedIntToBytes converts a signed integer to a pair of bytes
func SignedIntToBytes(input int) (byte, byte) {
	return byte(input & 0xff), byte(input >> 8)
}
