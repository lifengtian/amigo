//Play with channel and goroutine
package main

import "fmt"
import "time"

func sum(v []int, c chan int) {
	sum := 0
	for _, value := range v {
		sum += value
	}
	c <- sum
}

func main() {
	values := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	t := time.Now()
	c := make(chan int, 4)
	go sum(values[:len(values)/4], c)
	go sum(values[len(values)/4:len(values)/2], c)
	go sum(values[len(values)/2:len(values)*3/4], c)
	go sum(values[len(values)*3/4:], c)

	sum1, sum2, sum3, sum4 := <-c, <-c, <-c, <-c
	fmt.Printf("Result in %v : %d + %d + %d + %d = %d\n", time.Since(t), sum1, sum2, sum3, sum4, sum1+sum2+sum3+sum4)
}
