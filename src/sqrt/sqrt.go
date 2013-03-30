package main

import (
    "fmt"
    "math"
)

func sqrt(x float64) float64 {
    z := float64(1)
    for i := 1; i< 11; i++ {
    	z = z - (z * z -x)/2/z
    }
    return z
}

func main() {
    fmt.Println(sqrt(2)-math.Sqrt(2))
}
