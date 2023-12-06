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

var possibleGears map[string][]int

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func isValidSymbol(r rune) bool {
	return !unicode.IsLetter(r) && !unicode.IsNumber(r) && r != '.'
}

// recives 3 rows to compare and returns the sum of the numbers in the middle row if they are surrounded by valid symbols
// we check the middle row against the other 2 and when a valid symbol is found, we search for the whole number
// for part 2 we will also add a reference to the possible gears and their power
func processLine(rows []string, lineNumber int) int {
	result := 0
	currentRow := rows[1]
	rowLength := len(currentRow)
	for i := 0; i < rowLength; i++ {
		if unicode.IsDigit(rune(currentRow[i])) {
			found := false
			//corrdinates of the * (possible gear)
			coordinates := ""
			for j := 0; j <= 2; j++ {
				if i > 0 && isValidSymbol(rune(rows[j][i-1])) {
					found = true
					if rune(rows[j][i-1]) == '*' {
						coordinates = strconv.Itoa(lineNumber-1+j) + "," + strconv.Itoa(i-1)
					}
					break
				}

				if isValidSymbol(rune(rows[j][i])) {
					found = true
					if rune(rows[j][i]) == '*' {
						coordinates = strconv.Itoa(lineNumber-1+j) + "," + strconv.Itoa(i)
					}
					break
				}

				if i < rowLength-1 && isValidSymbol(rune(rows[j][i+1])) {
					found = true
					if rune(rows[j][i+1]) == '*' {
						coordinates = strconv.Itoa(lineNumber-1+j) + "," + strconv.Itoa(i+1)
					}
					break
				}
			}
			if found {
				number := ""
				//search where the number start and end
				start := i
				end := i
				for start > 0 && unicode.IsDigit(rune(currentRow[start-1])) {
					start--
				}
				for end < rowLength-1 && unicode.IsDigit(rune(currentRow[end+1])) {
					end++
				}
				number = currentRow[start : end+1]
				fmt.Println("found", i, number)

				//set i to the end of the number so we don't process it again
				i = end

				num, err := strconv.Atoi(number)
				check(err)
				result += num
				//add power to the possible gears
				possibleGears[coordinates] = append(possibleGears[coordinates], num)
			}
		}
	}
	return result
}

func processFile(f *os.File) (int, int) {
	scanner := bufio.NewScanner(f)
	lineNumber := 0
	var sum = 0
	var compareRows = []string{"", "", ""}
	scanner.Scan()
	line := scanner.Text()
	filler := strings.Repeat(".", len(line))
	compareRows[1] = filler
	compareRows[2] = line

	for scanner.Scan() {
		line := scanner.Text()
		compareRows[0] = compareRows[1]
		compareRows[1] = compareRows[2]
		compareRows[2] = line
		result := processLine(compareRows, lineNumber)
		sum += result
		lineNumber++
	}

	err := scanner.Err()
	check(err)

	//process last line
	compareRows[0] = compareRows[1]
	compareRows[1] = compareRows[2]
	compareRows[2] = filler
	result := processLine(compareRows, lineNumber)
	sum += result
	gearPower := 0
	for _, gear := range possibleGears {
		if len(gear) == 2 {
			gearPower += gear[0] * gear[1]
		}
	}
	return sum, gearPower
}

func main() {
	possibleGears = make(map[string][]int)
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	sum, gearpower := processFile(f)
	f.Close()

	fmt.Println(sum)
	fmt.Println(gearpower)

}
