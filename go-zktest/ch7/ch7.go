package ch7

import (
	"fmt"
	utils "github.com/peitalin/go-zktest"

	"math/big"
	"github.com/consensys/gnark-crypto/ecc/bn254"
    "github.com/consensys/gnark-crypto/ecc/bn254/fp"
)


func Ex1() {
    utils.PrintHeading("Ch7 Bilinear Maps - Ex1")

    x := fp.NewElement(1)
    y := fp.NewElement(2)
    var G1, a, b, c bn254.G1Affine
    G1 = bn254.G1Affine{ x, y }

    fmt.Println(">>> Elliptic Curve Operations:")
    fmt.Println("G1.String(): ", G1.String())
    fmt.Println("\nG1 + G1: ", a.Add(&G1, &G1))

    fmt.Println("\nG1 * 2: ", b.ScalarMultiplication(&G1, big.NewInt(2)).String())

    fmt.Println("\nG1 double: ", c.Double(&G1).String())

    fmt.Println("\nG1 + G1 ==  G1 * 2:", a.Equal(&b))
    fmt.Println("G1 * 2 == Double(G1):", b.Equal(&c))

    _, _, g1Aff, g2Aff := bn254.Generators()
    fmt.Println("\n>>> Generators: ")
    fmt.Println("G1: ", g1Aff.String())
    fmt.Println("G2: ", g2Aff.String())

    fmt.Println("\n>>> Alternative way to construct generator G2 point:")
    // Taking G2 generator points from:
    // https://github.com/ethereum/py_ecc/blob/main/py_ecc/bn128/bn128_curve.py
    var u, v bn254.E2
    u.SetString(
        "10857046999023057135944570762232829481370756359578518086990519993285655852781",
        "11559732032986387107991004021392285783925812861821192530917403151452391805634",
    )
    v.SetString(
        "8495653923123431417604973247489272438418190587263600148770280649306958101930",
        "4082367875863433681332203403145435568316851327593401208105741076214120093531",
    )

    g2 := bn254.G2Affine{ u, v }
    fmt.Println("Constructed g2.String() =>", g2.String())
    fmt.Println("g2Aff == g2 =>", g2.Equal(&g2Aff))

    fmt.Println("\n>>> G1 Jacobian Operations:")
    aJac := bn254.G1Jac{ x, y, fp.NewElement(0) }
    _aJac := bn254.G1Jac{ x, y, fp.NewElement(0) }

    fmt.Println("aJac: ", aJac)
    fmt.Println("aJac - aJac: ", aJac.SubAssign(&_aJac))

}

func Ex2() {
    utils.PrintHeading("\nCh7 Bilinear Maps - Ex2")

    var A, C bn254.G2Affine
    var B bn254.G1Affine
    var PairingBA, PairingG1C bn254.GT // type GT is a F12 point: fptower.E12

    _, _, g1, g2 := bn254.Generators()
    fmt.Println("g1: ", g1.String())
    fmt.Println("g2: ", g2.String())

    A = *A.ScalarMultiplication(&g2, big.NewInt(5))
    B = *B.ScalarMultiplication(&g1, big.NewInt(6))
    C = *C.ScalarMultiplication(&g2, big.NewInt(5 * 6))

    fmt.Println("\nA: g2 * 5 =", A.String())
    fmt.Println("\nB: g1 * 6 =", B.String())
    fmt.Println("\nC: g2 * 5 * 6 =", C.String())

    PairingBA, _ = bn254.Pair([]bn254.G1Affine{B}, []bn254.G2Affine{A})
    fmt.Println("\nPair(B, A): ", PairingBA.String())

    PairingG1C, _ = bn254.Pair([]bn254.G1Affine{g1}, []bn254.G2Affine{C})
    fmt.Println("\nPair(G1, C): ", PairingG1C.String())

    // Check equality on GT
    fmt.Println("\nPairing(B, A) == Pairing(G1, C):", PairingBA.Equal(&PairingG1C))

}

func Ex3() {
    utils.PrintHeading("\nCh7 Bilinear Maps - Ex3")

    var g1_a, g1_c bn254.G1Affine
    var g2_b, g2_d bn254.G2Affine

    _, _, g1, g2 := bn254.Generators();

    a := big.NewInt(4)
    b := big.NewInt(3)
    c := big.NewInt(6)
    d := big.NewInt(2)

    g1_a.ScalarMultiplication(&g1, a).Neg(&g1_a)
    fmt.Println("\nneg(multiply(G1, a)): ", g1_a.String())

    g2_b.ScalarMultiplication(&g2, b)
    fmt.Println("\nmultiply(G2, b): ", g2_b.String())

    g1_c.ScalarMultiplication(&g1, c)
    fmt.Println("\nmultiply(G1, c): ", g1_c.String())

    g2_d.ScalarMultiplication(&g2, d)
    fmt.Println("\nmultiply(G2, d): ", g2_d.String())

    fmt.Println(`
    These outputs are tested with the ecPairing precompile in sol-zktest: address(0x08).staticcall(input)
    where input = uint256[12] array of these elliptic curve points
    `)
}


