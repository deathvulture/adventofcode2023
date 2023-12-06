package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getSpelledNumber(s string) int {
	var numbers = []string{"one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
	for i, spelledNumber := range numbers {
		if strings.Contains(s, spelledNumber) {
			//quick fix to avoid missing a number that is before the spelled number
			if s[0] == spelledNumber[0] {
				return i + 1
			}
		}
	}
	return -1
}

func getNumber(s string) int {
	var firstNumber string
	var secondNumber string

	// Find first number
	for i := 0; i < len(s); i++ {
		if unicode.IsDigit(rune(s[i])) {
			firstNumber = string(s[i])
			break
		}

		//Calculate string length to check, to avoid out of range error
		var stringToCheckLength int
		if (i + 4) < len(s) {
			stringToCheckLength = 5
		} else if (len(s) - i) > 2 {
			stringToCheckLength = len(s) - i
		} else {
			continue
		}

		var spelledNumber = getSpelledNumber(s[i : i+stringToCheckLength])
		if spelledNumber > -1 {
			firstNumber = strconv.Itoa(spelledNumber)
			break
		}
	}

	// Find second number
	for i := len(s) - 1; i >= 0; i-- {
		var stringToCheckLength = 0
		if (i + 4) < len(s) {
			stringToCheckLength = 5
		} else if (len(s) - i) > 2 {
			stringToCheckLength = len(s) - i
		}

		if stringToCheckLength > 0 {
			var spelledNumber = getSpelledNumber(s[i : i+stringToCheckLength])
			if spelledNumber > -1 {
				secondNumber = strconv.Itoa(spelledNumber)
				break
			}
		}

		if unicode.IsDigit(rune(s[i])) {
			secondNumber = string(s[i])
			break
		}
	}

	//fmt.Println(firstNumber + secondNumber)

	result, err := strconv.Atoi(firstNumber + secondNumber)
	check(err)

	return result
}

func processFile(f *os.File) int {

	scanner := bufio.NewScanner(f)
	var sum = 0

	for scanner.Scan() {
		line := scanner.Text()
		sum += getNumber(line)
	}

	err := scanner.Err()
	check(err)

	return sum
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
