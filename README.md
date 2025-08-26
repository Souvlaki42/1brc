# 1 Billion Row Challenge in Go

This is my implementation for the 1 Billion Row Challenge. I know the competition was for Java and has ended, but I wanted to join the fun. My implementation is in Golang, as I am more familiar with it than I'm with Java.

## Results

My best average run on my 4-core, 16GB RAM machine is **11s**.

## Benchmarking

This repository includes an automated benchmarking script that handles setup, execution, and correctness checking.

**Prerequisites:**
*   Go (1.22+)
*   Java (21+) with `JAVA_HOME` set
*   Git, `sudo` access (for clearing caches), and `bc`
*   Optionaly, ripgrep as a faster replacement for grep

**To run the benchmark:**

```bash
# Make the script executable first
chmod +x benchmark.sh

# Run it!
./benchmark.sh
```

The script will automatically:
1.  Check for all required dependencies.
2.  Build the Go program.
3.  If necessary, clone the original 1BRC repository to generate the `measurements.txt` data file and the `solution.txt` file using the baseline Java implementation.
4.  Run your Go implementation multiple times, clearing the OS file cache before each run for consistent results.
5.  Check the output against the correct solution file.
6.  Calculate the average time between all of the runs.

## Manual Development

If you want to run the program manually for development or debugging:

```bash
# Run with default settings
go run main.go

# Specify a different measurements file
go run main.go -f path/to/your/measurements.txt
```

## Constraints

*   Temperatures are between -99.9 and 99.9 with exactly one fractional digit.
*   Station names are ≤ 100 bytes, with at most 10,000 unique stations.
*   Rounding follows IEEE 754 semantics ('round toward positive').

## Resources

A collection of useful links and articles related to the 1BRC.
*   [The One Billion Row Challenge](https://www.morling.dev/blog/one-billion-row-challenge/)
*   [A fun exploration of how quickly 1B rows from a text file can be aggregated with Java](https://github.com/gunnarmorling/1brc)
*   [The One Billion Row Challenge in Go: from 1m45s to 3.4s in nine solutions](https://benhoyt.com/writings/go-1brc/)
*   [Solution to One Billion Rows Challenge in Golang.](https://github.com/shraddhaag/1brc/)
*   [Why should the Java folk have all the fun?!](https://rmoff.net/2024/01/03/1%EF%B8%8F%E2%83%A3%EF%B8%8F-1brc-in-sql-with-duckdb/)
*   [One Billion Row Challenge in Go](https://mrkaran.dev/posts/1brc/)
*   [The One Billion Row Challenge Shows That Java Can Process a One Billion Rows File in Two Seconds](https://www.infoq.com/news/2024/01/1brc-fast-java-processing/)
*   [1 billion rows challenge in PostgreSQL and ClickHouse](https://ftisiot.net/posts/1brows/)
*   [1 BILLION row challenge in Go - 2.5 Seconds!](https://youtu.be/O1IFQav9FQg)
*   [I Parsed 1 Billion Rows Of Text (It Sucked)](https://youtu.be/e_9ziFKcEhw)

## Unlicense

This project is released into the public domain. See the [UNLICENSE](UNLICENSE) file for details.
