//This program is an immitation of Peter Norvig's Poker program
//TODO: implement complete Poker ranking rules
package main

import "fmt"
import "sort"

type Card struct {
	Rank int // 1-10  11:jack  12:queen 13:king 14:Ace
	Suit int // 0:diamond 1:heart 2:spade 3:club
}

type Hand [5]Card

func (h Hand) String() string { return fmt.Sprintf("%v %v %v %v %v", h[0], h[1], h[2], h[3], h[4]) }

func Poker ( h []Hand ) Hand {
	return  Max( h )
}

func Max ( h []Hand ) Hand {
	return h[0]	
}

func (s *Hand) Len() int      { return len(*s) }
func (s *Hand) Swap(i, j int) { s[i], s[j] = s[j], s[i] }
type ByRank struct{ *Hand }
func (s ByRank) Less(i, j int) bool { return s.Hand[i].Rank < s.Hand[j].Rank }

/* Straight Flush: 0.0015%
 * 4 of a Kind	:	0.024 %
 * Full House	:	0.140 %
 * Flush	:	0.196 %
 * Straight	:	0.39%
 * 3 of a Kind	:	2.11%
 * 2 Pair	:	4.75%
 *   Pair	:	42.25%
 * High Card	:	50.11%
 */
func straight ( h Hand ) bool {
	//we should sort the h.Rank by ascending order first
	if h[4].Rank-h[3].Rank == 1 && h[3].Rank-h[2].Rank == 1  && h[2].Rank-h[1].Rank == 1  && h[1].Rank-h[0].Rank == 1 {
		return true
	}	
	return false
}

func flush ( h Hand ) bool {
	if h[0].Suit == h[1].Suit && h[1].Suit == h[2].Suit && h[2].Suit == h[3].Suit && h[3].Suit == h[4].Suit {
		return true
	} 

	return false	
}

func kind ( h Hand ) int {
	c := make(map[int]int)
	for _,v := range h {
		c[v.Rank]++
	}
	return len(c)
}

func test () {
	sf := Hand{{6,0},{7,0},{8,0},{9,0},{10,0}}	
	sf2 := Hand{{6,1},{7,1},{8,1},{9,1},{10,1}}	
	sf3 := Hand{{6,1},{7,2},{8,1},{10,1},{9,1}}	
	sf4 := Hand{{7,1},{6,0},{8,2},{10,3},{9,3}}	
	sf5 := Hand{{7,1},{7,0},{8,2},{7,3},{7,2}}	
	
	game := []Hand{ sf, sf2, sf3, sf4,sf5}
	for i,v := range game {
		fmt.Printf("Hand %d: %v ", i, v)
		if flush(v) {
			fmt.Printf(" It is a flush ")
		}
		                if straight(v) {
                        fmt.Printf(" straight \n")
                }
		sort.Sort( ByRank{&v} )
		fmt.Printf("sorted now: %v; kind =%d\n", v, kind (v) )
	}
	fmt.Printf ("Wining Hand is %v\n", Poker ( game ) )
}


func main() {
	test()
}

