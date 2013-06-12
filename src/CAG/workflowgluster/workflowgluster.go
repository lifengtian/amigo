// A typical whole-exome pipeline on respublica
// Goal is to make writing such things easier for a new workflow
// but not trying to make it a generic workflow writter. Things changes much too fast in this field.
package main

import (
	"CAG/utilsgluster"
	"flag"
	"os"
	"net/http"	
	"time"
	"io"
	"fmt"
) 

var (
	sampleName   string
	jsonFileName string

	work         utilsgluster.Workflow

	dryrun *bool = flag.Bool("dryrun", false, "dry run")
	help   *bool = flag.Bool("help", false, "help")
)

/*  ./workflow -dryrun -config workflow.WES.00001.json -sn 0123456789  */
/*  ./workflow -dryrun -config pro.json -sn 0123456789  */

// hello world, the web server
func HelloServer(w http.ResponseWriter, req *http.Request) {
	io.WriteString(w, "hello, world!\n")
}

func ShowStates(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "It's %s\n%s\n", time.Now() , work.States() )
}

func SetPassword (w http.ResponseWriter, req *http.Request) {
	req.SetBasicAuth("test", "garyowen")	
}
	
func StartServer( ) {
	http.HandleFunc("/hello", HelloServer)
	http.HandleFunc("/states", ShowStates )
	http.HandleFunc("/setpwd", SetPassword )
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		fmt.Println("StartServer error ",  err ) 
	}
}

func main() {

	flag.StringVar(&jsonFileName, "config", "", "workflow configure file (json format)")
	flag.StringVar(&sampleName, "sn", "", "sample name")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if !*dryrun {
		utilsgluster.Begin()
	}

	
	work = make(utilsgluster.Workflow)
	work.Init( jsonFileName, sampleName, *dryrun )
	go StartServer()

	work.Run( *dryrun )
	work.Finish()

	
		
}


/* /Users/tianl/workspace/amigo/src/workflow/workflow.go  */
