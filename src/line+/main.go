package main

import (
	"bufio"
	"os"
	"strconv"
)

func main() {
	lineNo := 1
	r := bufio.NewReader(os.Stdin)
	for {
		line, err := r.ReadString('\n')
		os.Stdout.WriteString(strconv.Itoa(lineNo) )
		os.Stdout.Write([]byte{'\t'})
		os.Stdout.WriteString(line)
		if err != nil {
			break
		}
		lineNo++
	}
}
