package main

import (
	"fmt"
	"net"
	)

func main(){
	addrs, _ := net.LookupHost("localhost")
	fmt.Println(addrs)
	addrs2, _ := net.LookupIP("localhost")
	fmt.Println(addrs2)
	addrs3, _ := net.LookupAddr("localhost")
	fmt.Println(addrs3)
	addrs4, _ := net.InterfaceAddrs()
	fmt.Println(addrs4)
	}
