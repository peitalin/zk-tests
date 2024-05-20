package zktest

import (
	"fmt"
)


func ZkTest1() {
    fmt.Println("\n----------------------------------")
    fmt.Println("ZK Exercise One")
    fmt.Println("----------------------------------\n")
}

func ZkTest2() {
    fmt.Println("\n----------------------------------")
    fmt.Println("ZK Exercise Two")
    fmt.Println("----------------------------------\n")
}

func ZkTest3() {
    fmt.Println("\n----------------------------------")
    fmt.Println("ZK Exercise Three")
    fmt.Println("----------------------------------\n")
}

type Polynomial struct {
    Modulus int
}

func (poly *Polynomial) EvalPolyAt(p_coeffs []int, x int) int {

    y := 0
    power_of_x := 1

    for _, p_coeff := range p_coeffs {

        fmt.Println("p_coeff: ", p_coeff)
        fmt.Println("power_of_x: ", power_of_x)
        y += power_of_x * p_coeff
        power_of_x = (power_of_x * x) % poly.Modulus
        fmt.Println()

    }
    fmt.Printf("%v = %v mod %v\n", y % poly.Modulus, y, poly.Modulus)

    return y % poly.Modulus
}
