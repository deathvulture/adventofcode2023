package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

type mapRange struct {
	from int
	to   int
}

type relationMap struct {
	mapRange mapRange
	relation int
}

func check(e error) {
	if e != nil {
		log.Fatal(e)
	}
}

func createSeedRanges(str string) []mapRange {
	splitedStr := strings.TrimSpace(strings.Split(str, ":")[1])
	seedsValues := strings.Split(splitedStr, " ")
	var mapRanges []mapRange
	for i, value := range seedsValues {
		if i%2 == 1 {
			continue
		}

		from, err := strconv.Atoi(value)
		check(err)

		quantity, err := strconv.Atoi(seedsValues[i+1])
		check(err)

		var mr mapRange = mapRange{from, from + quantity}
		mapRanges = append(mapRanges, mr)
	}
	fmt.Println(mapRanges)

	return mapRanges
}

func createRelationMap(str string) relationMap {
	splitedStr := strings.Split(str, " ")
	var rm relationMap

	from, err := strconv.Atoi(splitedStr[1])
	check(err)

	to, err := strconv.Atoi(splitedStr[2])
	check(err)

	destination, err := strconv.Atoi(splitedStr[0])
	check(err)

	rm.mapRange.from = from
	rm.mapRange.to = from + to
	rm.relation = destination - from

	return rm
}

func resolveRelationMaps(mr mapRange, relationMap []relationMap) []mapRange {
	var mapRanges []mapRange
	lastRangePosition := mr.from
	for i := 0; i < len(relationMap); i++ {
		from := lastRangePosition
		to := mr.to
		relMap := relationMap[i]

		//if we get to the last position of the mapRange, we can stop
		if lastRangePosition == mr.to {
			break
		}

		//if we are ahead of the current range we can skip it
		if lastRangePosition > relMap.mapRange.to {
			continue
		}

		// last position is in the current range
		if lastRangePosition >= relMap.mapRange.from {
			from = lastRangePosition + relMap.relation
		} else {
			if to > relMap.mapRange.from {
				to = relMap.mapRange.from + relMap.relation - 1
				lastRangePosition = relMap.mapRange.from
				i--
				mapRanges = append(mapRanges, mapRange{from, to})
				continue
			} else {
				to = relMap.mapRange.from - 1
				lastRangePosition = relMap.mapRange.from
				mapRanges = append(mapRanges, mapRange{from, to})
				i--
				continue
			}
		}

		if to < relMap.mapRange.to {
			to += relMap.relation
			lastRangePosition = mr.to
		} else {
			to = relMap.mapRange.to - 1 + relMap.relation
			lastRangePosition = relMap.mapRange.to
		}

		mapRanges = append(mapRanges, mapRange{from, to})
	}

	if len(mapRanges) == 0 {
		mapRanges = append(mapRanges, mr)
	}

	return mapRanges
}

func processMap(mapRanges []mapRange, relationMaps [][]relationMap) []mapRange {
	if len(relationMaps) == 0 {
		return mapRanges
	}

	var nextRanges []mapRange
	//check if the mapRange is contained in the relation map,
	for _, mapRange := range mapRanges {
		result := resolveRelationMaps(mapRange, relationMaps[0])
		nextRanges = append(nextRanges, result...)
	}

	if len(nextRanges) > 1 {
		sort.Slice(nextRanges, func(i, j int) bool {
			return nextRanges[i].from < nextRanges[j].from
		})
	}

	ranges := processMap(nextRanges, relationMaps[1:])
	return ranges
}

func processFile(f *os.File) int {
	scanner := bufio.NewScanner(f)
	scanner.Scan()
	seedsString := scanner.Text()
	seedranges := createSeedRanges(seedsString)
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

	for _, relationMap := range relationMaps {
		sort.Slice(relationMap, func(i, j int) bool {
			return relationMap[i].mapRange.from < relationMap[j].mapRange.from
		})
	}

	result := processMap(seedranges, relationMaps)
	sort.Slice(result, func(i, j int) bool {
		return result[i].from < result[j].from
	})

	err := scanner.Err()
	check(err)

	for _, location := range result {
		if location.from > 0 {
			return location.from
		}
	}

	return result[0].from
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
