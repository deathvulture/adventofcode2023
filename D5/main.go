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

func processFile(f *os.File) (int, int) {
	scanner := bufio.NewScanner(f)

	scanner.Scan()
	seedsString := scanner.Text()
	mappingStage := 0

	var seedToSoilMap []relationMap
	var soilToFertilizerMap []relationMap
	var fertilizerToWaterMap []relationMap
	var waterToLightMap []relationMap
	var lightToTemperatureMap []relationMap
	var temperatureTohumidityMap []relationMap
	var humidityToLocationMap []relationMap

	for scanner.Scan() {
		line := scanner.Text()
		if line == "" {
			mappingStage++
			//skip the title line
			_ = scanner.Scan()
			continue
		}

		switch mappingStage {
		case 1:
			seedToSoilMap = append(seedToSoilMap, createRelationMap(line))
		case 2:
			soilToFertilizerMap = append(soilToFertilizerMap, createRelationMap(line))
		case 3:
			fertilizerToWaterMap = append(fertilizerToWaterMap, createRelationMap(line))
		case 4:
			waterToLightMap = append(waterToLightMap, createRelationMap(line))
		case 5:
			lightToTemperatureMap = append(lightToTemperatureMap, createRelationMap(line))
		case 6:
			temperatureTohumidityMap = append(temperatureTohumidityMap, createRelationMap(line))
		case 7:
			humidityToLocationMap = append(humidityToLocationMap, createRelationMap(line))
		}
	}

	relationMaps := [][]relationMap{
		seedToSoilMap, soilToFertilizerMap, fertilizerToWaterMap,
		waterToLightMap, lightToTemperatureMap, temperatureTohumidityMap,
		humidityToLocationMap}

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

	return minimumLocation, 0
}

func main() {
	dir, err := os.Getwd()
	check(err)

	fileName := "input.txt"
	file := filepath.Join(dir, fileName)

	f, err := os.Open(file)
	check(err)

	result, result2 := processFile(f)
	f.Close()
	fmt.Println("minimum location:", result)
	fmt.Println("total cards:", result2)
}
