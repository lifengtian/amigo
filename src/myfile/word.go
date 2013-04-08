// Read a file
// count occurrence of words
// output the words
package main

import (
"os" 
//"io"
"fmt"
"log"
"text/scanner"
)


func main(){
	//read a file
	file, err := os.Open("word.go")
	if err != nil { log.Fatal(err)}
	defer file.Close()

	//count occurrence of words

	m := make(map[string]int)

	var s scanner.Scanner
	s.Init(file)
	tok := s.Scan()
	for tok != scanner.EOF {
		// do sth
		m[s.TokenText()] += 1
		//m[scanner.TokenString(tok)] += 1
		tok = s.Scan()
	}

	for i,v := range m {
		fmt.Printf("%s: %d\n", i, v)
	}
}