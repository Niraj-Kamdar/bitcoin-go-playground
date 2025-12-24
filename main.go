package main

import (
	"fmt"

	"bitcoin.com/playground/pkg/fields"
)

func main() {
	fe, err := fields.NewFieldElement(7, 13)
	if err != nil {
		panic(err)
	}
	fmt.Println(fe)

	fe2, err := fields.NewFieldElement(12, 13)
	if err != nil {
		panic(err)
	}
	fmt.Println(fe.Eq(fe2))

	fe3, err := fields.NewFieldElement(5, 13)
	if err != nil {
		panic(err)
	}
	fmt.Println(fe.Eq(fe3))
}
