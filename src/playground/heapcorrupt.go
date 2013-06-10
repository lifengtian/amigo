package main

import (
	"fmt"
	"runtime"
	"strconv"
)

var s = "\xcc"

type a struct {
	b
	o *a
	t *func()
}

type b struct {
	c chan error
	d *b
}

func address(i interface{}) int {
	addr, err := strconv.ParseUint(fmt.Sprintf("%p", i), 0, 0)
	if err != nil {
		panic(err)
	}
	return int(addr)
}

func main() {
	var w *a
	for i := 0; i < 3; i++ {
		p := &a{}
		p.c = make(chan error)
		p.d = &w.b
		p.o = w
		p.t = new(func())
		*p.t = func() { fmt.Println("failed") }
		w = p
	}
	runtime.GC()
	r := address(&s)
	var l []*int
	for i := 0; i < 1000; i++ {
		b := new(int)
		*b = r
		l = append(l, b)
	}
	(*w.o.o.t)()
}