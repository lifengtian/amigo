package main

import "fmt"

func main() {
	fmt.Println("Hello, playground")

	s := "田立峰"
	fmt.Println(s)
	for _, v := range []byte(s) {
		fmt.Printf("%x\n", v)
	}

	b := string ( []byte(s) )
	fmt.Println( "str2byte2str ", b)
}