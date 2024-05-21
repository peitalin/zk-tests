package main

import (
	"fmt"
	ch7 "github.com/peitalin/go-zktest/ch7"
    ch11 "github.com/peitalin/go-zktest/ch11"
    ch13 "github.com/peitalin/go-zktest/ch13"
	utils "github.com/peitalin/go-zktest"
)

func main() {

    poly1 := utils.Polynomial{
        Modulus: 31,
    }

    coeffs := []int{4, 5, 6}
    res := poly1.EvalPolyAt(coeffs, 2)
    fmt.Println("EvalPolyAt: ", res)


    //// Ch7 Bilinear Pairings
    ch7.Ex1()
    // ch7.Ex2()
    // ch7.Ex3()

    //// Ch13 Encrypted Polynomial Evaluation
    ch13.Ex1()
    // ch13.Ex2()


    //// Ch11 Encrypted Polynomial Evaluation
    ch11.Ex1()

}
