package fields

import (
	"bitcoin.com/playground/pkg/assert"
	"math/big"
	"testing"
)

func TestNewFieldElement(t *testing.T) {
	fe, err := NewFieldElement(7, 13)

	assert.Equal(t, err, nil)
	assert.Equal(t, fe.num, 7)
	assert.Equal(t, fe.prime, 13)
}

func FuzzNewFieldElement(f *testing.F) {
	f.Fuzz(func(t *testing.T, num int, prime int) {
		if num >= prime || num < 0 {
			t.Skip()
		}
		if !big.NewInt(int64(prime)).ProbablyPrime(0) {
			t.Skip()
		}
		fe, err := NewFieldElement(num, prime)
		assert.Equal(t, err, nil)
		assert.Equal(t, fe.num, num)
		assert.Equal(t, fe.prime, prime)
	})
}

func TestNewFieldElementInvalidNum(t *testing.T) {
	fe, err := NewFieldElement(14, 13)

	assert.Equal(t, err, ErrNumOutOfFieldRange)
	assert.Equal(t, fe, nil)
}

func FuzzNewFieldElementInvalidNum(f *testing.F) {
	f.Fuzz(func(t *testing.T, num int, prime int) {
		if num < prime && num >= 0 {
			t.Skip()
		}
		if !big.NewInt(int64(prime)).ProbablyPrime(0) {
			t.Skip()
		}
		fe, err := NewFieldElement(num, prime)
		assert.Equal(t, err, ErrNumOutOfFieldRange)
		assert.Equal(t, fe, nil)
	})
}

func TestEqValid(t *testing.T) {
	fe1, err := NewFieldElement(7, 13)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe1.num, 7)
	assert.Equal(t, fe1.prime, 13)

	fe2, err := NewFieldElement(7, 13)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe2.num, 7)
	assert.Equal(t, fe2.prime, 13)
	assert.Equal(t, fe1.Eq(fe2), true)
}

func TestEqInvalidPrime(t *testing.T) {
	fe1, err := NewFieldElement(7, 13)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe1.num, 7)
	assert.Equal(t, fe1.prime, 13)

	fe2, err := NewFieldElement(7, 17)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe2.num, 7)
	assert.Equal(t, fe2.prime, 17)
	assert.Equal(t, fe1.Eq(fe2), false)
}

func TestEqInvalidNum(t *testing.T) {
	fe1, err := NewFieldElement(7, 13)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe1.num, 7)
	assert.Equal(t, fe1.prime, 13)

	fe2, err := NewFieldElement(12, 13)
	assert.Equal(t, err, nil)
	assert.Equal(t, fe2.num, 12)
	assert.Equal(t, fe2.prime, 13)
	assert.Equal(t, fe1.Eq(fe2), false)
}

func TestAdd(t *testing.T) {
	fe1, err := NewFieldElement(44, 57)
	assert.Equal(t, err, nil)

	fe2, err := NewFieldElement(33, 57)
	assert.Equal(t, err, nil)

	// Invariant 1: fe1.prime + fe2.prime == sum.prime
	sum, err := fe1.Add(fe2)
	assert.Equal(t, err, nil)
	assert.Equal(t, sum.prime, fe1.prime)
	assert.Equal(t, sum.prime, fe2.prime)

	// Invariant 2: (fe1.num + fe2.num) % fe1.prime == sum.num
	assert.Equal(t, sum.num, (fe1.num+fe2.num)%fe1.prime)
	assert.Equal(t, sum.num, (fe1.num+fe2.num)%fe2.prime)

	// Invariant 3: sum.num < fe1.prime
	assert.Less(t, sum.num, fe1.prime)
	assert.Less(t, sum.num, fe2.prime)

	// Invariant 4: sum.num >= 0
	assert.GreaterOrEqual(t, sum.num, 0)

	// Invariant 5: fe1 + fe2 == fe2 + fe1
	sum1, err := fe1.Add(fe2)
	assert.Equal(t, err, nil)
	assert.Equal(t, sum1.num, sum.num)

	sum2, err := fe2.Add(fe1)
	assert.Equal(t, err, nil)
	assert.Equal(t, sum2.num, sum.num)

	// Invariant 6: (fe1 + fe2) + fe3 == fe1 + (fe2 + fe3)
	fe3, err := NewFieldElement(12, 57)
	assert.Equal(t, err, nil)

	intermediate, err := fe1.Add(fe2)
	assert.Equal(t, err, nil)

	test1, err := intermediate.Add(fe3)
	assert.Equal(t, err, nil)

	intermediate, err = fe2.Add(fe3)
	assert.Equal(t, err, nil)

	test2, err := fe1.Add(intermediate)
	assert.Equal(t, err, nil)

	assert.Equal(t, test1.Eq(test2), true)
}

func TestSub(t *testing.T) {
	fe1, err := NewFieldElement(44, 57)
	assert.Equal(t, err, nil)

	fe2, err := NewFieldElement(33, 57)
	assert.Equal(t, err, nil)

	// Invariant 1: fe1.prime - fe2.prime == sub.prime
	sub, err := fe1.Sub(fe2)
	assert.Equal(t, err, nil)
	assert.Equal(t, sub.prime, fe1.prime)
	assert.Equal(t, sub.prime, fe2.prime)

	// Invariant 2: (fe1.num - fe2.num) % fe1.prime == sub.num
	assert.Equal(t, sub.num, (fe1.num-fe2.num+fe1.prime)%fe1.prime)
	assert.Equal(t, sub.num, (fe1.num-fe2.num+fe2.prime)%fe2.prime)

	// Invariant 3: sub.num < fe1.prime
	assert.Less(t, sub.num, fe1.prime)
	assert.Less(t, sub.num, fe2.prime)

	// Invariant 4: sub.num >= 0
	assert.GreaterOrEqual(t, sub.num, 0)
}

func TestMul(t *testing.T) {
	fe1, err := NewFieldElement(44, 57)
	assert.Equal(t, err, nil)

	fe2, err := NewFieldElement(33, 57)
	assert.Equal(t, err, nil)

	// Invariant 1: fe1.prime * fe2.prime == mul.prime
	mul, err := fe1.Mul(fe2)
	assert.Equal(t, err, nil)
	assert.Equal(t, mul.prime, fe1.prime)
	assert.Equal(t, mul.prime, fe2.prime)

	// Invariant 2: (fe1.num * fe2.num) % fe1.prime == mul.num
	assert.Equal(t, mul.num, (fe1.num*fe2.num)%fe1.prime)
	assert.Equal(t, mul.num, (fe1.num*fe2.num)%fe2.prime)

	// Invariant 3: mul.num < fe1.prime
	assert.Less(t, mul.num, fe1.prime)
	assert.Less(t, mul.num, fe2.prime)

	// Invariant 4: mul.num >= 0
	assert.GreaterOrEqual(t, mul.num, 0)

	// Invariant 5: fe1 * fe2 == fe2 * fe1
	mul1, err := fe1.Mul(fe2)
	assert.Equal(t, err, nil)
	assert.Equal(t, mul1.num, mul.num)

	mul2, err := fe2.Mul(fe1)
	assert.Equal(t, err, nil)
	assert.Equal(t, mul2.num, mul.num)

	// Invariant 6: (fe1 * fe2) * fe3 == fe1 * (fe2 * fe3)
	fe3, err := NewFieldElement(12, 57)
	assert.Equal(t, err, nil)

	intermediate, err := fe1.Mul(fe2)
	assert.Equal(t, err, nil)

	test1, err := intermediate.Mul(fe3)
	assert.Equal(t, err, nil)

	intermediate, err = fe2.Mul(fe3)
	assert.Equal(t, err, nil)

	test2, err := fe1.Mul(intermediate)
	assert.Equal(t, err, nil)

	assert.Equal(t, test1.Eq(test2), true)
}

func TestPow(t *testing.T) {
	fe1, err := NewFieldElement(44, 57)
	assert.Equal(t, err, nil)

	scalar := 3

	// Invariant 1: fe1.prime ^ scalar == pow.prime
	pow, err := fe1.Pow(scalar)
	assert.Equal(t, err, nil)
	assert.Equal(t, pow.prime, fe1.prime)

	// Invariant 2: (fe1.num ^ scalar) % fe1.prime == pow.num
	assert.Equal(t, pow.num, (fe1.num^scalar)%fe1.prime)

	// Invariant 3: pow.num < fe1.prime
	assert.Less(t, pow.num, fe1.prime)

	// Invariant 4: pow.num >= 0
	assert.GreaterOrEqual(t, pow.num, 0)
}
