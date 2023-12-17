package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getDifferenceArray(values [][]int) [][]int {
	currentIndex := len(values) - 1
	currentValues := values[currentIndex]

	allZero := true
	for _, value := range currentValues {
		if value != 0 {
			allZero = false
			break
		}
	}

	//on the last row we add the zero values
	if allZero {
		currentValues = append(currentValues, 0)
		currentValues = append([]int{0}, currentValues...)
		values[len(values)-1] = currentValues
		return values
	}

	var differences []int
	for i := 1; i < len(currentValues); i++ {
		differences = append(differences, currentValues[i]-currentValues[i-1])
	}

	values = append(values, differences)
	values = getDifferenceArray(values)
	lastValues := values[currentIndex+1]
	predictedPostValue := currentValues[len(currentValues)-1] + lastValues[len(lastValues)-1]
	predictedPreValue := currentValues[0] - lastValues[0]
	currentValues = append(currentValues, predictedPostValue)
	currentValues = append([]int{predictedPreValue}, currentValues...)
	values[currentIndex] = currentValues
	return values
}

func predictValue(history []string) (int, int) {
	var values []int
	for _, historyValue := range history {
		value, err := strconv.Atoi(historyValue)
		check(err)
		values = append(values, value)
	}

	var differenceArray [][]int
	differenceArray = append(differenceArray, values)
	differenceArray = getDifferenceArray(differenceArray)
	postValue := differenceArray[0][len(differenceArray[0])-1]
	preValue := differenceArray[0][0]
	return preValue, postValue
}

func processFile(f *os.File) (int, int) {

	scanner := bufio.NewScanner(f)
	var preSum = 0
	var postSum = 0

	for scanner.Scan() {
		line := scanner.Text()
		history := strings.Split(line, " ")
		preValue, postValue := predictValue(history)
		preSum += preValue
		postSum += postValue
	}

	err := scanner.Err()
	check(err)

	return preSum, postSum
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	//fileName = "testInput.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	preValue, postValue := processFile(f)
	f.Close()

	fmt.Println(preValue, postValue)
}
