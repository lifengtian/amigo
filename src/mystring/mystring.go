// learning strings package
// From: Mar 30, 2013
// Notes: Fields: like my @a=split/\s/
//            you can say for i,v := range s.Fields
//    
package main

import (
	"fmt"
	"strings"
	"os"
)

func main() {
	s := "hello world how are you? "
	fmt.Println("case 1: split strings")
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
	for _,v := range strings.Fields(s) {
		fmt.Printf("%s\n",v)
	}

	s2 := "Welcome to the Earth!"
	fmt.Println("case 2: concatenate strings")
	fmt.Println(s + s2)


	s = "/the/path/is/here/test.bam"
	fmt.Println("case 3: check suffix and prefix")
	if strings.HasSuffix(s, ".bam") { fmt.Println("Find a bam file")}
	s = "#CHROM\tblah\tblah"
	if strings.HasPrefix(s,"#") { fmt.Println("Find comment in VCF file")}

	s = "#CHROM\tlocation\tsample1\tsample2"
	a := strings.Fields(s)
	l := len(a)
	fmt.Fprintf(os.Stdout,"string:%s\nlen:%d\n", s, l)

}
