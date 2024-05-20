package main

import (
	"fmt"

	zktest "github.com/peitalin/go-zktest"

	"math/big"
	"github.com/consensys/gnark-crypto/ecc/bn254"
    "github.com/consensys/gnark-crypto/ecc/bn254/fp"

)

func TestG1 () {

    x := fp.NewElement(1)
    y := fp.NewElement(2)

    a := bn254.G1Affine{ x, y }
    b := bn254.G1Affine{ x, y }
    c := bn254.G1Affine{ x, y }

    fmt.Println("G1: ", a)

    fmt.Println("\nG1 + G1: ", a.Add(&a, &a))

    fmt.Println("\nG1 * 2: ", b.ScalarMultiplication(&b, big.NewInt(2)))

    fmt.Println("\nG1 double: ", c.Double(&c))

    fmt.Println("\nG1 + G1 ==  G1 * 2", a.Equal(&b))
    fmt.Println("G1 * 2 == Double(G1)", b.Equal(&c))

    _, _, g1Aff, g2Aff := bn254.Generators()
    fmt.Println("\n>>>>>>> G1: ", g1Aff)
    fmt.Println("\n>>>>>>> G2: ", g2Aff)


    // Taking G2 generator points from:
    // https://github.com/ethereum/py_ecc/blob/main/py_ecc/bn128/bn128_curve.py
    u := bn254.E2{ x, y }
    u.SetString("10857046999023057135944570762232829481370756359578518086990519993285655852781", "11559732032986387107991004021392285783925812861821192530917403151452391805634")
    v := bn254.E2{ x, y }
    v.SetString("8495653923123431417604973247489272438418190587263600148770280649306958101930", "4082367875863433681332203403145435568316851327593401208105741076214120093531")

    g2 := bn254.G2Affine{ u, v }
    fmt.Println("\nalternative way to construct g2 >>>>>>>: ", g2)
    fmt.Println("\nNote: g2.String() >>>>>>>: ", g2.String())

}

func TestG2 () {

    _, _, g1, g2 := bn254.Generators()
    g1b := g1 // make a copy of g1 for checking later
    g2b := g2 // make a copy of g2
    fmt.Println(">>>>>>> g1: ", g1.String())
    fmt.Println(">>>>>>> g2: ", g2.String())

    A := g2.ScalarMultiplication(&g2, big.NewInt(5))
    B := g1.ScalarMultiplication(&g1, big.NewInt(6))
    C := g2b.ScalarMultiplication(&g2b, big.NewInt(5 * 6))

    fmt.Println("\nA: ", A.String())
    fmt.Println("\nB: ", B.String())
    fmt.Println("\nC: ", C.String())

    AA := []bn254.G2Affine{*A}
    BB := []bn254.G1Affine{*B}
    CC := []bn254.G2Affine{*C}

    fmt.Println("\nAA: ", AA)
    fmt.Println("\nBB: ", BB)
    fmt.Println("\nCC: ", CC)

    // type GT = fptower.E12
    PairingBBAA, _ := bn254.Pair(BB, AA)
    fmt.Println("\nPair(BB, AA): ", PairingBBAA.String())

    PairingG1CC, _ := bn254.Pair([]bn254.G1Affine{g1b}, CC)
    fmt.Println("\nPair(G1, CC): ", PairingG1CC.String())

    // Check equality on GT
    // pairing(A, B) == pairing(C, G1)
    fmt.Println("\nPairing(BB, AA) == Pairing(G1, CC):", PairingBBAA.Equal(&PairingG1CC))




}

func TestG3() {

    _, _, g1, g2 := bn254.Generators();
    g1_b := g1 // make copies
    g2_b := g2

    a := big.NewInt(4)
    b := big.NewInt(3)
    c := big.NewInt(6)
    d := big.NewInt(2)

    g1.ScalarMultiplication(&g1, a)
    g1.Neg(&g1)
    fmt.Println("\nneg(multiply(G1, a)): ", g1.String())

    g2.ScalarMultiplication(&g2, b)
    fmt.Println("\nmultiply(G2, b): ", g2.String())

    g1_b.ScalarMultiplication(&g1_b, c)
    fmt.Println("\nmultiply(G1, c): ", g1_b.String())

    g2_b.ScalarMultiplication(&g2_b, d)
    fmt.Println("\nmultiply(G2, d): ", g2_b.String())

}


func main() {

    zktest.ZkTest1()
    TestG1()

    zktest.ZkTest2()
    TestG2()

    zktest.ZkTest3()
    TestG3()

    // poly := zktest.Polynomial{
    //     Modulus: 31,
    // }

    // coeffs := []int{4, 5, 6}
    // res := poly.EvalPolyAt(coeffs, 2)
    // fmt.Println("EvalPolyAt: ", res)

}
