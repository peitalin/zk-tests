package utils

import (
	"fmt"
)

func PrintHeading(s string) {
    fmt.Println("\n------------------------------------------")
    fmt.Println(s)
    fmt.Println("------------------------------------------\n")
}

func TestZipImplementations() {

	arr1 := []float64{1,2,3,4}
	arr2 := []float64{5,6,7,8}
	arr3 := []float64{10,11,12,15}

    PrintHeading("Three ways to implement zip()")
	fmt.Println("arr1:", arr1)
	fmt.Println("arr2:", arr2)
	fmt.Println("arr3:", arr2, "\n")
	zipped1 := Zip(arr1, arr2, arr3)
	fmt.Println("zipped1:", zipped1, "(basic for-loop)")

    zipped2 := Zip2(arr1, arr2, arr3)
	fmt.Println("zipped2:", zipped2, "(using iterator)")

    zipped3 := Zip3(arr1, arr2, arr3)
	fmt.Println("zipped3:", zipped3, "(using go channels)")

}

func Zip3(lists ...[]float64) [][]float64 {
	zipChannel := zipChannel(lists...)
    zipped := make([][]float64, 0)
	for tuple := range zipChannel {
        zipped = append(zipped, tuple)
	}
    return zipped
}

func zipChannel(lists ...[]float64) chan []float64 {
    out := make(chan []float64)
    go func() {
        defer close(out)
        for i := range len(lists[0]) {
            tup := make([]float64, len(lists))
            for j := range lists {
                tup[j] = lists[j][i]
            }
            out <- tup
        }
    }()
    return out
}

func Zip2(lists ...[]float64) [][]float64 {
	iter := zipIterator(lists...)
    zipped := make([][]float64, len(lists[0]))
    i := 0

	for tuple := iter(); tuple != nil; tuple = iter() {
        newTuple := []float64{}
        newTuple = append(newTuple, tuple...)
        zipped[i] = newTuple
        i++
	}
    return zipped
}

func zipIterator(lists ...[]float64) func() []float64 {
    tup := make([]float64, len(lists))
    i := 0
    return func() []float64 {
        for j := range lists {
            if i >= len(lists[j]) {
                return nil
            }
            tup[j] = lists[j][i]
        }
        i++
        return tup
    }
}

func Zip(lists ...[]float64) [][]float64 {
    zip := make([][]float64, len(lists[0]))

    for i := range len(lists[0])  {
        zip[i] = make([]float64, len(lists))
        for j := range lists {
            if i >= len(lists[j]) {
                return nil
            }
            n := lists[j][i]
            zip[i][j] = n
        }
    }
    return zip
}

