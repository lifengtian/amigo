//Play with channel and goroutine
//42 us for 2 go
//33 us for 4 go
//Need to figure out how to test with different number of channels
package main

import "fmt"
import "time"
import "os"
import "math/rand"
import "runtime"

// sum calculates the sum of all members of the slice
func sum(v []int, c chan int64) {
	var sum int64 = 0
	for _, value := range v {
		sum += int64(value)
	}

	// instead of return the value, we send it to the channel
	c <- sum
}

func main() {
	NCPU := runtime.NumCPU()
	N := NCPU / 2
	runtime.GOMAXPROCS(N)
	fmt.Printf("NCPU : %d N : %d\n", NCPU, N)

	values := rand.Perm(100000000)
	t := time.Now()
	c := make(chan int64, N)

	for i := 0; i < N; i++ {
		go sum(values[len(values)*i/N:len(values)*(i+1)/N], c)
	}

	var mysum int64 = 0

	// Drain the channel
	for i := 0; i < N; i++ {
		s := <-c
		mysum += s
		fmt.Println("go ", i, s)
	}

	fmt.Printf("Result in %v %v = %d\n", os.Args[0], time.Since(t), mysum)
}
