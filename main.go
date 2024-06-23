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

func processFile(path string, chunkSize int) int {
	ch := make(chan map[string]Station, chunkSize)
	var wg sync.WaitGroup

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	lines := make([]string, 0, chunkSize)

	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err.Error() != "EOF" {
				log.Fatal(err)
			}
			if len(lines) > 0 {
				wg.Add(1)
				go processChunk(lines, ch, &wg)
			}
			break
		}
		lines = append(lines, strings.TrimSpace(line))
		if len(lines) >= chunkSize {
			wg.Add(1)
			go processChunk(lines, ch, &wg)
			lines = make([]string, 0, chunkSize)
		}
	}

	go func() {
		wg.Wait()
		close(ch)
	}()

	aggregatedStations := make(map[string]Station, 10000)

	for stations := range ch {
		for name, station := range stations {
			if existing, found := aggregatedStations[name]; found {
				if station.min < existing.min {
					existing.min = station.min
				}
				if station.max > existing.max {
					existing.max = station.max
				}
				existing.mean = (existing.min + existing.max) / 2
				aggregatedStations[name] = existing
			} else {
				aggregatedStations[name] = station
			}
		}
	}

	i := 1
	fmt.Print("{ ")
	for name, station := range aggregatedStations {
		fmt.Printf("%s=%.1f/%.1f/%.1f", name, roundTowardsPositive(station.min), roundTowardsPositive(station.mean), roundTowardsPositive(station.max))
		if i < len(aggregatedStations) {
			fmt.Print(", ")
		}
		i++
	}
	fmt.Println(" }")

	return len(aggregatedStations)
}

func processChunk(chunk []string, ch chan<- map[string]Station, wg *sync.WaitGroup) {
	defer wg.Done()

	stations := make(map[string]Station, 10000)

	for _, line := range chunk {
		parts := strings.SplitN(line, ";", 2)
		name := parts[0]
		measurement, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			log.Fatal(err)
		}

		station, found := stations[name]
		if !found {
			station.min = measurement
			station.max = measurement
		} else {
			if measurement > station.max {
				station.max = measurement
			}
			if measurement < station.min {
				station.min = measurement
			}
		}

		station.mean = (station.min + station.max) / 2
		stations[name] = station
	}

	ch <- stations
}

func roundTowardsPositive(f float64) float64 {
	return math.Ceil(f*10) / 10
}

func main() {
	var file string
	var threads, chunkSize int

	flag.StringVar(&file, "file", "measurements.txt", "The file to take the measurements from")
	flag.IntVar(&threads, "threads", runtime.NumCPU(), "How many threads to spawn")
	flag.IntVar(&chunkSize, "chunkSize", 1000000, "How many lines per chunk")
	flag.Parse()

	runtime.GOMAXPROCS(threads)

	start := time.Now()

	uniqueStations := processFile(file, chunkSize)

	end := time.Now()
	duration := end.Sub(start)

	fmt.Printf("Unique stations: %d\n", uniqueStations)
	fmt.Printf("Execution time: %v\n", duration)
}
