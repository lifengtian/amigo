package main

import "fmt"
import "unicode"

type 字符串 string

func main() {
	var 姓名 字符串
	姓名 = "你好"
	tables := []*unicode.RangeTable{unicode.Han}
	fmt.Println("Hello, playground", 姓名)
	fmt.Printf("Han type: %T\n", unicode.Han)
	for i := 1; i < unicode.MaxRune ; i++ {
		if unicode.IsOneOf(tables, rune(i)) {
			fmt.Printf("%d : %c\n", i, i)
		}
	}
}
