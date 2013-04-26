// Package utils implement helper functions
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net" // net.InterfaceAddrs
	"os"  // os.Stat
	"os/exec"
	"strings" // strings.HasPrefix
	"time"

//	"github.com/stathat/jconfig"
)

func MyPrintf(s string, i int) {
	fmt.Fprintf(os.Stderr, s, i)
}

// TODO: Check S3 object
func isS3object(bucketName string) (ok bool) {
	return true
}

// Check file status
func IsFile(filename string) bool {
	if strings.HasPrefix(filename, "s3:") {
		if isS3object(filename) {
			return true
		} else {
			return false
		}
	}

	_, err := os.Stat(filename)
	if err != nil {
		return false
	}
	return true
}

func Begin() {
	fmt.Fprintf(os.Stdout, "Begin at %v\n", time.Now())
	addrs, _ := net.InterfaceAddrs()
	fmt.Fprintf(os.Stdout, "%v\n", addrs)
}

func Finish() {
	fmt.Fprintf(os.Stdout, "End at %v\n", time.Now())
}

// Cmd holds content of a command
type Cmd struct {
	Command    string
	Parameters string
}

func (cmd Cmd) Do(dryrun bool, c chan int) {

	if dryrun {
		fmt.Printf("%s %v\n", cmd.Command, cmd.Parameters)
	}

	command := exec.Command("bash", "-c", cmd.Command+" "+cmd.Parameters)
	var out bytes.Buffer
	command.Stdout = &out
	var errout bytes.Buffer
	command.Stderr = &errout

	if !dryrun {
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

func RunIt(dryrun bool, cmds []Cmd) (ok bool) {
	c := make(chan int, len(cmds))
	for i, v := range cmds {
		if !dryrun {
			MyPrintf("cmds %d channel\n", i)
		}
		go v.Do(dryrun, c)
	}
	for i := 0; i < len(cmds); i++ {
		<-c
	}

	return true
}

// Workflow contains all the information about the workflow
// right now, it is nothing but a hash of files, parameters
// path to database files
// path to executables (dependencies)
// parameters
// Each step will check their prerequisite and send error message
//     back to the controller and wait for approval to continue
type Workflow map[string]string

func (w Workflow) Parse(jsonFile string) (err error) {
	var workflow_json []byte

	//read from file
	workflow_json, err = ioutil.ReadFile(jsonFile)

	err = json.Unmarshal(workflow_json, &w)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	return
}

// A comprehensive error checking
// PIPELINE directory
// Files in PIPELINE
// Commands
func (w Workflow) Check(verbose bool) (err error, ok bool) {
	ok = true

	for i, v := range w {
		if IsFile(v) {
			if verbose {
				fmt.Fprintf(os.Stdout, "%s\t%s\n", i, v)
			}
		} else {
			if verbose {
				fmt.Fprintf(os.Stderr, "NNNNNNNNOOOOOO %s\t%s\n", i, v)
			}

			f := w["PIPELINE"] + "/" + v
			//fmt.Fprintf(os.Stderr, "try %s\n", f)
			if !IsFile(f) {
				if verbose {
					fmt.Fprintf(os.Stderr, "still NOT THERE %s\n", f)
				}
			}
			//return false
			ok = false
		}
	}

	//err, ok = IsDir(w["PIPELINE"])

	return
}

func (w Workflow) String() {
	for i, v := range w {
		fmt.Printf("%s\t%s\n", i, v)
	}
}

// init from a workflow configure file
// check local tools, dbs versions
// download the required version
func (w Workflow) Init(jsonFileName string, sampleName string, dryrun bool) {
	err := w.Parse(jsonFileName)
	if err != nil {
		log.Fatalf("Error work.Parse %s", err)
	}

	w["TMPDIR"] = w["TMPDIR"] + "/" + sampleName
	w["sampleName"] = sampleName
	w["JAVAGATK"] = w["JAVABIN"] + " -Xmx6g -Djava.io.tmpdir=" + w["TMPDIR"] + " -jar " + w["GATKJAR"] + " -et NO_ET -K " + w["GATKKEY"]
	// Prepare interval, only make sense for exome
	w["UGinterval"] = " -L " + w["Exons"] + "  -L " + w["TargetBed"] + " --interval_padding 50 "

	w.GetToolsDBs(dryrun)

	verbose := true

	if err, ok := w.Check(!verbose); !ok {
		fmt.Fprintf(os.Stderr, "Check failed %v ", err)
	}

}

func (w Workflow) Finish() {
	fmt.Println("Finish")

}

// GetToolsDBs
func (work Workflow) GetToolsDBs(dryrun bool) {
	prepDirs := " sudo mkdir -p " + work["PIPELINE"] + " ; sudo chmod 777 " + work["PIPELINE"] + " ; sudo mkdir -p " + work["TMPDIR"] + " ; sudo chmod 777 " + work["TMPDIR"] +
		" ; sudo mkdir -p " + work["Fastq"] + " ;  sudo chmod 777 " + work["Fastq"] +
		" ; mkdir -p " + work["Fastq"] + "/" + work["sampleName"]

	RunIt(dryrun, []Cmd{{prepDirs, ""}})

	if !IsFile(work["PIPELINE"] + "/" + work["Tools"]) {
		tools := work["S3cmd"] + " " + work["S3cfg.pipeline"] + " get s3://cagpipelines/" + work["Tools"] + " " + work["PIPELINE"] + "/"
		RunIt(dryrun, []Cmd{{tools, ""}})

		prepTools := " cd " + work["PIPELINE"] + ";" + " tar jxvf " + work["Tools"]
		RunIt(dryrun, []Cmd{{prepTools, ""}})
	}

	if !IsFile(work["PIPELINE"] + "/" + work["DBs"]) {
		dbs := work["S3cmd"] + " " + work["S3cfg.pipeline"] + " get s3://cagpipelines/" + work["DBs"] + " " + work["PIPELINE"] + "/"
		RunIt(dryrun, []Cmd{{dbs, ""}})

		prepDBs := " cd " + work["PIPELINE"] + ";" + " tar jxvf " + work["DBs"]
		RunIt(dryrun, []Cmd{{prepDBs, ""}})

	}

}

func (work Workflow) Run(dryrun bool) {
	work.GetFastq(dryrun)
	work.DoFastqc(dryrun)
	work.DoBWAAlnPE(dryrun)
	work.DoBWASampe(dryrun)
	work.DoDedup(dryrun)
	work.DoRealigner(dryrun)
	work.DoBQSR(dryrun)
	work.DoReduceBAM(dryrun)
	work.DoUnifiedGenotyper(dryrun)
	work.DoHaplotypeCaller(dryrun)
	work.CopyBack(dryrun)
}

func (work Workflow) GetFastq(dryrun bool) {
	sampleName := work["sampleName"]
	r1 := work["Fastq"] + "/" + sampleName + "/pe_1.fq.gz"
	r2 := work["Fastq"] + "/" + sampleName + "/pe_2.fq.gz"

	s1 := work["TMPDIR"] + "/pe_1.sai"
	s2 := work["TMPDIR"] + "/pe_2.sai"

	//timestamp := time.Now()
	lb := work["LB"]
	pl := work["PL"]
	pu := work["PU"]
	cn := work["CN"]

	getFastq1 := work["S3cmd"] + " -c " + work["S3cfg.base"] + "  sync " + work["FASTQBUCKET"] + "/" + sampleName + "/pe_1.fq.gz " + work["Fastq"] + "/" + sampleName + "/"
	getFastq2 := work["S3cmd"] + " -c  " + work["S3cfg.base"] + "   sync " + work["FASTQBUCKET"] + "/" + sampleName + "/pe_2.fq.gz " + work["Fastq"] + "/" + sampleName + "/"

	fmt.Fprintln(os.Stderr, r1, r2, s1, s2, lb, pl, pu, cn, getFastq1, getFastq2)

	cmds := []Cmd{{getFastq1, ""}, {getFastq2, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoFastqc(dryrun bool) {
	sampleName := work["sampleName"]
	r1 := work["Fastq"] + "/" + sampleName + "/pe_1.fq.gz"
	r2 := work["Fastq"] + "/" + sampleName + "/pe_2.fq.gz"
	fastqc1 := work["FASTQC"] + "/fastqc -o " + work["TMPDIR"] + " " + r1
	fastqc2 := work["FASTQC"] + "/fastqc -o " + work["TMPDIR"] + " " + r2

	cmds := []Cmd{{fastqc1, ""}, {fastqc2, ""}}
	RunIt(dryrun, cmds)
}

// BWA align paired end reads
func (work Workflow) DoBWAAlnPE(dryrun bool) {
	sampleName := work["sampleName"]
	r1 := work["Fastq"] + "/" + sampleName + "/pe_1.fq.gz"
	r2 := work["Fastq"] + "/" + sampleName + "/pe_2.fq.gz"

	s1 := work["TMPDIR"] + "/pe_1.sai"
	s2 := work["TMPDIR"] + "/pe_2.sai"

	bwa1 := work["BWA"] + " aln  -t 4 -I " + work["Reference"] + " " + r1 + " > " + s1
	bwa2 := work["BWA"] + " aln  -t 4 -I " + work["Reference"] + " " + r2 + " > " + s2

	cmds := []Cmd{{bwa1, ""}, {bwa2, ""}}
	RunIt(dryrun, cmds)
}

// BWA sampe
func (work Workflow) DoBWASampe(dryrun bool) {
	sampleName := work["sampleName"]
	r1 := work["Fastq"] + "/" + sampleName + "/pe_1.fq.gz"
	r2 := work["Fastq"] + "/" + sampleName + "/pe_2.fq.gz"

	s1 := work["TMPDIR"] + "/pe_1.sai"
	s2 := work["TMPDIR"] + "/pe_2.sai"

	lb := work["LB"]
	pl := work["PL"]
	pu := work["PU"]
	cn := work["CN"]

	sampe := work["BWA"] + ` sampe -r "@RG\tID:` + sampleName + `\tSM:` + sampleName + `" ` + work["Reference"] + " " + s1 + " " + s2 + " " + r1 + " " + r2 +
		" | " +
		work["SAMTOOLS"] + " view -h -S  - " +
		" | " +
		work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/SortSam.jar I=/dev/stdin O=/dev/stdout SO=coordinate VALIDATION_STRINGENCY=SILENT QUIET=true COMPRESSION_LEVEL=0 TMP_DIR=" + work["TMPDIR"] +
		" | " +
		work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/AddOrReplaceReadGroups.jar   I=/dev/stdin O=" + work["TMPDIR"] + "/SID" + sampleName + ".sorted.bam " + "SO=coordinate RGID=" + sampleName +
		" RGLB=" + lb + " RGPL=" + pl + " RGPU=" + pu + " RGCN=" + cn + " RGSM=" + sampleName +
		" CREATE_INDEX=true VALIDATION_STRINGENCY=SILENT QUIET=true TMP_DIR=" + work["TMPDIR"]

	cmds := []Cmd{{sampe, ""}}
	RunIt(dryrun, cmds)

}

func (work Workflow) DoDedup(dryrun bool) {
	sampleName := work["sampleName"]
	dedup := work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/MarkDuplicates.jar " + "I=" + work["TMPDIR"] + "/SID" + sampleName + ".sorted.bam " + "O=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " + " M=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.metric " + " VALIDATION_STRINGENCY=SILENT CREATE_INDEX=true REMOVE_DUPLICATES=false TMP_DIR=" + work["TMPDIR"]
	cmds := []Cmd{{dedup, ""}}
	RunIt(dryrun, cmds)

	flagstat := work["SAMTOOLS"] + " flagstat " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam  > " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam.flagstat"

	hsmetrics := work["JAVABIN"] + " -Xmx4g  -jar " + work["PICARD"] + "/CalculateHsMetrics.jar " + "I=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " + "O=" + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam.target_coverage " + " VALIDATION_STRINGENCY=SILENT " + " BI=" + work["TargetPicard"] + " TI=" + work["TargetPicard"] + " TMP_DIR=" + work["TMPDIR"]

	cmds = []Cmd{{flagstat, ""}, {hsmetrics, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoRealigner(dryrun bool) {
	sampleName := work["sampleName"]
	realigner1 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		"    -known " + work["Mills"] +
		"    -known " + work["Indels"] +
		"    -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.intervals " +
		"    -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " +
		"    -L " + work["TargetBed"] +
		"    -nt 4 " +
		" -T  RealignerTargetCreator "

	cmds := []Cmd{{realigner1, ""}}
	RunIt(dryrun, cmds)
	realigner2 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		"    -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.bam " +
		"    -targetIntervals " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.intervals " +
		"    -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.bam " +
		" -T IndelRealigner "

	cmds = []Cmd{{realigner2, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoBQSR(dryrun bool) {
	sampleName := work["sampleName"]
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

	cmds := []Cmd{{bqsr1, ""}}
	RunIt(dryrun, cmds)

	bqsr2 := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -BQSR " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.grp " +
		" -T PrintReads "

	cmds = []Cmd{{bqsr2, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoReduceBAM(dryrun bool) {
	sampleName := work["sampleName"]
	reduce := work["JAVAGATK"] +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" --out " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.reduced.bam " +
		" -T ReduceReads "

	cmds := []Cmd{{reduce, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoUnifiedGenotyper(dryrun bool) {
	sampleName := work["sampleName"]
	ug := work["JAVAGATK"] +
		" -T  UnifiedGenotyper " +
		" -glm BOTH -stand_call_conf 10 -stand_emit_conf 10 -nt 4 -out_mode EMIT_VARIANTS_ONLY " +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf  " +
		work["UGinterval"]

	cmds := []Cmd{{ug, ""}}
	RunIt(dryrun, cmds)

	ugzip := work["BGZIP"] + " -c " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf  > " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf.gz " + " ; " + work["TABIX"] + " -p vcf " + work["TMPDIR"] + "/SID" + sampleName + ".UG.vcf.gz "

	cmds = []Cmd{{ugzip, ""}}
	RunIt(dryrun, cmds)
}

func (work Workflow) DoHaplotypeCaller(dryrun bool) {
	sampleName := work["sampleName"]
	hc := work["JAVAGATK"] +
		" -T  HaplotypeCaller " +
		" -R " + work["Reference"] +
		" -I " + work["TMPDIR"] + "/SID" + sampleName + ".dedup.indelrealigner.recal.bam " +
		" -o " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf  " +
		" -stand_call_conf 10 -stand_emit_conf 10 " +
		work["UGinterval"]

	cmds := []Cmd{{hc, ""}}
	RunIt(dryrun, cmds)

	hczip := work["BGZIP"] + " -c " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf  > " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf.gz " + " ; " + work["TABIX"] + " -p vcf " + work["TMPDIR"] + "/SID" + sampleName + ".HC.vcf.gz "

	cmds = []Cmd{{hczip, ""}}
	RunIt(dryrun, cmds)

}

func (work Workflow) CopyBack(dryrun bool) {
	//copy back
	excludes := "--exclude dedup.bam --exclude indelrealigner.bam --exclude dedup.bai --exclude indelrealigner.bai --exclude sorted.bai --exclude sorted.bam --exclude .sai --exclude .vcf "
	copyBack := work["S3cmd"] + " -c " + work["S3cfg.base"] + " " + excludes + " sync " + work["TMPDIR"] + " " + work["BAMBUCKET"] + "/"

	cmds := []Cmd{{copyBack, ""}}
	RunIt(dryrun, cmds)

}

/**************** LOW level functions ***********************/

// IsDir wraps around the os.Stat, os.IsNotExist, FileInfo.IsDir
// http://golang.org/src/pkg/os/types.go

func IsDir(dir string) (err error, ok bool) {

	fileinfo, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(os.Stderr, "%s not exist", dir)
		}
	} else {

		ok = fileinfo.IsDir()
	}
	return
}

func Trace(s string) string {
	fmt.Printf("%v Entering %s\n", time.Now(), s)
	return s
}

func Un(s string) {
	fmt.Printf("%v Leaving %s\n", time.Now(), s)
}
