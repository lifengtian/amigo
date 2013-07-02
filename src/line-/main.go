package main

import (
	"bufio"
	"os"
	"strings"
)

func main() {
	r := bufio.NewReader(os.Stdin)
	firstLine := true
	for {
		line, err := r.ReadString('\n')
		lines := strings.Fields(line)
		if len(line) > 0 {
			if ! firstLine {
				os.Stdout.WriteString("\n")
			} else {
				firstLine = false
			}
			for i, v := range lines {
				if i > 0 {
					os.Stdout.WriteString(v + "\t")
				} 
			}
		} else {
			os.Stdout.WriteString(line)
		}
		if err != nil {
			break
		}
	}
}
