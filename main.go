package main

import (
	"bufio"
	"fmt"
	"os"
	"slices"
	"strconv"
	"strings"
)

type Measurement struct {
	Min   float64
	Max   float64
	Sum   float64
	Count int
}

func main() {
	dataFile, err := os.Open("./1brc-repo/measurements.txt")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	measurements := make(map[string]*Measurement)

	fileScanner := bufio.NewScanner(dataFile)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		rawString := fileScanner.Text()
		stationName, temperatureStr, found := strings.Cut(rawString, ";")
		if !found {
			continue
		}
		temperature, err := strconv.ParseFloat(temperatureStr, 32)
		if err != nil {
			panic(err)
		}

		measurement := measurements[stationName]
		if measurement == nil {
			measurements[stationName] = &Measurement{
				Min:   temperature,
				Max:   temperature,
				Sum:   temperature,
				Count: 1,
			}
		} else {
			measurement.Min = min(measurement.Min, temperature)
			measurement.Max = max(measurement.Max, temperature)
			measurement.Sum += temperature
			measurement.Count += 1
		}
	}

	printResults(measurements)
}

func printResults(results map[string]*Measurement) {
	stationNames := make([]string, 0, len(results))
	for stationName := range results {
		stationNames = append(stationNames, stationName)
	}

	slices.Sort(stationNames)

	fmt.Printf("{")
	for idx, stationName := range stationNames {
		measurement := results[stationName]
		mean := measurement.Sum / float64(measurement.Count)
		fmt.Printf("%s=%.1f/%.1f/%.1f",
			stationName, measurement.Min, mean, measurement.Max)
		if idx < len(stationNames)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}\n")
}
