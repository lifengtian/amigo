package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
)

/* A pipeline has multiple steps
* each step can be run on same or different node which depends on 
* the config e.g., shall I qsub it or run it via exec.Command?
*
* I prefer to build a pipeline first based on a prescription
* followed by generate actual running scripts (which is important)
* I want to see a BASH script file right there for each pipeline
* A step has parameters

* How can I chain the input and output naturally?

* Run a command wait till it finishes
* check the return errors

* common themes:
* 1. where to find the binaries for the command? what if it does not exist?
     where can I find it? Dependencies and packaging for external dependencies.
     I can't rewrite everything and I can't easily package them, can I? How about
     the datasets, vcfs?
  2. How do we know a step successfully done? what if return value was not correctly set e.g s3?
     Is the VCF or BAM corrupted? How about checksum? 
  3. How to log each step (abstracted as a computation)?

*/

// for each computation, there is a config (may be based on some template)
type Config struct {
	CmdName      string
	CmdArguments string
}

func add() {}

func PrepareCompute(c Config) (err error) {
	return nil
}

func DoCompute(c Config) (err error) {
	cmd := exec.Command(c.CmdName, c.CmdArguments)
	cmd.Stdout = os.Stdout
	err = cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	return nil
}

func Compute(c Config) (err error) {
	PrepareCompute(c)
	if err != nil {
		log.Fatalf("PrepareCompute failed: %q", err)
	} // need to learn to produce specific error 

	DoCompute(c)
	if err != nil {
		log.Fatalf("DoCompute failed: %q", err)
	} // need to learn to produce specific error 

	return nil
}

func main() {
	cmd := exec.Command("tr", "a-z", "A-Z")
	cmd.Stdin = strings.NewReader("some input")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("in all caps: %q\n", out.String())
	c := Config{CmdName: "echo", CmdArguments: " hello world!"}
	DoCompute(c)
}
