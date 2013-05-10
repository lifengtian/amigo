// Norvig Generap Problem Solver
// LFT May 2013
package main

import (
	"fmt"
)

// we don't have a set in Go yet; use map.
type State map[string]int

type Goals []string

// a single operator
type OP struct {
	Action       string
	Precondition []string
	Addlist      []string
	Dellist      []string
}

// all available operators
type OPs []OP

// return true if s
func MemberOf(s string, ss []string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func GPS(s State, g Goals, op OPs) {
	count := 0
	//fmt.Println("GPS: ", s, g, op)
	//every goal achieved
	for _, v := range g {
		if r := Achieve(s, v, op); r == true {
			count++
		}
	}
	if count == len(g) {
		fmt.Println("Solved")
	} else {
		fmt.Println("Not solved :(", count)
	}
}

func Achieve(s State, g string, op OPs) bool {
	count := 0
	//fmt.Println("Achieve: ", "g: ", g, "s: ", s, "op: ", op)
	if s[g] != 0 {
		return true
	}

	for _, n := range op {
		if Appropriate(g, n) {
			if Apply(n, s, op) {
				count++
			} else {
				//fmt.Println("Apply failed: ", n, s, op)
			}
		}
	}
	if count >= 1 {
		return true
	}

	return false

}

func Apply(op OP, s State, o OPs) bool {
	count := 0
	//fmt.Println("Apply : ", op, s)
	for _, v := range op.Precondition {
		if Achieve(s, v, o) {
			count++
		}
	}
	if count == len(op.Precondition) {
		fmt.Println(op.Action)

		//
		for _, v := range op.Dellist {
			delete(s, v)
		}

		for _, v := range op.Addlist {
			s[v] = 1
		}

		return true
	}

	return false
}

func Appropriate(g string, op OP) bool {
	//fmt.Println("Appropriate works: ", g, op)
	if MemberOf(g, op.Addlist) {
		return true
	}
	return false

}

func main() {
	testOP := OP{Action: "drive-son-to-school", Precondition: []string{"son-at-home", "car-works"}, Addlist: []string{"son-at-school"}, Dellist: []string{"son-at-home", "car-works"}}
	testOP2 := OP{Action: "make-car-works", Precondition: []string{"car-not-works"}, Addlist: []string{"car-works"}, Dellist: []string{"car-not-works"}}

	state := State{"son-at-home": 1, "car-not-works": 1}
	goals := Goals{"son-at-school"}
	fmt.Println(state, goals, testOP, testOP2)
	GPS(state, goals, OPs{testOP, testOP2})
	gatk := OPs{{Action: "fastq2sai", Precondition: []string{"fastq1","fastq2"}, Addlist: []string{"sai1","sai2"}, Dellist: []string{"fastq1", "fastq2"}},
			{Action: "sai2bam", Precondition: []string{"sai1","sai2"}, Addlist: []string{"bam"}, Dellist: []string{"sai1", "sai2"}},
			{Action: "dedup", Precondition: []string{"bam"}, Addlist: []string{"dedup"}, Dellist: []string{"bam"}},
			{Action: "realign", Precondition: []string{"dedup"}, Addlist: []string{"realign"}, Dellist: []string{"dedup"}},
			{Action: "recal", Precondition: []string{"realign"}, Addlist: []string{"recal"}, Dellist: []string{"realign"}},
			{Action: "recal2vcf", Precondition: []string{"recal"}, Addlist: []string{"vcf"}, Dellist: []string{"recal"}},}

	state = State{"fastq1":1, "fastq2":2}
	goals = Goals{"vcf"}
	GPS(state, goals, gatk)

}
