
//Learn to work with closure
package main

import "fmt"

// return a func(int) int
func adder() func(int) int {
	sum := 0
	return func(x int) int {
		sum += x
		return sum
	}
}

// create a bank account, one can cash, deposit, or check balance
func makeaccount (value int) func(string, int) int {
	amount := value
	return func(s string, v int) int {
		switch s {
		case "cash":
			amount -= v
		case "deposit":	
			amount += v
		}
		return amount
	}
}

func main() {
	pos, neg := adder(), adder()
	for i := 0; i < 10; i++ {
        fmt.Println(
			pos(i),
			neg(-2*i),
        )
	}

	acct := makeaccount(100)
	fmt.Println( acct ("cash", 40))
	fmt.Println( acct ("cash", 40))
	fmt.Println( acct ("dep", 40))

}

