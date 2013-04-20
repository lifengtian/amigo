//Package Cmd experiments with external processes
// https://gobyexample.com/spawning-processes
package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"
	"net/http"
)

import _ "net/http/pprof"

var (
	sampleName string
	filename   string

	dryrun *bool = flag.Bool("dryrun", false, "dry run")
	help   *bool = flag.Bool("help", false, "help")
)

func MyPrintf(s string, i int) {
	fmt.Fprintf(os.Stderr, s, i)
}

func isFile(filename string) bool {
	if strings.HasPrefix(filename, "s3:") {
		return true
	}

	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func Begin() {
	fmt.Printf("Begin at %v\n", time.Now())
	addrs, _ := net.InterfaceAddrs()
	fmt.Println(addrs)
}

func Finish() {
	fmt.Printf("End at %v\n", time.Now())
}

// Cmd holds content of a command
type Cmd struct {
	Command    string
	Parameters string
}

func (cmd Cmd) Do(c chan int) {

	if *dryrun {
		fmt.Printf("%s %v\n", cmd.Command, cmd.Parameters)
	}

	command := exec.Command("bash", "-c", cmd.Command+" "+cmd.Parameters)
	var out bytes.Buffer
	command.Stdout = &out
	var errout bytes.Buffer
	command.Stderr = &errout

	if !*dryrun {
		fmt.Printf("%v %s with %v Start\n", time.Now(), cmd.Command, cmd.Parameters)
		err := command.Run()
		if err != nil {
			log.Fatalf("Error %s %v %s\n", cmd.Command, err, errout.String())
		}

		// remove trailing newline
		s := strings.TrimRight(out.String(), "\n")
		fmt.Printf("%v %s with %v End output: %s\n", time.Now(), cmd.Command, cmd.Parameters, s)
	}
	//complete
	c <- 1
}

var work map[string]string

func main() {

	flag.StringVar(&filename, "config", "", "workflow configure file (json format)")
	flag.StringVar(&sampleName, "sn", "", "sample name")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(1)
	}

	if !*dryrun {
		Begin()
		go func() {
				log.Println(http.ListenAndServe("localhost:6060", nil))
		}()
	}

	var workflow_json []byte

	//processing arguments
	//if len(os.Args) < 3 {
	//	log.Fatalf("usage: %s json sampleName", os.Args[0])
	//}
	//filename := os.Args[1]

	//	filename := flag.Arg[0]
	//read from file
	workflow_json, err := ioutil.ReadFile(filename)

	err = json.Unmarshal(workflow_json, &work)
	if err != nil {
		fmt.Println("error:", err)
	}
	//fmt.Printf("workflow : %v", work)

	//fmt.Printf("Length of workflow json: %d\n", len(work))

	work["TMPDIR"] = "/mnt/run/" + sampleName
	work["JAVAGATK"] = work["JAVABIN"] + " -Xmx6g -Djava.io.tmpdir="+work["TMPDIR"] + " -jar "+ work["GATKJAR"] + " -et NO_ET -K " + work["GATKKEY"]

	work["UGinterval"] = strings.Join([]string{"-L", work["Exons"], " -L", work["TargetBed"], "--interval_padding 50 "}, " ")

	//sampleName := os.Args[2]

	//download db and tools from s3 
	prepDirs := " sudo mkdir -p " + work["PIPELINE"] + " ; sudo chmod 777 " + work["PIPELINE"] + " ; sudo mkdir -p " + work["TMPDIR"] + " ; sudo chmod 777 " + work["TMPDIR"] +
		" ; sudo mkdir -p " + work["Fastq"] + " ;  sudo chmod 777 " + work["Fastq"]  +
		" ; mkdir -p " + work["Fastq"] + "/" + sampleName

	cmds := []Cmd{{prepDirs, ""}}
	c := make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)

		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	prepTools := ""

	if !isFile(work["PIPELINE"] + "/" + work["Tools"]) {
		tools := work["S3cmd"] + " " + work["S3cfg.pipeline"] + " get s3://cagpipelines/" + work["Tools"] + " " + work["PIPELINE"] + "/"
		cmds = []Cmd{{tools, ""}}
		c = make(chan int, len(cmds))
		for i, v := range cmds {
			MyPrintf("cmds %d channel\n", i)
			go v.Do(c)
		}
		for i := 0; i < len(cmds); i++ {
			<-c
		}
		prepTools = " cd " + work["PIPELINE"] + ";" + " tar jxvf " + work["Tools"]
		cmds = []Cmd{{prepTools, ""}}
		c = make(chan int, len(cmds))
		for i, v := range cmds {
			MyPrintf("cmds %d channel\n", i)
			go v.Do(c)
		}
		for i := 0; i < len(cmds); i++ {
			<-c
		}

	}

	if !isFile(work["PIPELINE"] + "/" + work["DBs"]) {
		dbs := work["S3cmd"] + " " + work["S3cfg.pipeline"] + " get s3://cagpipelines/" + work["DBs"] + " " + work["PIPELINE"] + "/"
		cmds = []Cmd{{dbs, ""}}
		c = make(chan int, len(cmds))
		for i, v := range cmds {
			MyPrintf("cmds %d channel\n", i)
			go v.Do(c)
		}
		for i := 0; i < len(cmds); i++ {
			<-c
		}
		prepTools = " cd " + work["PIPELINE"] + ";" + " tar jxvf " + work["DBs"]
		cmds = []Cmd{{prepTools, ""}}
		c = make(chan int, len(cmds))
		for i, v := range cmds {
			MyPrintf("cmds %d channel\n", i)
			go v.Do(c)
		}
		for i := 0; i < len(cmds); i++ {
			<-c
		}

	}

	for i, v := range work {
		if isFile(v) {
			fmt.Fprintf(os.Stderr, "%s\t%s\n", i, v)
		} else {
			fmt.Fprintf(os.Stderr, "NNNNNNNNOOOOOO %s\t%s\n", i, v)
		}
	}

	r1 := work["Fastq"] + "/" + sampleName + "/pe_1.fq.gz"
	r2 := work["Fastq"] + "/" + sampleName + "/pe_2.fq.gz"

	s1 := work["TMPDIR"] + "/pe_1.sai"
	s2 := work["TMPDIR"] + "/pe_2.sai"

	timestamp := time.Now()
	lb := work["LB"]
	pl := work["PL"]
	pu := work["PU"]
	cn := work["CN"]

	getFastq1 := work["S3cmd"] + " -c /mnt/galaxyData/larus.s3cfg  sync " + work["FASTQBUCKET"] + "/" + sampleName + "/pe_1.fq.gz " + work["Fastq"] + "/" + sampleName + "/"
	getFastq2 := work["S3cmd"] + " -c /mnt/galaxyData/larus.s3cfg  sync " + work["FASTQBUCKET"] + "/" + sampleName + "/pe_2.fq.gz " + work["Fastq"] + "/" + sampleName + "/"

	fmt.Fprintln(os.Stderr, r1, r2, s1, s2, timestamp, lb, pl, pu, cn, getFastq1, getFastq2)

	cmds = []Cmd{{getFastq1, ""}, {getFastq2, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	fastqc1 := work["FASTQC"] + "/fastqc -o " + work["TMPDIR"] + " " + r1
	fastqc2 := work["FASTQC"] + "/fastqc -o " + work["TMPDIR"] + " " + r2

	bwa1 := work["BWA"] + " aln  -t 4 -I " + work["Reference"] + " " + r1 + " > " + s1
	bwa2 := work["BWA"] + " aln  -t 4 -I " + work["Reference"] + " " + r2 + " > " + s2

	fmt.Fprintf(os.Stderr, "bwa1: %s\n", bwa1)
	fmt.Fprintf(os.Stderr, "bwa2: %s\n", bwa2)
	cmds = []Cmd{{fastqc1, ""}, {fastqc2, ""}, {bwa1, ""}, {bwa2, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	sampe := work["BWA"] + ` sampe -r "@RG\tID:` + sampleName + `\tSM:` + sampleName + `" ` + work["Reference"] + " " + s1 + " " + s2 + " " + r1 + " " + r2 +
		" | " +
		work["SAMTOOLS"] + " view -h -S  - " +
		" | " +
		work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/SortSam.jar I=/dev/stdin O=/dev/stdout SO=coordinate VALIDATION_STRINGENCY=SILENT QUIET=true COMPRESSION_LEVEL=0 TMP_DIR=" + work["TMPDIR"] +
		" | " +
		work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/AddOrReplaceReadGroups.jar   I=/dev/stdin O=" + work["TMPDIR"] + "/SID" + sampleName + ".sorted.bam " + "SO=coordinate RGID=" + sampleName +
		" RGLB=" + lb + " RGPL=" + pl + " RGPU=" + pu + " RGCN=" + cn + " RGSM=" + sampleName +
		" CREATE_INDEX=true VALIDATION_STRINGENCY=SILENT QUIET=true TMP_DIR=" + work["TMPDIR"]

	fmt.Fprintf(os.Stderr, "sampe: %s\n", sampe)
	cmds = []Cmd{{sampe, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	dedup := work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/MarkDuplicates.jar " + "I=" + work["TMPDIR"] + "/SID" + sampleName + ".sorted.bam " + "O=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " + " M=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.metric " + " VALIDATION_STRINGENCY=SILENT CREATE_INDEX=true REMOVE_DUPLICATES=false TMP_DIR=" + work["TMPDIR"]

	fmt.Fprintf(os.Stderr, "dedup: %s\n", dedup)

	cmds = []Cmd{{dedup, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	flagstat := work["SAMTOOLS"] + " flagstat " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam  > " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam.flagstat"
	fmt.Fprintf(os.Stderr, "flagstat: %s\n", flagstat)

	cmds = []Cmd{{flagstat, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	hsmetrics := work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/CalculateHsMetrics.jar " + "I=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " + "O=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam.target_coverage " + " VALIDATION_STRINGENCY=SILENT " + " BI=" + work["TargetPicard"] + " TI=" + work["TargetPicard"] + " TMP_DIR=" + work["TMPDIR"]

	fmt.Fprintf(os.Stderr, "hsmetrics: %s\n", hsmetrics)
	cmds = []Cmd{{hsmetrics, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	realigner1 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		"    -known " + work["Mills"] +
		"    -known " + work["Indels"] +
		"    -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.intervals " +
		"    -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " +
		"    -L " + work["TargetBed"] +
		"    -nt 4 " +
		" -T  RealignerTargetCreator "

	fmt.Fprintf(os.Stderr, "realigner1: %s\n", realigner1)
	cmds = []Cmd{{realigner1, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	realigner2 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		"    -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " +
		"    -targetIntervals " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.intervals " +
		"    -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.bam " +
		" -T IndelRealigner "

	fmt.Fprintf(os.Stderr, "realigner2: %s\n", realigner2)
	cmds = []Cmd{{realigner2, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	bqsr1 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		" -knownSites " + work["Dbsnp"] +
		" -knownSites " + work["Mills"] +
		" -knownSites " + work["Indels"] +
		" -nct 4 " +
		work["UGinterval"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.grp " +
		" -T BaseRecalibrator "

	fmt.Fprintf(os.Stderr, "bqsr1: %s\n", bqsr1)

	cmds = []Cmd{{bqsr1, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	bqsr2 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -BQSR " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.grp " +
		" -T PrintReads "

	fmt.Fprintf(os.Stderr, "bqsr2: %s\n", bqsr2)
	cmds = []Cmd{{bqsr2, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	reduce := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" --out " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.reduced.bam " +
		" -T ReduceReads "

	fmt.Fprintf(os.Stderr, "reduce: %s\n", reduce)
	cmds = []Cmd{{reduce, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	hc := work["JAVAGATK"] +
		" -T  HaplotypeCaller " +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf  " +
		" -stand_call_conf 10 -stand_emit_conf 10 " +
		work["UGinterval"]

		//" -glm BOTH -stand_call_conf 10 -stand_emit_conf 10 -nt 4 -out_mode EMIT_ALL_CONFIDENT_SITES " +
	ug := work["JAVAGATK"] +
		" -T  UnifiedGenotyper " +
		" -glm BOTH -stand_call_conf 10 -stand_emit_conf 10 -nt 4 -out_mode EMIT_VARIANTS_ONLY " +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf  " +
		work["UGinterval"]

	fmt.Fprintf(os.Stderr, "ug: %s\n", ug)
	cmds = []Cmd{{ug, ""}, {hc, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	ugzip := work["BGZIP"] + " -c " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf  > " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf.gz " + " ; " + work["TABIX"] + " -p vcf " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf.gz "
	hczip := work["BGZIP"] + " -c " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf  > " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf.gz " + " ; " + work["TABIX"] + " -p vcf " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf.gz "

	fmt.Fprintf(os.Stderr, "ugzip: %s\n", ugzip)
	fmt.Fprintf(os.Stderr, "hczip: %s\n", hczip)
	cmds = []Cmd{{ugzip, ""}, {hczip, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	//copy back
	excludes := "--exclude dedup.bam --exclude indelrealigner.bam --exclude dedup.bai --exclude indelrealigner.bai --exclude sorted.bai --exclude sorted.bam --exclude .sai --exclude .vcf "
	copyBack := work["S3cmd"] + " " + work["S3cfg.output"] + " " + excludes + " sync " + work["TMPDIR"] + " " + work["BAMBUCKET"] + "/"
	fmt.Fprintf(os.Stderr, "copyBack: %s\n", copyBack)
	cmds = []Cmd{{copyBack, ""}}
	c = make(chan int, len(cmds))
	for i, v := range cmds {
		MyPrintf("cmds %d channel\n", i)
		go v.Do(c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	if !*dryrun {
		Finish()
	}
}
