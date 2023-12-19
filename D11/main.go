package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"strings"
)

const spaceMultiplier = 1000000

type galaxy struct {
	x int
	y int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

// we check for empty rows and columns to calculate the real position of each galaxy
func getGalaxyRealPositions(universe [][]string) ([]int, []int) {
	emptyCols := []int{}
	xPositions := []int{}
	yPositions := []int{}
	xPosition := 0
	yPosition := 0

	//aux array to keep track of empty columns
	for i := 0; i < len(universe); i++ {
		emptyCols = append(emptyCols, 1)
	}

	//check for empty rows and columns
	for _, row := range universe {
		isEmptyRow := true
		for j, col := range row {
			if col != "." {
				emptyCols[j] = 0
				isEmptyRow = false
			}
		}
		//add y position mapping
		if isEmptyRow {
			yPosition += spaceMultiplier
			yPositions = append(yPositions, yPosition)
		} else {
			yPosition++
			yPositions = append(yPositions, yPosition)
		}
	}

	//add x position mapping
	for _, col := range emptyCols {
		if col == 1 {
			xPosition += spaceMultiplier
			xPositions = append(xPositions, xPosition)
		} else {
			xPosition++
			xPositions = append(xPositions, xPosition)
		}
	}

	return xPositions, yPositions
}

func getGalaxies(universe [][]string) []galaxy {
	galaxies := []galaxy{}
	for i, row := range universe {
		for j, col := range row {
			if col == "#" {
				galaxies = append(galaxies, galaxy{j, i})
			}
		}
	}
	return galaxies
}

func getCombinations(galaxies []galaxy) [][]galaxy {
	combinations := [][]galaxy{}
	//get the combinations in reverse order to avoid duplicates
	for i, currentGalaxy := range galaxies {
		for j := i - 1; j >= 0; j-- {
			combinations = append(combinations, []galaxy{currentGalaxy, galaxies[j]})
		}
	}
	return combinations
}

func getDistance(combination []galaxy, xPositions, yPositions []int) int {
	g1 := combination[0]
	g2 := combination[1]
	xDistance := int(math.Abs(float64(xPositions[g1.x] - xPositions[g2.x])))
	yDistance := int(math.Abs(float64(yPositions[g1.y] - yPositions[g2.y])))
	return xDistance + yDistance
}

func processFile(f *os.File) int {
	scanner := bufio.NewScanner(f)
	var universe [][]string
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		universe = append(universe, line)
	}

	xPositions, yPositions := getGalaxyRealPositions(universe)
	galaxies := getGalaxies(universe)
	combinations := getCombinations(galaxies)
	var distances = 0
	for _, combination := range combinations {
		distances += getDistance(combination, xPositions, yPositions)
	}

	err := scanner.Err()
	check(err)

	return distances
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

	fmt.Println("Distance:", result)
}
