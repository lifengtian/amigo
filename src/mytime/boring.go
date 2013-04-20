/*Package main aims to practice the time package and channels, go routines

*/
package main

// import has to be said before func def (so it is ordered )
import (
	"fmt"
	"math/rand" // /Users/tianl/workspace/go/src/pkg/math/rand
	"time"
)

func doBoring() {

	c := make(chan int, 10)
	for i := 0; i < 10; i++ {
		go boring(i, c)
	}

	for i := 0; i < 10; i++ {
		<-c
	}
	fmt.Println("")
}

func boring(i int, c chan int) {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	time.Sleep(time.Duration(r.Intn(5)) * time.Second)
	fmt.Printf("%d ", i)
	c <- 1
	c <- 2
}

func main() {
	for i := 0; ; i++ {
		doBoring()
	}
}
