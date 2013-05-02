package main

import (
	"fmt"
)

func Ver1() {
	one := []string{"  1  ", "  1  ", "  1  ", "  1  ", "  1  "}
	two := []string{" 222 ", "2   2", "   2 ", "  2  ", " 2   ", "2    ", "22222"}

	sliceOfslice := [][]string{one, two}

	for _, v := range sliceOfslice {
		for _, m := range v {
			fmt.Printf("%s\n", m)
		}
	}

}

func Ver2() {
	sliceOfslice := [][]string{   {"    1     ", 
                                                      "    1     ", 
                                                      "    1     ",
						"    1     ",
						"    1     ",
						"    1     ",
						"    1     " }, 
                                                 {     "  222 ", 
                                                      "2      2 ", 
                                                      "      2   ",  
                                                      "    2     ", 
                                                      "  2       ", 
                                                      "2         ", 
                                                      "22222"}, 
					{       " 33333", 
       						"        3", 
						"         3", 
						"33333", 
						"         3", 
						"        3", 
                                                      "33333"} }
	result := []string{"","","","","","",""}

	for _, v := range sliceOfslice {
		for i, m := range v {
			result[i] += m
		}
		
	}
	
	for _, v := range result {
		fmt.Printf("%s\n", v)
	}

}
func main(){
	Ver2()
}
