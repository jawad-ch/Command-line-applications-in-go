package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	// Defining a boolean flag -l to count lines instead of words
	lines := flag.Bool("l", false, "Count lines")

	// Defining a boolean flag -b to count the number of bytes in addition to words and lines
	bytes := flag.Bool("b", false, "Count Bytes")

	// Parsing the flags provided by the user
	flag.Parse()

	fmt.Println(count(os.Stdin, *lines, *bytes))
}

func count(r io.Reader, countLines, bytes bool) int {
	// A scanner is used to read text from a Reader (such as files)
	scanner := bufio.NewScanner(r)

	// If the count lines flag is not set, we want to count words so we define
	//the scanner split type to words (default is split by lines)
	if !countLines {
		scanner.Split(bufio.ScanWords)
	}

	// Defining a counter
	wc := 0

	// Buffer to accumulate bytes
	var buffer []byte

	// For every word or line scanned, add 1 to the counter
	for scanner.Scan() {
		line := scanner.Bytes()
		buffer = append(buffer, line...)
		wc++
	}

	if bytes {
		return len(buffer)
	}

	// Return the total
	return wc
}
