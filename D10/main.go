package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var startDirection = vector2{1, 0}

const startingPipe = "L"

type vector2 struct {
	x int
	y int
}

type status struct {
	position  vector2
	direction vector2
	pipe      string
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func calcNextPoint(currentStatus status) (vector2, vector2) {
	nextPoint := currentStatus.position
	nextDirection := currentStatus.direction
	switch currentStatus.pipe {
	case "|":
		nextPoint.y -= currentStatus.direction.y
		nextDirection.x = 0
	case "-":
		nextPoint.x += currentStatus.direction.x * 1
		nextDirection.y = 0
	case "F":
		if currentStatus.direction.x == -1 {
			nextPoint.y += 1
			nextDirection = vector2{0, -1}
		} else {
			nextPoint.x += 1
			nextDirection = vector2{1, 0}
		}
	case "7":
		if currentStatus.direction.x == 1 {
			nextPoint.y += 1
			nextDirection = vector2{0, -1}
		} else {
			nextPoint.x -= 1
			nextDirection = vector2{-1, 0}
		}
	case "L":
		if currentStatus.direction.x == -1 {
			nextPoint.y -= 1
			nextDirection = vector2{0, 1}
		} else {
			nextPoint.x += 1
			nextDirection = vector2{1, 0}
		}
	case "J":
		if currentStatus.direction.x == 1 {
			nextPoint.y -= 1
			nextDirection = vector2{0, 1}
		} else {
			nextPoint.x -= 1
			nextDirection = vector2{-1, 0}
		}
	}

	return nextPoint, nextDirection
}

func isInside(pipeArray [][]string, position vector2) bool {
	intersections := 0

	//using the ray casting algorithm to check if the point is inside the area
	for i := position.x; i < len(pipeArray[0]); i++ {
		if pipeArray[position.y][i] != "." && pipeArray[position.y][i] != "-" && pipeArray[position.y][i] != "L" && pipeArray[position.y][i] != "J" {
			intersections++
		}
	}

	return intersections%2 != 0
}

func processArray(pipeArray [][]string, startPoint vector2) (int, int) {
	var maxDistance int

	//Copy the array to keep track of the original pipes
	pipeArrayCopy := make([][]string, len(pipeArray))
	for i, line := range pipeArray {
		pipeArrayCopy[i] = make([]string, len(line))
		copy(pipeArrayCopy[i], line)
	}

	var currentStatus = status{startPoint, startDirection, startingPipe}

	//get the pipe loop
	for currentStatus.pipe != "S" {
		if currentStatus.position.x != startPoint.x || currentStatus.position.y != startPoint.y {
			pipeArray[currentStatus.position.y][currentStatus.position.x] = "X"
		}
		nextPoint, nextDirection := calcNextPoint(currentStatus)
		currentStatus.position = nextPoint
		currentStatus.direction = nextDirection
		currentStatus.pipe = pipeArray[nextPoint.y][nextPoint.x]
		maxDistance++
	}

	//format the array for the area calculation
	for i, line := range pipeArrayCopy {
		for j, v := range line {
			if pipeArray[i][j] == "X" {
				pipeArray[i][j] = v
			} else if pipeArray[i][j] == "S" {
				pipeArray[i][j] = startingPipe

			} else {
				pipeArray[i][j] = "."
			}
		}
	}

	//calculate the area
	area := 0
	for i, line := range pipeArray {
		for j, v := range line {
			if v == "." {
				if isInside(pipeArray, vector2{j, i}) {
					pipeArray[i][j] = "I"
					area++
				} else {
					pipeArray[i][j] = "0"
				}
			}
		}
		fmt.Println(strings.Join(pipeArray[i], ""))
	}

	return maxDistance / 2, area
}

func processFile(f *os.File) (int, int) {
	scanner := bufio.NewScanner(f)

	var pipeArray [][]string
	for scanner.Scan() {
		line := strings.Split(scanner.Text(), "")
		pipeArray = append(pipeArray, line)
	}

	var startPoint = vector2{-1, -1}
	for i, line := range pipeArray {
		if startPoint.x != -1 {
			break
		}
		for j, v := range line {
			if v == "S" {
				startPoint.x = j
				startPoint.y = i
				break
			}
		}
	}
	distance, area := processArray(pipeArray, startPoint)
	err := scanner.Err()
	check(err)

	return distance, area
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	//fileName = "testInput.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	distance, area := processFile(f)
	f.Close()

	fmt.Println(distance, area)
}
