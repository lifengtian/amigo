//This program is an immitation of Peter Norvig's Poker program
//TODO: implement complete Poker ranking rules
package main

import "fmt"
import "sort"
import "math/rand"
import "time"

const cardNumbers=52
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

func kind ( h Hand ) (int, int) {
	c := make(map[int]int)
	max := 0
	for _,v := range h {
		c[v.Rank]++
		if c[v.Rank] > max {
			max = c[v.Rank]
		}
	}
	return len(c), max
}

func HandType( h Hand ) (result []int) {
	result = make([]int,6)
	c, m := kind( h )
	switch {
		case straight(h) && flush (h):
			result[0] = 1
			result[1] = h[4].Rank  //highest rank
			return 
		case straight(h) :
			result[0] = 5
			result[1] = h[4].Rank //highest rank
			return
		case flush (h):
			result[0] = 4
			result[1] = h[4].Rank //highest rank
			result[2] = h[3].Rank
			result[3] = h[2].Rank
			result[4] = h[1].Rank
			result[5] = h[0].Rank
			return
		case c == 2 && m == 4 : // 4 of a kind
			result[0] = 2
			if  h[0].Rank == h[1].Rank {
				result[1] = h[0].Rank
				result[2] = h[4].Rank
			} else {
				result[1] = h[4].Rank 
				result[2] = h[0].Rank
			}
			return
		case c == 2 && m == 3 : // Full house  33322
			result[0] = 3
			if h[0].Rank == h[1].Rank && h[1].Rank == h[2].Rank {
				result[1] = h[0].Rank
				result[2] = h[4].Rank
			} else {
				result[1] = h[4].Rank
				result[2] = h[0].Rank
			}
			return
		case c == 3 && m == 3 : // 3 of a kind 66645
			result[0] = 6
			result[1] = h[2].Rank
			switch {
			   case h[4].Rank == h[2].Rank:
				result[2] = h[1].Rank
				result[3] = h[0].Rank
			  case h[0].Rank == h[2].Rank:
				result[2] = h[4].Rank
				result[3] = h[3].Rank
			  default:
				result[2] = h[4].Rank
				result[3] = h[0].Rank
			}
			return
		case c == 3 && m == 2 : // 2 pairs 66554
			result[0] = 7
			switch {
				case h[4].Rank == h[3].Rank :
					result[1] = h[4].Rank
					if h[2].Rank == h[1].Rank {
						result[2] = h[2].Rank
						result[3] = h[0].Rank
					} else {
						result[2] = h[0].Rank
						result[3] = h[2].Rank
					}
				default:
					result[1] = h[3].Rank	
					result[2] = h[0].Rank
					result[3] = h[4].Rank
			}
			return
		case c == 4 && m == 2 : // pair   66543
			result[0] = 8
			switch {
				case h[0].Rank == h[1].Rank:
					result[1] = h[0].Rank
					result[2] = h[4].Rank
					result[3] = h[3].Rank
					result[4] = h[2].Rank
				case h[1].Rank == h[2].Rank:
					result[1] = h[1].Rank
					result[2] = h[4].Rank
					result[3] = h[3].Rank
					result[4] = h[0].Rank
				case h[2].Rank == h[3].Rank:
					result[1] = h[2].Rank
					result[2] = h[4].Rank
					result[3] = h[1].Rank
					result[4] = h[0].Rank
				case h[3].Rank == h[4].Rank:	
					result[1] = h[3].Rank
					result[2] = h[2].Rank
					result[3] = h[1].Rank
					result[4] = h[0].Rank
			}	
			return
		case c == 5 && m == 1 : // High Card
			result[0] = 9
					result[1] = h[4].Rank
					result[2] = h[3].Rank
					result[3] = h[2].Rank
					result[4] = h[1].Rank
					result[5] = h[0].Rank
			return
		default:
			fmt.Printf("Error! %v\n", h)
			return
	}


	
}

func test () {
	sf1 := Hand{{6,0},{5,0},{2,0},{3,0},{4,0}}	
	sf2 := Hand{{6,0},{6,1},{6,2},{6,3},{4,0}}	
	sf3 := Hand{{6,0},{6,1},{6,2},{4,0},{4,1}}	
	sf4 := Hand{{10,0},{8,0},{6,0},{5,0},{4,0}}	
	sf5 := Hand{{10,0},{9,1},{8,2},{7,0},{6,0}}	
	sf6 := Hand{{6,0},{6,1},{6,2},{4,0},{5,0}}	
	sf7 := Hand{{6,0},{6,1},{5,0},{5,1},{4,0}}	
	sf8 := Hand{{6,0},{6,1},{5,0},{4,1},{3,1}}	
	sf9 := Hand{{10,0},{8,1},{6,2},{5,3},{4,0}}	
	
	game := []Hand{ sf1, sf2, sf3, sf4,sf5, sf6, sf7, sf8, sf9}
	for i,v := range game {
		fmt.Printf("Hand %d: %v ", i, v)
		sort.Sort( ByRank{&v} )
		n, k := kind(v)
		fmt.Printf("sorted now: %v; no_ranks=%d max_runs_ranks=%d\tTransformed=%v\n", v, n, k, HandType(v) )
	}
	fmt.Printf ("Wining Hand is %v\n", Poker ( game ) )
}


func MakeHand () Hand {
	index := rand.Perm(cardNumbers)
	
	return Hand{AllCards[index[0]], AllCards[index[10]], AllCards[index[20]],AllCards[index[30]], AllCards[index[40]] }	
}

func MakeHand2 () Hand {
	all := make([]int, cardNumbers)
	result := make([]int, 5)
	count := 0
	//first draw
	index := rand.Intn( cardNumbers )
	all[index] = 1
	result[count] = index
	//second ... draw	
	for count = 1;count <= 4 ; count++ {
	index = rand.Intn( cardNumbers ) 
	// if pulled out try again and again
	for ;all[index] != 0; {
		index = rand.Intn( cardNumbers )
	}
	all[index] = 1
	result[count] = index
	}
	
	return Hand{AllCards[result[0]], AllCards[result[1]], AllCards[result[2]], AllCards[result[3]], AllCards[result[4]]}
}

func ReturnHand() {
	for m:=0;m<100000;m++ {
	h := MakeHand()
	 sort.Sort( ByRank{&h} )
	ht := HandType(h)
	//fmt.Printf("%v\t%v\t%v\n", h, ht, ht[0])
	fmt.Printf("%d\n", ht[0])
	}	
}


var AllCards [cardNumbers]Card

func shuffle ( c [cardNumbers]Card ) [cardNumbers]Card {
	p := rand.Perm(cardNumbers)
	var result [cardNumbers]Card
	for i,v := range p {
		result[i] = c[v]
	}
	return result	
}


func main() {
	rand.Seed ( time.Now().Unix() )
//	test()

	        for i:=2; i<=14; i++ {
                for j := 0; j<=3 ; j++ {
                        AllCards[j*13+(i-2)] = Card{i,j}
                }
        }
	
	//fmt.Printf("AllCards:\n%v\n", AllCards)
	AllCards = shuffle ( AllCards )
	//fmt.Printf("Shuffled AllCards:\n%v\n", AllCards)
	ReturnHand()
//	fmt.Printf( "Make a Hand\t%v\n", MakeHand() )
}

