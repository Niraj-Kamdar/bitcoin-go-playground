package fields

import (
	"errors"
	"fmt"
)

type FieldElement struct {
	num   int
	prime int
}

var ErrNumOutOfFieldRange = errors.New("num not in field range")
var ErrPrimesMustBeSame = errors.New("primes must be the same")
var ErrDivisionByZero = errors.New("division by zero")

func NewFieldElement(num int, prime int) (*FieldElement, error) {
	if num >= prime || num < 0 {
		return nil, ErrNumOutOfFieldRange
	}

	return &FieldElement{num: num, prime: prime}, nil
}

func (fe *FieldElement) Num() int {
	return fe.num
}

func (fe *FieldElement) Prime() int {
	return fe.prime
}

func (fe *FieldElement) Eq(other *FieldElement) bool {
	return fe.num == other.num && fe.prime == other.prime
}

// Create representation of the field element
func (fe *FieldElement) String() string {
	return fmt.Sprintf("FieldElement(%d %% %d)", fe.num, fe.prime)
}

func (fe *FieldElement) Add(other *FieldElement) (*FieldElement, error) {
	if fe.prime != other.prime {
		return nil, ErrPrimesMustBeSame
	}

	sum := (fe.num + other.num) % fe.prime
	return NewFieldElement(sum, fe.prime)
}

func (fe *FieldElement) Sub(other *FieldElement) (*FieldElement, error) {
	if fe.prime != other.prime {
		return nil, ErrPrimesMustBeSame
	}

	sum := (fe.num - other.num + fe.prime) % fe.prime
	return NewFieldElement(sum, fe.prime)
}

func (fe *FieldElement) Mul(other *FieldElement) (*FieldElement, error) {
	if fe.prime != other.prime {
		return nil, ErrPrimesMustBeSame
	}

	sum := (fe.num * other.num) % fe.prime
	return NewFieldElement(sum, fe.prime)
}

// Fermat's little theorem Inverse n ^ (p - 1) = 1 , n > 0
func (fe *FieldElement) Pow(exp int) (*FieldElement, error) {
	n := exp % (fe.prime - 1)
	num := (fe.num ^ n) % fe.prime
	return NewFieldElement(num, fe.prime)
}

// Fermat's little theorem Inverse n ^ (p - 1) = 1 , n > 0
func (fe *FieldElement) Inv() (*FieldElement, error) {
	return fe.Pow(fe.prime - 2)
}

func (fe *FieldElement) Div(other *FieldElement) (*FieldElement, error) {
	if fe.prime != other.prime {
		return nil, ErrPrimesMustBeSame
	}

	if other.num == 0 {
		return nil, ErrDivisionByZero
	}

	// fe.num * (other.num ^ -1) % fe.prime
	inv, err := other.Inv()
	if err != nil {
		return nil, err
	}

	return fe.Mul(inv)
}
