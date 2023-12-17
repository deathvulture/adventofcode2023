package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type paths struct {
	left  string
	right string
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return a / gcd(a, b) * b
}

func lcmOfMany(numbers []int) int {
	result := numbers[0]
	for _, num := range numbers[1:] {
		result = lcm(result, num)
	}
	return result
}

func getRoute(line string) (string, paths) {
	//black magic regex to get the key and paths
	re := regexp.MustCompile(`(\w+)\s*=\s*\((\w+),\s*(\w+)\)`)
	match := re.FindStringSubmatch(line)

	//key, left, right
	return match[1], paths{match[2], match[3]}
}

func traversePaths(routes map[string]paths, start string, finish string, instructions []string) int {
	current := start
	movements := 0

	for current != finish {
		for _, instruction := range instructions {
			movements++
			//get the path for the current key
			route := routes[current]

			//move to the next key
			if instruction == "L" {
				current = route.left

			} else {
				current = route.right
			}
			if current == finish {
				break
			}
		}
	}

	return movements
}

func fromAtoZ(routes map[string]paths, start string, instructions []string) int {
	current := start
	movements := 0
	finished := false
	for finished == false {
		for _, instruction := range instructions {
			movements++
			//get the path for the current key
			route := routes[current]

			//move to the next key
			if instruction == "L" {
				current = route.left

			} else {
				current = route.right
			}
			if current[2] == 'Z' {
				finished = true
				break
			}
		}
	}

	return movements
}

func processFile(f *os.File) (int, int) {

	scanner := bufio.NewScanner(f)

	//get direction instructions
	scanner.Scan()
	instructions := strings.Split(scanner.Text(), "")

	_ = instructions
	routes := make(map[string]paths)
	var startingPoints []string

	//blank line
	scanner.Scan()

	for scanner.Scan() {
		line := scanner.Text()
		key, route := getRoute(line)
		if key[2] == 'A' {
			startingPoints = append(startingPoints, key)
		}
		routes[key] = route
	}
	err := scanner.Err()
	check(err)

	//get all the movments from A to Z (Part 2)
	var movements []int
	for _, start := range startingPoints {
		currentMovements := fromAtoZ(routes, start, instructions)
		movements = append(movements, currentMovements)
	}

	part1 := traversePaths(routes, "AAA", "ZZZ", instructions)
	part2 := lcmOfMany(movements)
	return part1, part2
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	//fileName = "testInput.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	part1, part2 := processFile(f)
	f.Close()

	fmt.Println(part1, part2)
}
