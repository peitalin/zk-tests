package main

import (
	utils "github.com/peitalin/go-zktest"
	blsSignatures "github.com/peitalin/go-zktest/blsSignatures"
	ch13 "github.com/peitalin/go-zktest/ch13"
	ch7 "github.com/peitalin/go-zktest/ch7"
	polyutils "github.com/peitalin/go-zktest/polyutils"
	vbuterin "github.com/peitalin/go-zktest/vbuterin"
)

func main() {

	//// Ch7 Bilinear Pairings
	ch7.Ex1()
	ch7.Ex2()
	ch7.Ex3()

	//// Ch13 Encrypted Polynomial Evaluation
	ch13.Ex1()
	ch13.Ex2()

	//// Vbuterin Examples
	vbuterin.Ex1()

	//// Gnark Library Examples
	polyutils.Ex1()
	utils.TestZipImplementations()

	blsSignatures.Ex1()
}
