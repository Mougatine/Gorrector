package main

import (
	"flag"

	"fmt"
	"os"

	trie "../trie"
)

type toto struct {
	t bool
}

func main() {
	flag.Parse()
	wordsPath, dictPath := flag.Arg(0), flag.Arg(1)

	root, err := trie.CreateTrie(wordsPath)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	root.CompactTrie()
	if err = trie.SaveTrie(dictPath, root); err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
}
