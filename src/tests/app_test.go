package test

import (
	"os"
	"testing"

	trie "../trie"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func BenchmarkSearchWords(b *testing.B) {
	root, _ := trie.CreateTrie("../../words.txt")
	root.CompactTrie()
	trie.SaveTrie("dict.bin", root)
	dict, _ := trie.LoadTrie("dict.bin")

	b.ResetTimer()
	// FIXME
	//dict.SearchCloseWords("test", 0)
}
