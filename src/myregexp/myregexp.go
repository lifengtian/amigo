package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

func main() {
	var file io.Reader
	var err error
	s := regexp.MustCompile(`AF=([\.\d]+)`)
	var l int // length

	if l = len(os.Args); l > 1 {
		file, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		fmt.Fprintf(os.Stdin, "os.Args len:%d", l)
	} else {
		file = os.Stdin
	}

	scanner := bufio.NewScanner(file)
	var m string
	for scanner.Scan() {
		m = s.FindString(scanner.Text())
		if l = len(m); l > 0 {
			fmt.Println(m)
		}
	}
}
