package device

import (
	"testing"
)

func TestConvertFloat64Unsigned(t *testing.T) {
	var float64Value, float64Result float64
	var divider, digitsToRound int
	errorFormat := "The initial '%v' and the loop transformed float64 '%v' are not equal (divider is %v, digits to round %v)."
	// initial
	float64Value = 33
	divider = 1
	digitsToRound = 1
	float64Result = loopConvertFloat64NDigitsAfterCommaUnsigned(float64Value, digitsToRound, divider)
	// should be equal
	if float64Value != float64Result {
		t.Errorf(errorFormat, float64Value, float64Result, divider, digitsToRound)
	}

	// initial
	float64Value = 33.3
	divider = 100
	digitsToRound = 1
	float64Result = loopConvertFloat64NDigitsAfterCommaUnsigned(float64Value, digitsToRound, divider)
	// should be equal
	if float64Value != float64Result {
		t.Errorf(errorFormat, float64Value, float64Result, divider, digitsToRound)
	}
}

func loopConvertFloat64NDigitsAfterCommaUnsigned(float64Value float64, digitsToRound, divider int) float64 {
	var bytes [2]byte
	// integer
	intValue := int(float64Value * float64(divider))
	// to byte
	bytes[0] = byte(intValue)
	bytes[1] = byte(intValue >> 8)
	// back to float
	return prepareInt(getFloatFrom2Bytes(bytes[1], bytes[0]), digitsToRound, divider)
}

func TestConvertFloat64Signed(t *testing.T) {
	var float64Value, float64Result float64
	var divider, digitsToRound int
	errorFormat := "The initial '%v' and the loop transformed float64 '%v' are not equal (divider is %v, digits to round %v)."
	// initial
	float64Value = -33
	divider = 1
	digitsToRound = 1
	float64Result = loopConvertFloat64NDigitsAfterCommaSigned(float64Value, digitsToRound, divider)
	// should be equal
	if float64Value != float64Result {
		t.Errorf(errorFormat, float64Value, float64Result, divider, digitsToRound)
	}

	// initial
	float64Value = -322.3
	divider = 100
	digitsToRound = 1
	float64Result = loopConvertFloat64NDigitsAfterCommaSigned(float64Value, digitsToRound, divider)
	// should be equal
	if float64Value != float64Result {
		t.Errorf(errorFormat, float64Value, float64Result, divider, digitsToRound)
	}
}

func loopConvertFloat64NDigitsAfterCommaSigned(float64Value float64, digitsToRound, divider int) float64 {
	var bytes [2]byte
	// integer
	intValue := int(float64Value * float64(divider))
	// to byte
	bytes[0] = byte(intValue & 0xff)
	bytes[1] = byte(intValue >> 8)
	// back to float
	return prepareInt(getSignedFloatFrom2Bytes(bytes[1], bytes[0]), digitsToRound, divider)
}
