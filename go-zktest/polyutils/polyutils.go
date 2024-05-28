package polyutils

import (
	"fmt"
    "math"
    "errors"
	utils "github.com/peitalin/go-zktest"

	fr "github.com/consensys/gnark-crypto/ecc/bn254/fr"
	p "github.com/consensys/gnark-crypto/ecc/bn254/fr/polynomial"
)


func Ex1() {
	utils.PrintHeading("Polynomial Tools")

	// curveOrder := fr.Modulus()
	// fmt.Println("curveOrder: ", curveOrder)

	p1 := p.Polynomial{
		fr.NewElement(1),
		fr.NewElement(2),
		fr.NewElement(3),
		fr.NewElement(4),
	}
	p2 := p.Polynomial{
		fr.NewElement(5),
		fr.NewElement(6),
		fr.NewElement(7),
		fr.NewElement(8),
	}
	fmt.Println("highest degree polynomial on the bottom.\n")


    apoly := add_polys(p1, p2)
    fmt.Println("\nadd poly:", apoly.Text(10))

    spoly := subtract_polys(p1, p2)
    fmt.Println("\nsubtract poly:", spoly.Text(10))

    mpoly := multiply_polys(p1, p2)
    fmt.Println("\nmultiply poly:", mpoly.Text(10))

    epoly := eval_poly(p1, 3)
    fmt.Println("\nepoly: ", epoly.Text(10), "\n")

    pint := lagrange_interp([]uint64{ 12, 10, 15, 15 })
    fmt.Println("pint: ", pint.Text(10))
    for i := range len(pint) {
        _p := eval_poly(pint, float64(i))
        fmt.Println(fmt.Sprintf("Polynomial Interp: pint(%v): %v", i, _p.Text(10)))
    }

    A := [][]uint64{
        []uint64{ 0, 1, 0, 0, 0, 0 },
        []uint64{ 0, 0, 0, 1, 0, 0 },
        []uint64{ 0, 1, 0, 0, 1, 0 },
        []uint64{ 5, 0, 0, 0, 0, 1 },
    }

    B := [][]uint64{
        []uint64{ 0, 1, 0, 0, 0, 0 },
        []uint64{ 0, 1, 0, 0, 0, 0 },
        []uint64{ 1, 0, 0, 0, 0, 0 },
        []uint64{ 1, 0, 0, 0, 0, 0 },
    }
    C := [][]uint64{
        []uint64{ 0, 0, 0, 1, 0, 0 },
        []uint64{ 0, 0, 0, 0, 1, 0 },
        []uint64{ 0, 0, 0, 0, 0, 1 },
        []uint64{ 0, 0, 1, 0, 0, 0 },
    }

    Ap, Bp, Cp, Z := r1cs_to_qap(A, B, C)
    fmt.Println("\nAp: ")
    for i := range len(Ap) { fmt.Println(Ap[i].Text(10)) }

    fmt.Println("\nBp: ")
    for i := range len(Bp) { fmt.Println(Bp[i].Text(10)) }

    fmt.Println("\nCp: ")
    for i := range len(Cp) { fmt.Println(Cp[i].Text(10)) }

    fmt.Println("\nZ: ", Z)

}



func vecToPolynomial(vec []uint64) p.Polynomial {
    poly := make(p.Polynomial, len(vec))
    for i := range len(vec) {
        poly[i] = fr.NewElement(vec[i])
    }
    return poly
}

func printMatrix(m []p.Polynomial) {
    for i := range len(m) {
        fmt.Println(m[i].Text(10))
    }
}


func transpose(matrix [][]uint64) [][]uint64 {
    return utils.ZipUint64(matrix...)
}

func _sum_polys(a, b p.Polynomial, subtract bool) p.Polynomial {

    var e fr.Element

    poly := make(p.Polynomial, max(len(a), len(b)))
    for i := range len(a) {
        poly[i] = a[i]
    }

    for i := range len(b) {
        if (subtract) {
            poly[i] = *e.Sub(&poly[i], &b[i])
        } else {
            poly[i] = *e.Add(&poly[i], &b[i])
        }
    }
    return poly
}

func add_polys(a, b p.Polynomial) p.Polynomial {
    return _sum_polys(a, b, false)
}

func subtract_polys(a, b p.Polynomial) p.Polynomial {
    return _sum_polys(a, b, true)
}

func multiply_polys(a, b p.Polynomial) p.Polynomial {

    var e fr.Element
    poly := make(p.Polynomial, len(a) + len(b) - 1)

    for i := range len(a) {
        for j := range len(b) {
            term := *e.Mul(&a[i], &b[j])
            poly[i + j] = *e.Add(&poly[i + j], &term)
        }
    }
    return poly
}


func lagrange_interp(vec []uint64) p.Polynomial {

    fr_elements := make([]fr.Element, len(vec))

    for i := range len(vec) {
        fr_elements[i] = fr.NewElement(vec[i])
    }

    return p.InterpolateOnRange(fr_elements)
}


func eval_poly(poly p.Polynomial, x float64) fr.Element {
    // Evaluate polynomial at point x
    var e fr.Element
    out := fr.NewElement(0)

    for i := range len(poly) {
        // out += poly[i] * math.Pow(x, float64(i))
        pow_value := fr.NewElement(uint64(math.Pow(x, float64(i))))
        mvalue := e.Mul(&poly[i], &pow_value)
        out = *e.Add(&out, mvalue)
    }
    return out
}

func r1cs_to_qap(A, B, C [][]uint64) (
    []p.Polynomial,
    []p.Polynomial,
    []p.Polynomial,
    p.Polynomial,
) {
    var e fr.Element

    A = transpose(A)
    B = transpose(B)
    C = transpose(C)

    new_A := []p.Polynomial{}
    new_B := []p.Polynomial{}
    new_C := []p.Polynomial{}

    for i := range len(A) {
        a := lagrange_interp(A[i])
        new_A = append(new_A, a)
    }

    for j := range len(B) {
        b := lagrange_interp(B[j])
        new_B = append(new_B, b)
    }

    for k := range len(C) {
        c := lagrange_interp(C[k])
        new_C = append(new_C, c)
    }

    Z := vecToPolynomial([]uint64{ 1 })

    for l := 1; l <= len(A[0]); l++ {

        ll := fr.NewElement(uint64(l))
        negll := *e.Neg(&ll)
        zPoly := p.Polynomial{ negll, fr.NewElement(uint64(1)) }
        Z = multiply_polys(Z, zPoly)
    }

    return new_A, new_B, new_C, Z
}


func assertZeroFrElement(a fr.Element) (bool, error) {
    if !a.IsZero() {
        errMsg := errors.New(fmt.Sprintf("Term non-zero: %v", a))
        fmt.Println(errMsg)
        return false, errMsg
    } else {
        return true, nil
    }
}