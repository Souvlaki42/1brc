package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Station struct {
	min  float64
	mean float64
	max  float64
}

func processFile(path string, threads int) int {
	ch := make(chan map[string]Station)
	var wg sync.WaitGroup
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lines := make([]string, 0)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	length := len(lines)
	partSize := int(math.Ceil(float64(length) / float64(threads)))

	for i := 0; i < length; i += partSize {
		end := i + partSize
		if end > length {
			end = length
		}
		chunk := lines[i:end]
		wg.Add(1)
		go processChunk(chunk, ch, &wg)
	}

	i := 1
	uniqueStations := 0

	go func() {
		wg.Wait()
		close(ch)
	}()

	fmt.Print("{ ")
	for stations := range ch {
		uniqueStations += len(stations)
	}

	for stations := range ch {
		for name, station := range stations {
			fmt.Printf("%s=%.1f/%.1f/%.1f", name, math.Round(station.min*10)/10, math.Round(station.mean*10)/10, math.Round(station.max*10)/10)
			if i < uniqueStations {
				fmt.Print(", ")
			}
			i++
		}
		fmt.Println(" }")
	}

	return uniqueStations
}

func processChunk(chunk []string, ch chan<- map[string]Station, wg *sync.WaitGroup) {
	stations := make(map[string]Station)
	defer wg.Done()

	for _, line := range chunk {
		parts := strings.SplitN(line, ";", 2)
		name := parts[0]
		measurement, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Fatal(err)
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

	ch <- stations
}

func main() {
	var file string
	var threads int
	flag.StringVar(&file, "file", "measurements.txt", "The file to take the measurements from")
	flag.IntVar(&threads, "threads", runtime.NumCPU(), "How many threads to spawn")
	flag.Parse()

	runtime.GOMAXPROCS(threads)

	start := time.Now()

	numberOfStations := processFile(file, threads)

	end := time.Now()
	duration := end.Sub(start)

	fmt.Printf("Unique Stations: %d\n", numberOfStations)
	fmt.Printf("Execution time: %v\n", duration)
}
