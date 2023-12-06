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

const RED = 12
const GREEN = 13
const BLUE = 14

type roll struct {
	red   int
	green int
	blue  int
}

type game struct {
	id    int
	rolls []roll
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func validateGame(g game) bool {
	error := false
	for _, roll := range g.rolls {
		if roll.red > RED || roll.green > GREEN || roll.blue > BLUE {
			error = true
			break
		}
	}
	return !error
}

func getGameFromLine(line string) game {
	var g game
	gameParts := strings.Split(line, ":")

	//get id
	id, err := strconv.Atoi(gameParts[0][5:])
	check(err)
	g.id = id

	//get rolls
	rollsString := strings.Split(gameParts[1], ";")
	for _, rollData := range rollsString {
		colors := strings.Split(rollData, ",")
		roll := roll{}
		for _, color := range colors {
			//get color
			colorParts := strings.Split(strings.Trim(color, " "), " ")
			colorName := colorParts[1]
			colorValue := 0
			colorValue, err := strconv.Atoi(colorParts[0])
			check(err)
			//add color to roll
			switch colorName {
			case "red":
				roll.red = colorValue
			case "green":
				roll.green = colorValue
			case "blue":
				roll.blue = colorValue
			}
			//add roll to game
		}
		g.rolls = append(g.rolls, roll)
	}
	return g
}

func getCubeSetPower(g game) int {
	red, green, blue := 0, 0, 0

	for _, roll := range g.rolls {
		if roll.red > red {
			red = roll.red
		}
		if roll.green > green {
			green = roll.green
		}
		if roll.blue > blue {
			blue = roll.blue
		}
	}
	return red * green * blue
}

func processFile(f *os.File) (int, int) {
	scanner := bufio.NewScanner(f)
	var sum = 0
	var powerSum = 0

	for scanner.Scan() {
		line := scanner.Text()
		game := getGameFromLine(line)
		//first part
		if validateGame(game) {
			sum += game.id
		}

		//second part
		powerSum += getCubeSetPower(game)
	}

	err := scanner.Err()
	check(err)

	return sum, powerSum
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	sum, powerSum := processFile(f)
	f.Close()

	fmt.Println(sum)
	fmt.Println(powerSum)
}
