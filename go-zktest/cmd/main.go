package main

import (
	"fmt"

	"github.com/arnaucube/go-snark/bn128"
	zktest "github.com/peitalin/go-zktest"
)

func main() {
    zktest.ZkTest1()

    poly := zktest.Polynomial{
        Modulus: 31,
    }

    coeffs := []int{4, 5, 6}
    res := poly.EvalPolyAt(coeffs, 2)
    fmt.Println("EvalPolyAt: ", res)

    b, err := bn128.NewBn128()
    if err == nil {
        fmt.Println("Bn128: ", b)
    }

    g1 := bn128.NewG1()
    fmt.Println("G1: ", g1)
    // fmt.Println("G2: ", bn128.G2)
}
