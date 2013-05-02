// Hello is a trivial example of a main package.
package main

import (
        "newmath"
        "fmt"
	"time"
	"math/rand"
)

func main() {
	start := time.Now()
	rand.Seed ( start.Unix() )
        fmt.Printf("Hello, world. I can do math :) Look: \nSqrt(2) =")
	time.Sleep( time.Duration ( rand.Intn(10) ) * time.Second )
	fmt.Printf(" %v\tin %d seconds\n", newmath.Sqrt(2 ),  int (time.Since(start).Seconds() ) )
	
}
