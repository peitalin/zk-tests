package ch11

import (
	"fmt"
	utils "github.com/peitalin/go-zktest"
	// "math/big"
	// "github.com/consensys/gnark-crypto/ecc/bn254"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
)

func Ex1() {
	utils.PrintHeading("Ch11 Quadratic Arithmetic Programs - Ex1")

	curveOrder := fr.Modulus()
	fmt.Println("curveOrder: ", curveOrder)

	p1 := polynomial.Polynomial{
		NegCoeff(6),
		fr.NewElement(4),
		fr.NewElement(9),
	}
	p2 := polynomial.Polynomial{
		fr.NewElement(2),
		fr.NewElement(4),
	}
	fmt.Println("p1: ", p1.Text(10))
	fmt.Println("p2: ", p2.Text(10))
	fmt.Println("p1 + p2: ", p1.Add(p1, p2).Text(10))
	fmt.Println("highest degree polynomial on the bottom.\n")


	p3 := []fr.Element{
		fr.NewElement(4),
		fr.NewElement(12),
		fr.NewElement(6),
	}
	// InteroplateOnRange faulty
	pInt := polynomial.InterpolateOnRange(p3)
	fmt.Println("pInt: ", pInt.Text(10))

	e1 := fr.NewElement(1)
	e1 = pInt.Eval(&e1)
	e2 := fr.NewElement(2)
	e2 = pInt.Eval(&e2)
	e3 := fr.NewElement(3)
	e3 = pInt.Eval(&e3)
	fmt.Println("pInt(1): ", e1.String())
	fmt.Println("pInt(2): ", e2.String())
	fmt.Println("pInt(3): ", e3.String())


}

// func poly1(v []fr.Element) polynomial.Polynomial {

// 	p3 := []fr.Element{
// 		fr.NewElement(0),
// 		fr.NewElement(0),
// 		fr.NewElement(1),
// 	}
// }

func NegCoeff(a uint64) fr.Element {
	d := fr.NewElement(a)
	return *d.Neg(&d)
}
