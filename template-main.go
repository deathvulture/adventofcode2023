package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func processFile(f *os.File) int {

	scanner := bufio.NewScanner(f)
	var sum = 0

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
	}

	err := scanner.Err()
	check(err)

	return sum
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	//fileName = "testInput.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	result := processFile(f)
	f.Close()

	fmt.Println(result)
}
