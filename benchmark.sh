#!/usr/bin/env bash
# TODO Add a check for colors and emojis at some point

# --- Progress Spinner ---
spinner() {
    local pid=$1
    local message=${2:-"Working..."}
    local spin_chars="⠗⠗⠗⠗⠗⠗⠖⠖⠖⠖⠖⠖⠘⠘⠘⠘⠘⠘⠰⠰⠰⠰⠰⠰⠤⠤⠤⠤⠤⠤⠦⠦⠦⠦⠦⠦"
    echo -n "$message "
    while ps -p $pid > /dev/null; do
        for i in $(seq 0 $((${#spin_chars}-1))); do
            echo -ne "\033[0;36m${spin_chars:$i:1}\033[0m" # Blue spinner
            echo -ne "\033[D" # Move cursor back
            sleep 0.05
        done
    done
    echo -e "\033[0;32m✓ Done\033[0m" # Green checkmark
}

# --- Toolchain Selection ---
if command -v rg &> /dev/null; then GREP_CMD="rg"; else GREP_CMD="grep"; fi

# --- Prerequisite Checks ---
echo "Checking prerequisites..."
if ! command -v git &> /dev/null; then echo "❌ Error: 'git' not found."; exit 1; fi
if ! command -v go &> /dev/null; then echo "❌ Error: 'go' not found."; exit 1; fi
if ! command -v java &> /dev/null; then echo "❌ Error: 'java' not found."; exit 1; fi
if ! command -v bc &> /dev/null; then echo "❌ Error: 'bc' not found."; exit 1; fi
if ! java -version 2>&1 | $GREP_CMD -q 'version "21'; then echo "❌ Error: Java is not version 21."; exit 1; fi
if [ -z "$JAVA_HOME" ]; then echo "❌ Error: JAVA_HOME is not set."; exit 1; fi
echo "✅ Prerequisites met."

# --- Configuration ---
NUM_RUNS=5
GO_PROGRAM="main.go"
EXECUTABLE_NAME="1brc"
ORIGINAL_REPO="https://github.com/gunnarmorling/1brc"
ORIGINAL_NAME="1brc-repo"
DATA_FILE="./$ORIGINAL_NAME/measurements.txt"
SOLUTION_FILE="./solution.txt"
OUTPUT_FILE="output.txt"

# --- Build the Go program ---
echo "Building the Go program..."
go build -o $EXECUTABLE_NAME $GO_PROGRAM

# --- File Generation ---
if [ ! -f "$DATA_FILE" ]; then
    read -p "❓ Data file '$DATA_FILE' not found. Generate it now? (y/n) " -n 1 -r; echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        git clone --depth 1 $ORIGINAL_REPO $ORIGINAL_NAME &> /dev/null &
        spinner $! "Cloning 1BRC repository..."

        (cd "$ORIGINAL_NAME" && ./mvnw -q clean verify) &> /dev/null &
        spinner $! "Building the Java reference implementation (mvnw)..."
        
        echo "Generating 1 billion row data file (this will take a while)..."
        # We let this one run in the foreground to see its own progress updates
        (cd "$ORIGINAL_NAME" && ./create_measurements.sh 1000000000)
        echo "✅ Data file generated."
    else
        echo "Aborting."; exit 1
    fi
fi

if [ ! -f "$SOLUTION_FILE" ]; then
    read -p "❓ Solution file '$SOLUTION_FILE' not found. Generate it now? (y/n) " -n 1 -r; echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        (cd "$ORIGINAL_NAME" && ./calculate_average_baseline.sh > "../$SOLUTION_FILE") &
        spinner $! "Generating solution file with reference Java implementation..."
    else
        echo "Aborting. Cannot run benchmark without a solution file."; exit 1
    fi
fi

echo -e "\n--- Starting Benchmark (running $NUM_RUNS times) ---"

total_time="0"
for i in $(seq 1 $NUM_RUNS); do
    echo -e "\n---> Run $i of $NUM_RUNS <---"
    echo "Clearing OS page cache (requires sudo)..."
    sync
    echo 3 | sudo tee /proc/sys/vm/drop_caches > /dev/null
    
    echo -n "Starting Go program..."
    run_time=$( { /usr/bin/time -f "%e" ./$EXECUTABLE_NAME -f "$DATA_FILE" > "$OUTPUT_FILE"; } 2>&1 )
    echo " finished in ${run_time}s"
    
    total_time=$(echo "$total_time + $run_time" | bc)
    
    echo -n "Checking correctness... "
    if diff -q "$SOLUTION_FILE" "$OUTPUT_FILE" > /dev/null; then
        echo -e "\033[0;32m✓ CORRECT\033[0m"
    else
        echo -e "\033[0;31m❌ ERROR: Output is INCORRECT. Differences:\033[0m"
        diff "$SOLUTION_FILE" "$OUTPUT_FILE"
        exit 1 
    fi
done

average_time=$(echo "scale=4; $total_time / $NUM_RUNS" | bc)
echo -e "\n--- Benchmark Finished ---"
echo "Average execution time over $NUM_RUNS runs: ${average_time}s"

rm "$OUTPUT_FILE"
