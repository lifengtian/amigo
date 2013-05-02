package main


import (
	"math/rand"
	"fmt"
	"time"
)


func main(){
	now := time.Now()
	n := make ( []int, 10 )


	rand.Seed ( now.Unix() ) 
	for i:=0; i<=100000; i++ {
	 	n[ rand.Intn ( 10 ) ] ++
	}
	
	for i,v := range n {
		fmt.Printf("%d : %d\n", i, v )
	}

}