package ch13

import (
	"fmt"
	utils "github.com/peitalin/go-zktest"

	"math/big"
	"github.com/consensys/gnark-crypto/ecc/bn254"
	// "github.com/consensys/gnark-crypto/ecc/bn254/fp"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr"
	"github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
)

func Ex1() {
	utils.PrintHeading("Ch13 Encrypted Polynomial Evaluation - Ex1")

	// Prover
	// x := 5
	fmt.Println("5**3 = ", PowInt(5, 3))

	var X3, X2, X, LHS bn254.G1Affine
	_, _, g1, _ := bn254.Generators()

	X3 = Multiply(g1, PowInt(5, 3))
	X2 = Multiply(g1, PowInt(5, 2))
	X  = Multiply(g1, big.NewInt(5))

	fmt.Println("X3: ", X3.String())
	fmt.Println("X2: ", X2.String())
	fmt.Println("X: ", X.String())

	LHS = Multiply(g1, big.NewInt(39))
	fmt.Println("\nLHS: ", LHS.String())

	//// Polynomial
	// 39 = x**3 - 4x**2 + 3x - 1
	fmt.Println("\nEvaluate: 39 = x**3 - 4x**2 + 3x - 1")
	RHS :=
		Add(
			Add(
				Add(
					Multiply(X3, big.NewInt(1)),
					Multiply(Neg(X2), big.NewInt(4)),
				),
				Multiply(X, big.NewInt(3)),
			),
			Multiply(Neg(g1), big.NewInt(1)),
		)
	fmt.Println("RHS: ", RHS.String())

	fmt.Println("\nLHS == RHS:", LHS.Equal(&RHS))

}

func Ex2() {
	utils.PrintHeading("Ch13 Encrypted Polynomial Evaluation - Ex2")

	curveOrder := fr.Modulus()
	fmt.Println("curveOrder: ", curveOrder)

	a := fr.NewElement(1)
	b := fr.NewElement(2)
	c := fr.NewElement(5)
	p := polynomial.Polynomial{ a, b, c }
	fmt.Println("p.Degree(): ", p.Degree())
	fmt.Println("p: ", p.Text(10))

}

func Neg(a bn254.G1Affine) bn254.G1Affine {
	var g bn254.G1Affine
	return *g.Neg(&a)
}

func Add(a bn254.G1Affine, b bn254.G1Affine) bn254.G1Affine {
	var g bn254.G1Affine
	return *g.Add(&a, &b)
}

func Multiply(a bn254.G1Affine, scalar *big.Int) bn254.G1Affine {
	var g bn254.G1Affine
	return *g.ScalarMultiplication(&a, scalar)
}

func PowInt(base int64, exponent int64) (*big.Int) {
	var y big.Int
	return y.Exp(big.NewInt(base), big.NewInt(exponent), nil)
}

