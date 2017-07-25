package test

import (
	"testing"

	trie "../trie"
)

// BenchmarkCreateTrie gives the execution time and memory consumption of the CreateTrie method
func BenchmarkCreateTrie(b *testing.B) {
	b.ReportAllocs()
	trie.CreateTrie("../../words.txt")
}

// BenchmarkSaveTrie gives the execution time and memory consumption of the SaveTrie method
/*func BenchmarkSaveTrie(b *testing.B) {
	root, _ := trie.CreateTrie("../../words.txt")
	b.ResetTimer()
	trie.SaveTrie("dict.bin", root)
}*/
