package main

import "fmt"
import "unsafe"
import "reflect"
import "runtime"

func main() {
	var b []byte
	s := "hello world"
	h := (*reflect.SliceHeader)(unsafe.Pointer(&b))
        h.Data = (*reflect.StringHeader)(unsafe.Pointer(&s)).Data
        h.Len = len(s)
        h.Cap = len(s) 
	fmt.Printf("h has Type %T value %v \n", h , h.Data )
	g := string( b ) 
	fmt.Printf("g has Type %T value %v \n", g , g )

	c := string ( []byte(s)[:3] )
	fmt.Printf("c has Type %T value %v \n", c, c )


	var m *runtime.MemStats = new(runtime.MemStats)

	runtime.ReadMemStats(m)

	fmt.Printf("Memory usage\nSys %d\nAlloc %d\nStack %d\nHeap %d\n", m.Sys, m.Alloc, m.StackInuse, m.HeapAlloc)



}