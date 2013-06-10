package main

import (
	"fmt"
	"time"
)

type Ball struct{ hits int } //we record the number of hits in it

func main() {
	//the communication channel is called a 
	table := make(chan *Ball) // where Ball bounces to and fro the players
	go player("ping", table)
	go player("pong", table)
	go player("XXXX", table)

	table <- new(Ball) // game bebins
	time.Sleep(1 * time.Second)
	<-table
}

func player(name string, table chan *Ball) {
	for {
		ball := <-table
		ball.hits++
		fmt.Println(name, ball.hits)
		time.Sleep(100 * time.Millisecond)
		table <- ball
	}

}

