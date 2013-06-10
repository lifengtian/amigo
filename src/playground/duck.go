package main

import (
	"fmt"
)

type Duck struct {
	Name string
}

func (d *Duck) Quack() {
	fmt.Println("Duck Quack!")
}

type DonaldDuck struct {
	Age int
	Duck
}

func (d *DonaldDuck) Quack() {
	fmt.Println("DonaldDuck Also Quack!")
}

func main() {
	dduck := new(DonaldDuck)
	fmt.Printf("%T %v %T %v\n", dduck,dduck,dduck.Duck, dduck.Duck)
	dduck.Duck.Quack()
}
