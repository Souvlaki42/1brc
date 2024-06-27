package main

import (
	"bytes"
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

type HashBucket struct {
	key   []byte
	value *Measurement
}

const (
	// FNV-1 64-bit constants from hash/fnv.
	offset64   = 14695981039346656037
	prime64    = 1099511628211
	numBuckets = 1 << 17 // 2^17
)

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

func readMemoryChunk(ch chan map[string]*Measurement, data mmap.MMap, start int, end int) {
	temperature := 0
	prev := start
	hash := uint64(offset64)
	buckets := make([]HashBucket, numBuckets) // hash buckets, linearly probed.
	active := 0                               // number of active buckets.

	for i := start; i <= end; i++ {
		hash ^= uint64(data[i]) // FNV-1a is XOR then *
		hash *= prime64

		if data[i] == ';' {
			station := data[prev:i]
			temperature = 0
			i++
			negative := false

			for data[i] != '\n' {
				ch := data[i]
				if ch == '.' {
					i++
					continue
				}
				if ch == '-' {
					negative = true
					i++
					continue
				}
				ch -= '0'
				if ch > 9 {
					panic("Invalid character")
				}
				temperature = temperature*10 + int(ch)
				i++
			}

			if negative {
				temperature = -temperature
			}

			// Go to correct bucket in hash table.
			hashIndex := int(hash & uint64(numBuckets-1))
			for {
				if buckets[hashIndex].key == nil {
					// Found empty slot, add new item (copying key).
					buckets[hashIndex] = HashBucket{
						key: station,
						value: &Measurement{
							Min:   temperature,
							Max:   temperature,
							Sum:   int64(temperature),
							Count: 1,
						},
					}
					active++
					if active > numBuckets/2 {
						panic("Too many items in hash table.")
					}
					break
				}
				if bytes.Equal(buckets[hashIndex].key, station) {
					// Found matching slot, add to existing value.
					s := buckets[hashIndex].value
					s.Min = min(s.Min, temperature)
					s.Max = max(s.Max, temperature)
					s.Sum += int64(temperature)
					s.Count++
					break
				}
				// Slot already holds another key, try next slot (linear probe).
				hashIndex++
				if hashIndex >= numBuckets {
					hashIndex = 0
				}
			}

			prev = i + 1
			temperature = 0
			hash = uint64(offset64)
		}
	}

	measurements := make(map[string]*Measurement)
	for _, bucket := range buckets {
		if bucket.key == nil {
			continue
		}
		measurements[string(bucket.key)] = bucket.value
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
		go readMemoryChunk(measurementChan, data, chunks[i].start, chunks[i].end)
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
