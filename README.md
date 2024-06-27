# 1 Billion Row Challenge

This is my implementation for the 1 Billion Row Challenge. \
I know the competition was for Java initially and that it have ended, but I wanted to join the fun. \
My implementation is in Golang/Go because, I'm much more familiar with it than I'm with Java.

## Setup
- Install java 21:
```bash
sudo pacman -S jdk21-openjdk
```
- Clone or fork the repo:
```bash
git clone https://github.com/gunnarmorling/1brc.git 1brc-repo
```
- Go inside the repo and build the test data generator:
```bash
cd ./1brc-repo
./mvnw clean verify
```
- Generate the test data:
```bash
# You can pass any number of bytes
# The 1B row file is about 14GB in size
./create_measurements.sh 1000000000
```
- Run the program

## Resources
[https://benhoyt.com/writings/go-1brc/](https://benhoyt.com/writings/go-1brc/)
[https://youtu.be/O1IFQav9FQg?si=uBaalVeGkevBOWT4](https://youtu.be/O1IFQav9FQg?si=uBaalVeGkevBOWT4) \
[https://youtu.be/e_9ziFKcEhw?si=GmluAFpm5fslQdvl](https://youtu.be/e_9ziFKcEhw?si=GmluAFpm5fslQdvl) \
[https://github.com/gunnarmorling/1brc](https://github.com/gunnarmorling/1brc) \
[https://github.com/shraddhaag/1brc/](https://github.com/shraddhaag/1brc/) \
[https://rmoff.net/2024/01/03/1%EF%B8%8F%E2%83%A3%EF%B8%8F-1brc-in-sql-with-duckdb/](https://rmoff.net/2024/01/03/1%EF%B8%8F%E2%83%A3%EF%B8%8F-1brc-in-sql-with-duckdb/) \
[https://mrkaran.dev/posts/1brc/](https://mrkaran.dev/posts/1brc/) \
[https://www.morling.dev/blog/one-billion-row-challenge/](https://www.morling.dev/blog/one-billion-row-challenge/) \
[https://www.infoq.com/news/2024/01/1brc-fast-java-processing/](https://www.infoq.com/news/2024/01/1brc-fast-java-processing/) \
[https://ftisiot.net/posts/1brows/](https://ftisiot.net/posts/1brows/)

## Constraints

The temperatures are between -99.9 and 99.9 and they have exactly one fractional digit. \
Rounding should be done using the semantics of IEEE 754 with a rounding direction 'round toward positive'. \
Each station name is less than or equal to 100 bytes long while they are at most 10000 unique stations in the file.

## Results

Fastest run of my current program is 10s for the 1 billion lines file.
