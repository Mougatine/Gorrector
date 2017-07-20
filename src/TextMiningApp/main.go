package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	trie "../trie"
)

func main() {
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
		dist, err := strconv.Atoi(distance)
		if err != nil {
			continue
		}

		answers := dict.ExactSearch(word, dist)
		trie.PrettyPrint(answers)
	}
}
