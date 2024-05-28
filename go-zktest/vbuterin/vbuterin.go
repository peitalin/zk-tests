package vbuterin

import (
	"fmt"
    "math"
    "errors"
	utils "github.com/peitalin/go-zktest"
)

func Ex1() {
	utils.PrintHeading("Vitalik's R1CS to QAP Exercise")

	arr1 := []float64{1,2,3,4}
	arr2 := []float64{5,6,7,8}

    apoly := add_polys(arr1, arr2)
    fmt.Println("\nadd poly:", apoly)

    spoly := subtract_polys(arr1, arr2)
    fmt.Println("\nsubtract poly:", spoly)

    mpoly := multiply_polys(arr1, arr2)
    fmt.Println("\nmultiply poly:", mpoly, "\n")

    mk := mk_singleton(1, 2, 3)
    fmt.Println("\nmk: ", mk)

    epoly := eval_poly([]float64{1,2,3,4}, 3)
    fmt.Println("\nepoly: ", epoly, "\n")

    ipoly, _ := lagrange_interp([]float64{12, 10, 15, 15})
    fmt.Println("\nInterpolated Poly: ", ipoly)

    var pint float64
    for i := range 5 {
        pint = eval_poly(ipoly, float64(i))
        fmt.Println(fmt.Sprintf("eval_poly(%v): %v", i, pint))
    }

    dpoly, remainder := div_polys([]float64{4,5,6,7}, []float64{1,2,3,4})
    fmt.Println("dpoly: ", dpoly, remainder)
    fmt.Println("\n")


    r := []float64{ 1, 3, 35, 9, 27, 30 }

    A := [][]float64{
        []float64{ 0, 1, 0, 0, 0, 0 },
        []float64{ 0, 0, 0, 1, 0, 0 },
        []float64{ 0, 1, 0, 0, 1, 0 },
        []float64{ 5, 0, 0, 0, 0, 1 },
    }
    B := [][]float64{
        []float64{ 0, 1, 0, 0, 0, 0 },
        []float64{ 0, 1, 0, 0, 0, 0 },
        []float64{ 1, 0, 0, 0, 0, 0 },
        []float64{ 1, 0, 0, 0, 0, 0 },
    }
    C := [][]float64{
        []float64{ 0, 0, 0, 1, 0, 0 },
        []float64{ 0, 0, 0, 0, 1, 0 },
        []float64{ 0, 0, 0, 0, 0, 1 },
        []float64{ 0, 0, 1, 0, 0, 0 },
    }

    Ap, Bp, Cp, Z := r1cs_to_qap(A, B, C)
    fmt.Println("\nAp: ")
    for i := range len(Ap) { fmt.Println(Ap[i]) }

    fmt.Println("\nBp: ")
    for i := range len(Bp) { fmt.Println(Bp[i]) }

    fmt.Println("\nCp: ")
    for i := range len(Cp) { fmt.Println(Cp[i]) }

    fmt.Println("\nZ: ", Z)

    Apoly, Bpoly, Cpoly, sol := create_solution_polynomials(r, Ap, Bp, Cp)

    fmt.Println("\nApoly: ", Apoly)
    fmt.Println("\nBpoly: ", Bpoly)
    fmt.Println("\nCpoly: ", Cpoly)
    fmt.Println("\nSolution: ", sol)

    quot := create_divisor_polynomial(sol, Z)
    fmt.Println("\nQuotient: ", quot)

}





func transpose(matrix [][]float64) [][]float64 {
    return utils.Zip2(matrix...)
}

func _sum_polys(a, b []float64, subtract bool) []float64 {
    poly := make([]float64, max(len(a), len(b)))

    for i := range len(a) {
        poly[i] += a[i]
    }

    for j := range len(b) {
        if (subtract) {
            poly[j] -= b[j]
        } else {
            poly[j] += b[j]
        }
    }
    return poly
}

func add_polys(a, b []float64) []float64 {
    return _sum_polys(a, b, false)
}

func subtract_polys(a, b []float64) []float64 {
    return _sum_polys(a, b, true)
}

func multiply_polys(a, b []float64) []float64 {
    poly := make([]float64, len(a) + len(b) - 1)

    for i := range len(a) {
        for j := range len(b) {
            poly[i + j] += a[i] * b[j]
        }
    }
    return poly
}

func mk_singleton(point_loc, height, total_pts float64) []float64 {
    // Make a polynomial which is zero at {1, 2 ... total_pts}, except
    // for `point_loc` where the value is `height`
    fac := 1.0

    for i := 1.0; i <= total_pts; i++ {
        if i != point_loc {
            fac *= point_loc - i
        }
    }

    poly := []float64{ height * 1.0 / fac }

    for i := 1.0; i <= total_pts; i++ {
        if i != point_loc {
            neg_poly := []float64{ -i, 1 }
            poly = multiply_polys(poly, neg_poly)
        }
    }
    return poly
}

func lagrange_interp(vec []float64) ([]float64, error) {
    // Assumes vec[0] = p(1), vec[1] = p(2), etc, tries to find p,
    // expresses result as [deg 0 coeff, deg 1 coeff...]
    poly := make([]float64, len(vec))

    for i := range len(vec) {
        poly = add_polys(poly, mk_singleton(float64(i) + 1, vec[i], float64(len(vec))))
    }
    for j := range len(vec) {
        epoly := eval_poly(poly, float64(j + 1)) - vec[j]
        isZero, errMsg := assertZeroF64(epoly)
        if !isZero {
            return nil, errMsg
        }
    }

    return poly, nil
}

func div_polys(a, b []float64) ([]float64, []float64){
    // Divide a/b, return quotient and remainder
    poly := make([]float64, len(a) - len(b) + 1)
    remainder := a
    isDivisible := func() bool { return len(remainder) >= len(b) }

    for ok := isDivisible(); ok; ok = isDivisible() {

        leading_fac := remainder[len(remainder)-1] / b[len(b)-1]
        pos := len(remainder) - len(b)
        poly[pos] = leading_fac

        mpoly := make([]float64, pos)
        mpoly = append(mpoly, leading_fac)
        mpoly = multiply_polys(b, mpoly)

        spoly := subtract_polys(remainder, mpoly)
        remainder = spoly[:len(spoly) - 1]
    }

    return poly, remainder
}


func eval_poly(poly []float64, x float64) float64 {
    // Evaluate polynomial at point x
    var out float64
    for i := range len(poly) {
        out += poly[i] * math.Pow(x, float64(i))
    }
    return out
}

func r1cs_to_qap(A, B, C [][]float64) ([][]float64, [][]float64, [][]float64, []float64) {
    A = transpose(A)
    B = transpose(B)
    C = transpose(C)

    new_A := [][]float64{}
    new_B := [][]float64{}
    new_C := [][]float64{}

    for i := range len(A) {
        a, _ := lagrange_interp(A[i])
        new_A = append(new_A, a)
    }

    for j := range len(B) {
        b, _ := lagrange_interp(B[j])
        new_B = append(new_B, b)
    }

    for k := range len(C) {
        c, _ := lagrange_interp(C[k])
        new_C = append(new_C, c)
    }

    Z := []float64{ 1 }

    for l := 1; l <= len(A[0]); l++ {
        Z = multiply_polys(Z, []float64{ float64(-l), 1.0 })
    }

    return new_A, new_B, new_C, Z
}


func create_solution_polynomials(
    r []float64,
    A, B, C [][]float64,
) ([]float64, []float64, []float64, []float64) {

    Apoly := []float64{}
    Bpoly := []float64{}
    Cpoly := []float64{}

    for i := range len(A) {
        mpoly := multiply_polys([]float64{ r[i] }, A[i])
        Apoly = add_polys(Apoly, mpoly)
    }

    for i := range len(B) {
        Bpoly = add_polys(Bpoly, multiply_polys([]float64{ r[i] }, B[i]))
    }

    for i := range len(C) {
        Cpoly = add_polys(Cpoly, multiply_polys([]float64{ r[i] }, C[i]))
    }

    o := subtract_polys(multiply_polys(Apoly, Bpoly), Cpoly)

    for i := 1; i <= len(A[0]); i++ {
        assertZeroF64(eval_poly(o, float64(i)))
    }

    return Apoly, Bpoly, Cpoly, o
}


func create_divisor_polynomial(sol, Z []float64) []float64 {
    quot, rem := div_polys(sol, Z)
    for i := range len(rem) {
        assertZeroF64(rem[i])
    }
    return quot
}


func assertZeroF64(a float64) (bool, error) {

    dust := math.Pow(10, -8)
    isCloseToZero := math.Abs(a) < dust

    if !isCloseToZero {
        errMsg := errors.New(fmt.Sprintf("Term non-zero: %v", a))
        fmt.Println(errMsg)
        return false, errMsg
    } else {
        return true, nil
    }
}