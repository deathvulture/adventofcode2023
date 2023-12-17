package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

type race struct {
	time     int
	distance int
}

func getRaces(timesStr string, distancesStr string) ([]race, race) {
	re := regexp.MustCompile(`\s+`)
	times := re.Split(strings.TrimSpace(timesStr), -1)
	distances := re.Split(strings.TrimSpace(distancesStr), -1)

	var races []race
	for i, t := range times {
		distance, err := strconv.Atoi(distances[i])
		time, err := strconv.Atoi(t)
		check(err)
		races = append(races, race{time: time, distance: distance})
	}

	megaRaceTime, err := strconv.Atoi(strings.ReplaceAll(timesStr, " ", ""))
	megaRaceDistance, err := strconv.Atoi(strings.ReplaceAll(distancesStr, " ", ""))
	check(err)
	megarace := race{time: megaRaceTime, distance: megaRaceDistance}
	return races, megarace
}

func getRoots(a int, b int, c int) (float64, float64) {
	discriminant := b*b - 4*a*c
	if discriminant < 0 {
		return 0, 0
	}

	discriminantRoot := math.Sqrt(float64(discriminant))
	return (-float64(b) + discriminantRoot) / (2 * float64(a)), (-float64(b) - discriminantRoot) / (2 * float64(a))
}

func calcWinRacePosibities(race race) float64 {
	//in the example we got t = 7 and d = 9 so we need every value that return a distance greater than 9
	//we use this inequality x * (7-x) > 9 ---> x^2 - 7x + 9 > 0
	// we calculate the roots of the equation and then get the values between them
	max, min := (getRoots(1, -1*race.time, race.distance))
	if min != math.Floor(min) {
		min = math.Ceil(min)
	} else {
		min += 1
	}
	if max != math.Floor(max) {
		max = math.Floor(max)
	} else {
		max -= 1
	}

	//return the number of win posibilities
	return max - min + 1
}

func processFile(f *os.File) (int, int) {

	scanner := bufio.NewScanner(f)

	scanner.Scan()
	times := strings.Split(scanner.Text(), ":")[1]

	scanner.Scan()
	distances := strings.Split(scanner.Text(), ":")[1]

	races, megarace := getRaces(times, distances)
	megaRacewinChances := int(calcWinRacePosibities(megarace))

	winChances := 1
	for _, race := range races {
		winChances = winChances * int(calcWinRacePosibities(race))
	}

	fmt.Println(races)
	err := scanner.Err()
	check(err)

	return winChances, megaRacewinChances
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	chances, megaraceChances := processFile(f)
	f.Close()

	fmt.Println(chances, megaraceChances)
}
