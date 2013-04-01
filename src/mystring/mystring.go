// learning strings package
// Mar 30, 2013
// Fields: like my @a=split/\s/
//      you can say for i,v := range s.Fields
package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "hello world how are you? "
	fmt.Printf("Fields are: %q", strings.Fields("  foo bar  baz   "))
	for _,v := range strings.Fields(s) {
		fmt.Printf("%s\n",v)
	}
}
