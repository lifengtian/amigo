package main

import (
	"fmt"
	"crypto/md5"
	"io"
)

func main(){
	h := md5.New()
	io.WriteString(h, "password")
	fmt.Printf("%x\n", h.Sum(nil) )
	io.WriteString(h, "this will be encrypted")
	fmt.Printf("%x", h.Sum(nil) )

}
