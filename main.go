package main

import (
	"fmt"
	"os"
	"runtime"
	"slices"

	"github.com/edsrzf/mmap-go"
)

type Measurement struct {
	Min   int
	Max   int
	Sum   int64
	Count int
}

type MemoryChunk struct {
	start int
	end   int
}

func splitMemory(memory mmap.MMap, n int) []MemoryChunk {
	total := len(memory)
	chunkSize := total / n
	chunks := make([]MemoryChunk, n)

	chunks[0].start = 0
	for i := 1; i < n; i++ {
		for j := i * chunkSize; j < i*chunkSize+50; j++ {
			if memory[j] == '\n' {
				chunks[i-1].end = j
				chunks[i].start = j + 1
				break
			}
		}
	}
	chunks[n-1].end = total - 1
	return chunks
}

func readMemoryChannel(ch chan map[string]*Measurement, data mmap.MMap, start int, end int) {
	station := ""
	temperature := 0
	prev := start
	measurements := make(map[string]*Measurement)
	for i := start; i <= end; i++ {
		if data[i] == ';' {
			station = string(data[prev:i])
			temperature = 0
			i += 1
			negative := false

			for data[i] != '\n' {
				ch := data[i]
				if ch == '.' {
					i += 1
					continue
				}
				if ch == '-' {
					negative = true
					i += 1
					continue
				}
				ch -= '0'
				if ch > 9 {
					panic("Invalid character")
				}
				temperature = temperature*10 + int(ch)
				i += 1
			}

			if negative {
				temperature = -temperature
			}

			measurement := measurements[station]
			if measurement == nil {
				measurements[station] = &Measurement{
					Min:   temperature,
					Max:   temperature,
					Sum:   int64(temperature),
					Count: 1,
				}
			} else {
				measurement.Min = min(measurement.Min, temperature)
				measurement.Max = max(measurement.Max, temperature)
				measurement.Sum += int64(temperature)
				measurement.Count += 1
			}

			prev = i + 1
			station = ""
			temperature = 0
		}
	}
	ch <- measurements
}

func main() {
	maxGoroutines := min(runtime.NumCPU(), runtime.GOMAXPROCS(0))

	dataFile, err := os.Open("./1brc-repo/measurements.txt")
	if err != nil {
		panic(err)
	}
	defer dataFile.Close()

	data, err := mmap.Map(dataFile, mmap.RDONLY, 0)
	if err != nil {
		panic(err)
	}
	defer data.Unmap()

	chunks := splitMemory(data, maxGoroutines)
	totals := make(map[string]*Measurement)
	measurementChan := make(chan map[string]*Measurement)

	for i := 0; i < maxGoroutines; i++ {
		go readMemoryChannel(measurementChan, data, chunks[i].start, chunks[i].end)
	}

	for i := 0; i < maxGoroutines; i++ {
		measurements := <-measurementChan
		for station, measurement := range measurements {
			total := totals[station]
			if total == nil {
				totals[station] = measurement
			} else {
				total.Min = min(total.Min, measurement.Min)
				total.Max = max(total.Max, measurement.Max)
				total.Sum += measurement.Sum
				total.Count += measurement.Count
			}
		}
	}

	printResults(totals)
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
		mean := float64(measurement.Sum/10) / float64(measurement.Count)
		max := float64(measurement.Max) / 10
		min := float64(measurement.Min) / 10
		fmt.Printf("%s=%.1f/%.1f/%.1f", stationName, min, mean, max)
		if idx < len(stationNames)-1 {
			fmt.Printf(", ")
		}
	}
	fmt.Printf("}\n")
}
