//package vcf provides an interface to VCF files
package main


import (
	"fmt" 
	"os"
	"regexp"
	"bufio"
	"log"
	"io"
	"strings"
	"time"
	"strconv"
)



func main(){
	var file io.Reader
	var err error
	s := regexp.MustCompile(`AF=([\.\d]+)`)
	var l int // length

	t0 := time.Now()
	if l = len(os.Args); l> 1 {
		file, err = os.Open(os.Args[1])
		if err != nil { log.Fatal(err)}
		fmt.Fprintf(os.Stdin, "os.Args len:%d", l )
	} else {
		file = os.Stdin
	}

	scanner := bufio.NewScanner(file)
	var m string
	var text string
	var fields []string
	var n_fields int
	var line_no int 

	for scanner.Scan() {
		text = scanner.Text()
		line_no++
		// retrieve number of samples
		// should it be that number of samples be provided in the header???
		if strings.HasPrefix(text, "#CHROM") {
			//parse CHROM line
			// #CHROM POS ID REF ALT QUAL FILTER INFO FORMAT
			fields = strings.Fields(text)
			n_fields = len(fields)
			n_samples := n_fields - 9
			fmt.Println("Number of samples: ", n_samples)
		} else if ! strings.HasPrefix(text,"##"){
			fields = strings.Fields(text)
			if n_fields != len(fields) {log.Fatal("n_fields != len(fields) at line ", line_no)}
			m = s.FindString( text )
			if l = len(m); l>0 {
				if i,_ := strconv.ParseFloat(strings.Trim(m,"AF="), 64); i < 0.1 {
					fmt.Println(line_no, i, m)
				}
			}
		}
	}

	t1 := time.Now()
	fmt.Printf("Total run time: %v\n", t1.Sub(t0) )
}