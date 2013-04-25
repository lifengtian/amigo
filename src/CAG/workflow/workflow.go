// A typical whole-exome pipeline
// Goal is to make writing such things easier for a new workflow
// but not trying to make it a generic workflow writter. Things changes much too fast in this field.
package main

import (
	"CAG/utils"
	"flag"
	"os"
	
) 

var (
	sampleName   string
	jsonFileName string

	work         utils.Workflow

	dryrun *bool = flag.Bool("dryrun", false, "dry run")
	help   *bool = flag.Bool("help", false, "help")
)

/*  ./workflow -dryrun -config workflow.WES.00001.json -sn 0123456789  */
/*  ./workflow -dryrun -config pro.json -sn 0123456789  */

func main() {

	flag.StringVar(&jsonFileName, "config", "", "workflow configure file (json format)")
	flag.StringVar(&sampleName, "sn", "", "sample name")
	flag.Parse()

	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if !*dryrun {
		utils.Begin()
	}

	
	work = make(utils.Workflow)
	work.Init( jsonFileName, sampleName, *dryrun )
	work.Run( *dryrun )
	work.Finish()

	
		
}


/* /Users/tianl/workspace/amigo/src/workflow/workflow.go  */