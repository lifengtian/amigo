package main

import (
	"fmt"
	"flag"
)

var (
	v1 int
)

func GetName() (firstName, lastName, nickName string) {
	return "Jacky", "Chan", "Long"
}

func MyPrintf(args ...interface{}) {
	for _, arg := range args {
		switch arg.(type) {
		case int:
			fmt.Printf("%d\n", arg)
		case string:
			fmt.Printf("%s\n", arg)
		case float32:
			fmt.Printf("%f\n", arg)
		default:
			fmt.Printf("default %T %v\n", arg, arg)
		}

	}

}

func main() {
	var infile *string = flag.String("i", "infile", "File input")

	v10 := 10
	fmt.Printf("v1 : %d v10 : %d\n", v1, v10)
	v1, v10 = v10, v1
	fmt.Printf("v1 : %d v10 : %d\n", v1, v10)

	// challenge: multiple-value to single-value
	fmt.Println(GetName())

	const Pi = 3.14159265358979323846 + 1
	const (
		c0 = iota
		c1
		c2
	)
	v2 := 2.1
	fmt.Printf("Pi = %f\nType(Pi) = %T\nint32(Pi) = %d\n", Pi, Pi, int32(v2)) // can't say int32(Pi)
	fmt.Printf("c0 = %d\nc1 = %d\nc2 = %d\n", c0, c1, c2)

	v3 := complex(3, 3)
	fmt.Printf("complex(3,3) = %v\nReal = %f\nImag = %f\n", v3, real(v3), imag(v3))

	v4 := "BHello world!"
	ch := v4[0]
	fmt.Printf("string =%s\nch = %d\nch = %c\nType(ch) = %T\n", v4, ch, ch, ch)

	modify := func(array *[5]int) {
		array[0] = 10
		fmt.Println("In modify, array values:", *array)
	}

	v5 := [5]int{1, 2, 3, 4, 5}
	modify(&v5)
	fmt.Println("In main, array values:", v5)

	v6 := v5[2:5]
	fmt.Println("v6 = ", v6)

	v7 := make(map[string]int)
	v7["one"] = 1
	v7["two"] = 2
	_, ok := v7["one"]
	if ok {
		fmt.Printf("v7 = %v\n", v7)
	} else {
		fmt.Println("Can't find key=one")
	}

	v8 := []interface{}{1, "one", 1.2, 0123}
	MyPrintf(v8)
	MyPrintf(v7)
	MyPrintf(v5)
	MyPrintf(1, "hello", 1.2)

	flag.Parse()
	fmt.Println("Input infile is ", *infile)


}
