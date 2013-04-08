// Learn to play with the file processing and io.Copy

package main

import (
	"io"
	"log"
	"os"
)


func main(){
	from, err := os.Open("myfile.go")
	if err != nil { log.Fatal(err)}

	to, err := os.OpenFile("to.go", os.O_WRONLY | os.O_CREATE, 0644)
	if err != nil { log.Fatal(err)}

	_, err := io.Copy(to, from)
	if err != nil { log.Fatal(err)}
	
}