package main

import (
	"flag"

	"fmt"
	"os"

	trie "../trie"
)

func main() {
	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Fprintln(os.Stderr,
			"Provide the compiler with words.txt path and dict path.")
		os.Exit(1)
	}
	wordsPath, dictPath := flag.Arg(0), flag.Arg(1)

	root, err := trie.CreateTrie(wordsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	root.CompactTrie()
	if err = trie.SaveTrie(dictPath, root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
