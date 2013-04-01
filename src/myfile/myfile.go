
package main

import (
	"fmt"
	"os"
	"log"
)

// cat the content of the file
func CatFile (fn string) error  {
	fmt.Println("CatFile: ", fn)

	file, err := os.Open(fn) // simply open a file with default settings
	if err != nil {
		return fmt.Errorf("failed opening %s: %s", fn, err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return fmt.Errorf("failed file stat for %s %s:", stat, err)
	}

	fmt.Printf("File %s has stat: %s\n", fn, stat)
	
	fmt.Printf("ModTime: %s\n", stat.ModTime() )
	fmt.Printf("Size: %d\n", stat.Size() )
	fmt.Printf("Mode: %d\n", stat.Mode() )

	return nil
}


func main() {
	file, err := os.OpenFile("output", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	if err != nil {
		log.Fatalf("failed opening file for writing: %s", err)
	}
	defer file.Close()

	CatFile("myfile.go")
	
}