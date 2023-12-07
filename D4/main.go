package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func processLine(row string) int {

	numberPart := strings.Split(row, ":")[1]
	numbers := strings.Split(numberPart, "|")
	spaces := regexp.MustCompile(`\s+`)
	winningNumbers := spaces.Split(strings.TrimSpace(numbers[0]), -1)
	ownNumbers := spaces.Split(strings.TrimSpace(numbers[1]), -1)

	totalWinnings := 0

	//check if winningNumbers are in ownNumbers
	for _, winningNumber := range winningNumbers {
		for _, ownNumber := range ownNumbers {
			if winningNumber == ownNumber {
				totalWinnings++
			}
		}
	}
	return totalWinnings
}

func processCards(cards []int) int {
	sum := 0
	queue := []int{}
	// Add all cardsIndexes to queue
	for i := range cards {
		queue = append(queue, i)
	}

	// Process queue and add new cards to queue if corresponding card has winnings
	for len(queue) > 0 {
		sum++
		cardIndex := queue[0]
		queue = queue[1:]
		winnings := cards[cardIndex]
		if winnings > 0 {
			for i := 0; i < winnings; i++ {
				queue = append(queue, cardIndex+i+1)
			}
		}
	}
	return sum
}

func processFile(f *os.File) (int, int) {
	cards := []int{}
	scanner := bufio.NewScanner(f)
	var sum int = 0

	for scanner.Scan() {
		line := scanner.Text()
		winnings := processLine(line)
		cards = append(cards, winnings)

		// Calculate winnings
		if winnings > 1 {
			winnings = int(math.Pow(2, (float64(winnings) - 1)))
		}

		sum += int(winnings)
	}

	err := scanner.Err()
	check(err)
	totalCards := processCards(cards)
	return sum, totalCards
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	winnings, totalCards := processFile(f)
	f.Close()
	fmt.Println("winnings:", winnings)
	fmt.Println("total cards:", totalCards)
}
