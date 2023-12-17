package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

var cardsMap = map[rune]int{
	'J': 0,
	'1': 1,
	'2': 2,
	'3': 3,
	'4': 4,
	'5': 5,
	'6': 6,
	'7': 7,
	'8': 8,
	'9': 9,
	'T': 10,
	'Q': 12,
	'K': 13,
	'A': 14,
}

type hand struct {
	cards []int
	bid   int
	value int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getHand(line string) hand {
	var parts = strings.Split(line, " ")
	bid, err := strconv.Atoi(parts[1])
	check(err)

	var cardsRunes []rune = []rune(parts[0])
	var cards []int
	for _, card := range cardsRunes {
		cards = append(cards, cardsMap[card])
	}

	value := evaluateHand(cards)

	return hand{cards: cards, bid: bid, value: value}
}

func evaluateHand(cards []int) int {
	cardQty := make(map[int]int)
	maxCard := -1
	jokerCard := cardsMap['J']

	//get the card with the highest quantity
	for _, card := range cards {
		cardQty[card]++
		if card != jokerCard && cardQty[card] > cardQty[maxCard] {
			maxCard = card
		}
	}

	//joker replacement for part 2
	if cardQty[jokerCard] > 0 && maxCard != jokerCard {
		cardQty[maxCard] += cardQty[jokerCard]
		cardQty[jokerCard] = 0
	}

	tok := 0
	pairs := 0

	//assign a value to the hand
	for _, qty := range cardQty {
		switch qty {
		case 2:
			pairs++
		case 3:
			tok++
		case 4:
			return 6
		case 5:
			return 7
		}
	}

	if tok == 1 {
		if pairs == 1 {
			return 5
		}
		return 4
	}

	if pairs == 2 {
		return 3
	}

	if pairs == 1 {
		return 2
	}

	return 1
}

func getWinnings(hands []hand) int {
	var winnings int
	//order hands by value ascending
	sort.Slice(hands, func(i, j int) bool {

		//first by value
		if hands[i].value < hands[j].value {
			return true
		}

		if hands[i].value > hands[j].value {
			return false
		}

		//if value is the same, then by cards
		for k := 0; k < len(hands[i].cards); k++ {
			if hands[i].cards[k] < hands[j].cards[k] {
				return true
			}

			if hands[i].cards[k] > hands[j].cards[k] {
				return false
			}
		}

		return true
	})

	//multiply bid by position in array (rank)
	for i, hand := range hands {
		winnings += hand.bid * (i + 1)
	}

	return winnings
}

func processFile(f *os.File) int {

	scanner := bufio.NewScanner(f)
	var hands []hand

	for scanner.Scan() {
		line := scanner.Text()
		hands = append(hands, getHand(line))
	}

	err := scanner.Err()
	check(err)

	return getWinnings(hands)
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	sum := processFile(f)
	f.Close()

	fmt.Println(sum)
}
