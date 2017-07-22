package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"runtime"

	trie "../trie"
)

func main() {
	runtime.GOMAXPROCS(1)

	flag.Parse()
	if len(flag.Args()) != 1 {
		fmt.Println("Usage: ./TextMiningApp /path/to/compiled/dict.bin")
		os.Exit(134)
	}

	dictPath := flag.Arg(0)
	dict, err := trie.LoadTrie(dictPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		distance, word := fields[1], fields[2]
		dist, err := strconv.ParseUint(distance, 10, 8)
		if err != nil {
			panic("Error")
		}

		answers := dict.ExactSearch(word, uint8(dist))
		trie.PrettyPrint(answers)
	}
}
