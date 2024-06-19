package main

import (
	"bufio"
	"fmt"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"
)

type Station struct {
	min  float64
	mean float64
	max  float64
}

func processFile(path string) (int, error) {
	file, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	stations := make(map[string]Station)

	for scanner.Scan() {
		parts := strings.SplitN(scanner.Text(), ";", 2)
		name := parts[0]
		measurement, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return 0, err
		}
		station := stations[name]
		if measurement > station.max {
			station.max = measurement
			stations[name] = station
		}
		if measurement < station.min {
			station.min = measurement
			stations[name] = station
		}
		station.mean = (station.min + station.max) / 2
		stations[name] = station
	}

	if err := scanner.Err(); err != nil {
		return 0, err
	}

	i := 1
	fmt.Print("{ ")
	for name, station := range stations {
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, math.Round(station.min*10)/10, math.Round(station.mean*10)/10, math.Round(station.max*10)/10)
		if i < len(stations) {
			fmt.Print(", ")
		}
		i++
	}
	fmt.Println(" }")

	return len(stations), nil
}

func main() {
	start := time.Now()

	numberOfStations, err := processFile("measurements.txt")

	if err != nil {
		log.Fatal(err)
	}

	end := time.Now()
	duration := end.Sub(start)

	fmt.Printf("Unique Stations: %d\n", numberOfStations)
	fmt.Printf("Execution time: %v\n", duration)
}
