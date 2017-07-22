package main

import (
	"flag"
	"runtime"

	"fmt"
	"os"
	"runtime/pprof"

	trie "../trie"
)

func main() {
	runtime.GOMAXPROCS(1)

	flag.Parse()
	if len(flag.Args()) != 2 {
		fmt.Println("Usage: ./TextMiningCompiler /path/to/word/freq.txt /path/to/output/dict.bin")
		os.Exit(134)
	}

	// Used for profiling
	f, err := os.Create("dump.profile")
	if err != nil {
		panic(err)
	}
	pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()

	wordsPath, dictPath := flag.Arg(0), flag.Arg(1)

	root, err := trie.CreateTrie(wordsPath)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	//root.CompactTrie()
	if err = trie.SaveTrie(dictPath, root); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
