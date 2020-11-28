package helpers

import "math/big"

func Float64ToBigInt(value float64, multiplier float64) *big.Int {
	floatValue := big.NewFloat(value)
	floatValue.Mul(floatValue, big.NewFloat(multiplier))
	intValue := big.NewInt(0)
	floatValue.Int(intValue)

	return intValue
}
