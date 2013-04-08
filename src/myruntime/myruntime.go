package main

import (
	"fmt"
	"runtime"
	"os"
)

func main(){
	fmt.Fprintf(os.Stdout, "NumCPU: %d\n", runtime.NumCPU() )
	fmt.Fprintf(os.Stdout, "Version: %s\n", runtime.Version() )
	



}
