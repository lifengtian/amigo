//you have a list/slice of numbers
//calculate the sum of it
package main

import "fmt"
import "math/rand"
import "math"

func sum(n []int) int {
	s := 0
	for _, v := range n {
		s += v
	}
	return s
}

func sumOfN( n []int, N int) int {
	s := 0
	for _, v := range n {
		s += int(math.Pow(float64(v), float64(N) ) ) 
	}
	return s
}

func main() {
	//to make things interesting, we ask rand.Perm to produce a random permutation of integer numbers 1 to N
	numbers := rand.Perm(21)
	//numbers := [...]int{1, 2, 3, 4, 5}
	//fmt.Println( sum(numbers[:] ) )
	//fmt.Println( sum(n2 ) )

	for i :=1 ; i<10 ; i++ {
		fmt.Printf("n=%d\t%d\n",i, sumOfN ( numbers[:], i ) )
	}
}
