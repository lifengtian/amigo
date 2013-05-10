//package vcf provides an interface to VCF files
//bufio has a maxTokenSize 64 * 1024 which turns out to be too small for some INDEL records
//I changed it to 640 * 1024 in scanner.go
//TODO: use Reader instead
// LFT May 2013
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

//	"strconv"
)

func main() {
	var file io.Reader
	var err error
	af := regexp.MustCompile(`AF=([\.\de\-,]+)`)
	an := regexp.MustCompile(`AN=([\.\d]+)`)
	ac := regexp.MustCompile(`AC=([\.\d,]+)`)
	gene := regexp.MustCompile(`Gene=([^;]+)`)
	var l int // length

	t0 := time.Now()
	if l = len(os.Args); l > 1 {
		file, err = os.Open(os.Args[1])
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Fprintf(os.Stdin, "os.Args len:%d", l )
	} else {
		file = os.Stdin
	}

	scanner := bufio.NewScanner(file)
	var text string
	var fields []string
	var n_fields int
	var line_no int
	var m_af, m_an, m_ac, m_gene string

	for scanner.Scan() {
		text = scanner.Text()
		//fmt.Printf("text:\n%s\n", text)
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
			fmt.Printf("CHROM\tPOS\tID\tREF\tALT\tFILTER\tAltAlleleFreq\tTotalNumberofAllAlleles\tAltAlleleNumber\tGene\n")
		} else if !strings.HasPrefix(text, "##") {
			fields = strings.Fields(text)
			if n_fields != len(fields) {
				fmt.Println("error", text)
				log.Fatal("n_fields != len(fields) at line ", line_no)
			}

			m_af = af.FindString(text)
			m_an = an.FindString(text)
			m_ac = ac.FindString(text)
			m_gene = gene.FindString(text)
			if l = len(m_af); l > 0 {
				//if i,_ := strconv.ParseFloat(strings.Trim(m,"AF="), 64); i < 0.1 {
				fmt.Printf("%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\t%s\n", fields[0], fields[1], fields[2], fields[3], fields[4], fields[6], m_af, m_an, m_ac, m_gene)
				//}
			}
		}
		//if a := line_no % 10 ; a == 0  {
		//	fmt.Fprintf(os.Stderr,"%d\n",line_no)
		//}
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

	t1 := time.Now()
	fmt.Fprintf(os.Stderr, "Total run time: %v\n", t1.Sub(t0))
}
