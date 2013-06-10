package main

import (
	"sort"
	"fmt"
	"genome" 
)

func main() {
	l := genome.Locs{{"chr1", 4}, {"chr2", 5}, {"chr1", 1}, {"chr2",1} }
	fmt.Println(l)
	sort.Sort(genome.ByLoc{l})
	fmt.Println(l)
}
