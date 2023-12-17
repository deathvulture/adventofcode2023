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

type relationMap struct {
	from     int
	to       int
	relation int
}

type mapRange struct {
	from int
	to   int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func getSeeds(str string, soilMap []relationMap) []string {
	seedsString := strings.TrimSpace(strings.Split(str, ":")[1])
	seeds := strings.Split(seedsString, " ")

	return seeds
}

func createRelationMap(str string) relationMap {
	splitedStr := strings.Split(str, " ")
	var rm relationMap
	from, err := strconv.Atoi(splitedStr[1])
	to, err := strconv.Atoi(splitedStr[2])
	destination, err := strconv.Atoi(splitedStr[0])
	check(err)
	rm.from = from
	rm.to = from + to
	rm.relation = destination - from

	return rm
}

func resolveRelationMap(rm relationMap, value int) int {
	result := value
	if value >= rm.from && value < rm.to {
		result = rm.relation + value
	}
	return result
}

func mapSeedToLocation(maps [][]relationMap, seed int) int {
	value := seed
	for _, relationMap := range maps {
		for _, rm := range relationMap {
			result := resolveRelationMap(rm, value)
			if result != value {
				value = result
				break
			}

		}
	}
	return value
}

func processFile(f *os.File) int {
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	seedsString := scanner.Text()
	mappingStage := -1

	var seedToSoilMap, soilToFertilizerMap, fertilizerToWaterMap, waterToLightMap,
		lightToTemperatureMap, temperatureTohumidityMap, humidityToLocationMap []relationMap

	relationMaps := [][]relationMap{seedToSoilMap, soilToFertilizerMap, fertilizerToWaterMap,
		waterToLightMap, lightToTemperatureMap, temperatureTohumidityMap, humidityToLocationMap}

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mappingStage++
			//skip the title line
			_ = scanner.Scan()
			continue
		}
		relationMaps[mappingStage] = append(relationMaps[mappingStage], createRelationMap(line))
	}

	seeds := getSeeds(seedsString, seedToSoilMap)

	var minimumLocation int
	for _, seedStr := range seeds {
		seed, err := strconv.Atoi(seedStr)
		check(err)
		location := mapSeedToLocation(relationMaps, seed)
		if location < minimumLocation || minimumLocation == 0 {
			minimumLocation = location
		}
	}
	err := scanner.Err()
	check(err)

	return minimumLocation
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	result := processFile(f)
	f.Close()
	fmt.Println("minimum location:", result)
}
