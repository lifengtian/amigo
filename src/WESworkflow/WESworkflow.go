package main

import (
	"encoding/json"
	"os"
	//"log"
	"bufio"
	"fmt"
	"io/ioutil"
	"time"
)

type WorkflowConfig struct {
	Ref   string
	Dbsnp string
	Mills string
}

func ( c WorkflowConfig) Print() {
	fmt.Printf("\t%#v\n", c)
}


func Config(filename string) (v WorkflowConfig) {
	content, _ := Contents(filename)
	json.Unmarshal(content, &v)
	//fmt.Printf("%v Read config from file: %s\n%+v\n", time.Now(), filename, v)
	return
}

func Contents(filename string) (contents []byte, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)

	contents, err = ioutil.ReadAll(reader)
	return
}

func Trace(s string) string {
	fmt.Printf("%v Entering %s\n", time.Now(), s)
	return s
}

func Un(s string) {
	fmt.Printf("%v Leaving %s\n", time.Now(), s)
}


var config WorkflowConfig

func main() {
	// play with marshal
	//t := WorkflowConfig{"/mnt/here/ref", "/mnt/there/dbsnp", "/mnt/here/again/mills"}
	//out, err := os.OpenFile("output", os.O_WRONLY | os.O_CREATE | os.O_TRUNC, 0644)
	//defer out.Close()
	//b, _ := json.Marshal(t)

	//fmt.Printf("json : %#v", b)
	//ioutil.WriteFile("output.txt",b, 0644)

	// play with unmarshal from a file
	//Config("output.txt")
	config.Print()
}

func init() {
	defer Un(Trace("init"))
	config = Config("output.txt")

}
