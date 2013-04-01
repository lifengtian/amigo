package main

import (
	"io"
	"os"
	"strings"
	//"fmt"
)

type rot13Reader struct {
	r io.Reader
}

type tianReader struct {
	r io.Reader
}

func rot13(b byte) byte {
	switch {
	case 'A' <= b && b <= 'M':
		b = (b - 'A') + 'N'
	case 'N' <= b && b <= 'Z':
		b = (b - 'N') + 'A'
	case 'a' <= b && b <= 'm':
		b = (b - 'a') + 'n'
	case 'n' <= b && b <= 'z':
		b = (b - 'n') + 'a'
	}
	return b
}

// make rot13Reader conform to the Reader interface
func (b *rot13Reader) Read(p []byte) (n int, err error) {
	n, err = b.r.Read(p)

	for i := range p[:n] {
		p[i] = rot13(p[i])
	}
	return
}


func (b *tianReader) Read(p []byte) (n int, err error) {
	n, err = b.r.Read(p)
	
	for i := range p[:n] {
		p[i] = 'B' - 'B' + p[i]
	}
	return
}
func main() {
	// fmt.Println(string(rot13('A')), string(rot13('x')))
	s := strings.NewReader(
		"Lbh penpxrq gur pbqr!")
	r := rot13Reader{s}
	t := tianReader{&r}
	io.Copy(os.Stdout, &t)
}
