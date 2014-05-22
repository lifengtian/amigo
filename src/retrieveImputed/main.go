package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

func Ver1() {
	var file io.Reader
	var err error

	file, err = os.Open(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(file)
	var text string
	var fields []string
	var n_fields int
	var line_no int
	var m_chr, m_pos string

	n_fields = 3
	for scanner.Scan() {
		text = scanner.Text()
		//fmt.Printf("text:\n%s\n", text)
		line_no++
		fields = strings.Fields(text)
		if n_fields != len(fields) {
			fmt.Println("error", text)
			log.Fatal("n_fields != len(fields) at line ", line_no)
		}
		m_chr = fields[0]
		m_pos = fields[1]
		n, _ := strconv.Atoi(m_pos)
		d := strconv.Itoa(n / 5000000 + 1)
		//fmt.Fprintf(os.Stdout, "chr:%s\npos:%d\npos/5000000:%s\n", m_chr, n, d)
		fmt.Fprintf(os.Stdout, `echo "~/s3cmd/s3cmd -c ~/patrick.s3cfg get `+"s3://cagpatrick_go_imputation/merged-trad-forward-chr"+m_chr+"_"+d+".imputed "+` " | qsub -cwd -l mem_free=20G,h_vmem=24G` + "\n")
	}
}

func main() {
	Ver1()
}
