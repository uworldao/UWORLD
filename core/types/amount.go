// Copyright 2017-2018 The qitmeer developers
// Copyright 2015 The Decred developers
// Copyright 2013, 2014 The btcsuite developers
// Use of this source code is governed by an ISC
// license that can be found in the LICENSE file.

package types

import (
	"errors"
	"math"
)

const (
	// AtomsPerCent is the number of atomic units in one coin cent.
	// TODO, relocate the coin related item to chain's params
	AtomsPerCent = 1e6

	// AtomsPerCoin is the number of atomic units in one coin.
	// TODO, relocate the coin related item to chain's params
	AtomsPerCoin = 1e8

	// MaxAmount is the maximum transaction amount allowed in atoms.
	// TODO, relocate the coin related item to chain's params
	MaxAmount = 21e6 * AtomsPerCoin
)

// AmountUnit describes a method of converting an Amount to something
// other than the base unit of a coin.  The value of the AmountUnit
// is the exponent component of the decadic multiple to convert from
// an amount in coins to an amount counted in atomic units.
type AmountUnit int

// These constants define various units used when describing a coin
// monetary amount.
const (
	AmountMegaCoin  AmountUnit = 6
	AmountKiloCoin  AmountUnit = 3
	AmountCoin      AmountUnit = 0
	AmountMilliCoin AmountUnit = -3
	AmountMicroCoin AmountUnit = -6
	AmountAtom      AmountUnit = -8
)

// Amount represents the base coin monetary unit (colloquially referred
// to as an `Atom').  A single Amount is equal to 1e-8 of a coin.
type Amount int64

// round converts a floating point number, which may or may not be representable
// as an integer, to the Amount integer type by rounding to the nearest integer.
// This is performed by adding or subtracting 0.5 depending on the sign, and
// relying on integer truncation to round the value to the nearest Amount.
func round(f float64) uint64 {
	if f < 0 {
		return uint64(f - 0.5)
	}
	return uint64(f + 0.5)
}

// NewAmount creates an Amount from a floating point value representing
// some value in the currency.  NewAmount errors if f is NaN or +-Infinity,
// but does not check that the amount is within the total amount of coins
// producible as f may not refer to an amount at a single moment in time.
//
// NewAmount is for specifically for converting qitmeer to Atoms (atomic units).
// For creating a new Amount with an int64 value which denotes a quantity of
// Atoms, do a simple type conversion from type int64 to Amount.
func NewAmount(f float64) (uint64, error) {
	// The amount is only considered invalid if it cannot be represented
	// as an integer type.  This may happen if f is NaN or +-Infinity.
	switch {
	case math.IsNaN(f):
		fallthrough
	case math.IsInf(f, 1):
		fallthrough
	case math.IsInf(f, -1):
		return 0, errors.New("invalid coin amount")
	}

	return round(f * AtomsPerCoin), nil
}

// ToUnit converts a monetary amount counted in coin base units to a
// floating point value representing an amount of coins.
func (a Amount) ToUnit(u AmountUnit) float64 {
	return float64(a) / math.Pow10(int(u+8))
}

// ToCoin is the equivalent of calling ToUnit with AmountCoin.
func (a Amount) ToCoin() float64 {
	return a.ToUnit(AmountCoin)
}
