
// time report time since
package main

import (
	"fmt"
	"time"
)

func main() {
	start:=time.Now()
	//fetch("http://www.google.com")
	fmt.Println("hello")
	fmt.Println(time.Since(start))
}
