package main

import (
	"bufio"
	"fmt"
	"os"

	"../cards"
	"../scorer"
)

const numCards = 52

func main() {
	// countAllHands()
	generateCachedScores()
}

type countedPoints struct {
	hand int
	crib int
}

func generateCachedScores() {
	// map of lead card to "serialized" hand value to the points for the hand
	pointCache := make(map[int8]map[int32]countedPoints, 30)

	for l := 0; l < numCards; l++ {
		lead := cards.NewCardFromNumber(l)
		leadVal := cards.CardToInt8(lead)

		if _, ok := pointCache[leadVal]; !ok {
			pointCache[leadVal] = map[int32]countedPoints{}
		}
		for c1 := 0; c1 < numCards; c1++ {
			if c1 == l {
				continue
			}
			card1 := cards.NewCardFromNumber(c1)
			for c2 := c1 + 1; c2 < numCards; c2++ {
				if c2 == c1 || c2 == l {
					continue
				}
				card2 := cards.NewCardFromNumber(c2)
				for c3 := c2 + 1; c3 < numCards; c3++ {
					if c3 == c2 || c3 == c1 || c3 == l {
						continue
					}
					card3 := cards.NewCardFromNumber(c3)
					for c4 := c3 + 1; c4 < numCards; c4++ {
						if c4 == c3 || c4 == c2 || c4 == c1 || c4 == l {
							continue
						}
						card4 := cards.NewCardFromNumber(c4)

						hand := []cards.Card{card1, card2, card3, card4}
						handVal := cards.HandToInt32(hand)
						hPoints := scorer.HandPoints(lead, hand)
						cPoints := scorer.CribPoints(lead, hand)

						pointCache[leadVal][handVal] = countedPoints{
							hand: hPoints,
							crib: cPoints,
						}
					}
				}
			}
		}
	}

	fmt.Println(`finished calcing!!`)
	printPointCache(pointCache)
	fmt.Println(`printed a file which is the point cache!!`)
}

func printPointCache(pointCache map[int8]map[int32]countedPoints) {
	f, err := os.Create(`../scorer/cached_scores.g.go`)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	w := bufio.NewWriter(f)

	_, err = w.WriteString(`
package scorer

type countedPoints struct {
	hand int
	crib int
}

var pointCache = map[int8]map[int32]countedPoints{
`)
	if err != nil {
		panic(err)
	}
	w.Flush()
	for l := int8(0); l < numCards; l++ {
		_, err = w.WriteString(fmt.Sprintf("\t%d: map[int32]countedPoints{\n", l))
		if err != nil {
			panic(err)
		}
		w.Flush()

		for handVal, cp := range pointCache[l] {
			_, err = w.WriteString(fmt.Sprintf("\t\t%d: countedPoints{hand: %d, crib: %d},\n", handVal, cp.hand, cp.crib))
			if err != nil {
				panic(err)
			}
			w.Flush()
		}

		_, err = w.WriteString(fmt.Sprintf("\t},\n"))
		if err != nil {
			panic(err)
		}
		w.Flush()
	}

	_, err = w.WriteString("\n}\n")
	w.Flush()
}

func countAllHands() {
	scoreCounts := make(map[int]int, 30)
	cribScoreCounts := make(map[int]int, 30)

	for l := 0; l < numCards; l++ {
		lead := cards.NewCardFromNumber(l)
		for c1 := 0; c1 < numCards; c1++ {
			if c1 == l {
				continue
			}
			card1 := cards.NewCardFromNumber(c1)
			for c2 := c1 + 1; c2 < numCards; c2++ {
				if c2 == c1 || c2 == l {
					continue
				}
				card2 := cards.NewCardFromNumber(c2)
				for c3 := c2 + 1; c3 < numCards; c3++ {
					if c3 == c2 || c3 == c1 || c3 == l {
						continue
					}
					card3 := cards.NewCardFromNumber(c3)
					for c4 := c3 + 1; c4 < numCards; c4++ {
						if c4 == c3 || c4 == c2 || c4 == c1 || c4 == l {
							continue
						}
						card4 := cards.NewCardFromNumber(c4)

						hand := []cards.Card{card1, card2, card3, card4}
						hPoints := scorer.HandPoints(lead, hand)
						fmt.Printf("%v-%v-h:%d\n", lead, hand, hPoints)
						scoreCounts[hPoints] += 1

						cPoints := scorer.CribPoints(lead, hand)
						fmt.Printf("%v-%v-c:%d\n", lead, hand, cPoints)
						cribScoreCounts[cPoints] += 1
					}
				}
			}
		}
	}

	fmt.Println(`finished!`)
	fmt.Printf("the hand counts look like: %v\n", scoreCounts)
	fmt.Printf("the crib hand counts look like: %v\n", cribScoreCounts)

	for i := 0; i < 30; i++ {
		fmt.Printf("%d,%d,%d\n", i, scoreCounts[i], cribScoreCounts[i])
	}
}
