package main

import "math/big"
import "fmt"

func main(){
	var z big.Int
	x := big.NewInt ( 30.0 )
	y := big.NewInt ( 23.0 )
	m := big.NewInt ( 0.0 )
	z.Exp( x, y, m )
	fmt.Printf ("%v\n", z )

}
