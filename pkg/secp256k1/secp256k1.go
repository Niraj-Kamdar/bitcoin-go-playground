package secp256k1

// Test-only implementation of low-level secp256k1 field and group arithmetic
// It is designed for ease of understanding, not performance.
// WARNING: This code is slow and trivially vulnerable to side channel attacks. Do not use for
// anything but tests.
// Exports:
// * FE: class for secp256k1 field elements
// * GE: class for secp256k1 group elements
// * G: the secp256k1 generator point

import (
	"errors"
	"math/big"
)

type FE struct {
	Num *big.Int
	Den *big.Int
}

var ErrDenIsZero = errors.New("denominator is zero")
var ErrSqrtFailed = errors.New("sqrt failed")
var ErrNumOutOfFeSize = errors.New("num out of FE_SIZE")
var ErrGeCreationFailed = errors.New("ge creation failed")
var ErrAddFailed = errors.New("addition of ge failed")

func GetFeSize() *big.Int {
	var a = new(big.Int).Exp(new(big.Int).SetInt64(2), new(big.Int).SetInt64(256), nil)
	var b = new(big.Int).Exp(new(big.Int).SetInt64(2), new(big.Int).SetInt64(32), nil)
	var c = new(big.Int).SetInt64(977)

	return new(big.Int).Sub(new(big.Int).Sub(a, b), c)
}

var FE_SIZE = GetFeSize()

func NewFE(_num *big.Int, _den *big.Int) (*FE, error) {
	var num = new(big.Int).Mod(_num, FE_SIZE)
	var den = new(big.Int).Mod(_den, FE_SIZE)

	if den == new(big.Int).SetInt64(0) {
		return nil, ErrDenIsZero
	}

	if num == new(big.Int).SetInt64(0) {
		den = new(big.Int).SetInt64(1)
	}

	return &FE{Num: num, Den: den}, nil
}

func Add(a *FE, b *FE) (*FE, error) {
	var newNum = new(big.Int).Add(new(big.Int).Mul(a.Num, b.Den), new(big.Int).Mul(b.Den, a.Num))
	var newDen = new(big.Int).Mul(a.Den, b.Den)
	return NewFE(newNum, newDen)
}

func Sub(a *FE, b *FE) (*FE, error) {
	var newNum = new(big.Int).Sub(new(big.Int).Mul(a.Num, b.Den), new(big.Int).Mul(b.Den, a.Num))
	var newDen = new(big.Int).Mul(a.Den, b.Den)
	return NewFE(newNum, newDen)
}

func Mul(a *FE, b *FE) (*FE, error) {
	var newNum = new(big.Int).Mul(a.Num, b.Num)
	var newDen = new(big.Int).Mul(a.Den, b.Den)
	return NewFE(newNum, newDen)
}

func ScalarMul(a *FE, b *big.Int) (*FE, error) {
	var newNum = new(big.Int).Mul(a.Num, b)
	var newDen = a.Den
	return NewFE(newNum, newDen)
}

func Pow(a *FE, b *big.Int) (*FE, error) {
	var newNum = new(big.Int).Exp(a.Num, b, FE_SIZE)
	var newDen = new(big.Int).Exp(a.Den, b, FE_SIZE)
	return NewFE(newNum, newDen)
}

func Inv(a *FE) (*FE, error) {
	return Pow(a, new(big.Int).SetInt64(-1))
}

func ToBigInt(a *FE) *big.Int {
	var num *big.Int
	if a.Den != big.NewInt(1) {
		var x = new(big.Int).Exp(a.Den, big.NewInt(-1), FE_SIZE)
		num = new(big.Int).Mul(a.Num, x)
		a.Num = num
		a.Den = big.NewInt(1)
	}
	return num
}

func Sqrt(a *FE) (*FE, error) {
	var v = ToBigInt(a)
	var s = new(big.Int).Exp(v, new(big.Int).Div(new(big.Int).Add(FE_SIZE, big.NewInt(1)), big.NewInt(4)), FE_SIZE)
	var ss = new(big.Int).Exp(s, big.NewInt(2), FE_SIZE)
	if ss == v {
		return NewFE(s, big.NewInt(1))
	}
	return nil, ErrSqrtFailed
}

func HasSqrt(a *FE) bool {
	_, e := Sqrt(a)
	if e == ErrSqrtFailed {
		return false
	}
	return true
}

func IsEven(a *FE) bool {
	return new(big.Int).Mod(ToBigInt(a), big.NewInt(2)) == big.NewInt(0)
}

func Eq(a *FE, b *FE) bool {
	return new(big.Int).Mod(new(big.Int).Sub(new(big.Int).Mul(a.Num, b.Den), new(big.Int).Mul(b.Num, a.Den)), FE_SIZE) == big.NewInt(0)
}

func ToBytes(a *FE) [32]byte {
	var result [32]byte
	num := ToBigInt(a)
	byteSlice := num.Bytes()
	copy(result[32-len(byteSlice):], byteSlice)
	return result
}

func ToHex(a *FE) string {
	var result [64]byte
	num := ToBigInt(a)
	stringSlice := num.Text(16)
	copy(result[64-len(stringSlice):], stringSlice)
	return string(result[:])
}

func FromBytes(a [32]byte) (*FE, error) {
	num := new(big.Int).SetBytes(a[:])
	cmp := num.Cmp(FE_SIZE)
	if cmp == -1 {
		return NewFE(num, big.NewInt(1))
	}
	return nil, ErrNumOutOfFeSize
}

func GetOrder() *big.Int {
	order, err := new(big.Int).SetString("FFFFFFFFFFFFFFFFFFFFFFFFFFFFFFFEBAAEDCE6AF48A03BBFD25E8CD0364141", 16)
	if !err {
		panic("error expected to parse order of elliptic curve!")
	}
	return order
}

var ORDER = GetOrder()
var ORDER_HALF = new(big.Int).Div(ORDER, big.NewInt(2))

type GE struct {
	X   *FE
	Y   *FE
	Inf bool
}

func NewGE(x *FE, y *FE) (*GE, error) {
	fy, err := Pow(y, big.NewInt(2))
	if err != nil {
		return nil, err
	}

	fx, err := Pow(x, big.NewInt(3))
	if err != nil {
		return nil, err
	}

	scalar, err := NewFE(big.NewInt(7), big.NewInt(1))
	if err != nil {
		return nil, err
	}

	add, err := Add(fx, scalar)
	if err != nil {
		return nil, err
	}

	if Eq(add, fy) {
		return &GE{X: fx, Y: fy, Inf: false}, nil
	}

	return nil, ErrGeCreationFailed
}

func NewInfGE() *GE {
	return &GE{X: nil, Y: nil, Inf: true}
}

func Add(a *GE, b *GE) (*GE, error) {
	var lam *FE
	if a.Inf {
		return b, nil
	} else if b.Inf {
		return a, nil
	} else if Eq(a.X, b.X) {
		if !Eq(a.Y, b.Y) {
			add, err := Add(a.Y, b.Y)
			if ToBigInt(add) == big.NewInt(0) {
				return &GE{X: nil, Y: nil, Inf: true}, nil
			} else if err != nil {
				return nil, err
			} else {
				return nil, ErrAddFailed
			}
		} else {
			// For Identical input use tangent
			pow1, err := Pow(a.X, big.NewInt(2))
			if err != nil {
				return nil, err
			}

			num1, err := ScalarMul(pow1, big.NewInt(3))
			if err != nil {
				return nil, err
			}
			den1, err := ScalarMul(a.Y, big.NewInt(2))
			if err != nil {
				return nil, err
			}
			num2 := ToBigInt(num1)
			den2 := ToBigInt(den1)
			fe, err := NewFE(num2, den2)
			if err != nil {
				return nil, err
			}
			lam = fe
		}
	} else {
		num1, err := Sub(a.Y, b.Y)
		if err != nil {
			return nil, err
		}

		den1, err := Sub(a.X, b.X)
		if err != nil {
			return nil, err
		}

		num2 := ToBigInt(num1)
		den2 := ToBigInt(den1)

		fe, err := NewFE(num2, den2)
		if err != nil {
			return nil, err
		}
		lam = fe

	}

	f1, err := Pow(lam, big.NewInt(2))
	if err != nil {
		return nil, err
	}
	f2, err := Add(a.X, b.X)
	if err != nil {
		return nil, err
	}
	x, err := Sub(f1, f2)
	if err != nil {
		return nil, err
	}

	f3, err := Sub(a.X, x)
	if err != nil {
		return nil, err
	}

	f4, err := Mul(lam, f3)
	if err != nil {
		return nil, err
	}

	y, err := Sub(f4, a.Y)
	if err != nil {
		return nil, err
	}
	return NewGE(x, y)
}
