package trie

import (
	"encoding/gob"
	"os"
)

type Trie struct {
	isWord   bool
	children map[string]*Trie
	char     string
}

func CreateRootTrie() *Trie {
	return &Trie{false, make(map[string]*Trie), ""}
}

func (t *Trie) AddWord(letters string) {
	char := string(letters[0])

	if len(letters) == 1 {
		// End of recursion.
		t.isWord = true
	} else if _, ok := t.children[char]; ok {
		// Child path already exists.
		t.children[char].AddWord(letters[1:])
	} else {
		t.children[char] = &Trie{false, make(map[string]*Trie), char}
	}
}

func (t *Trie) FetchAllWords(prefix string) []string {
	words := make([]string, 1)
	current := prefix + t.char

	if t.isWord {
		words = append(words, current)
	}

	for _, child := range t.children {
		words = append(words, child.FetchAllWords(current)...)
	}

	return words
}

func (t *Trie) CompactTrie() {
	if len(t.children) == 1 {
		for child := range t.children {
			// Only looping once: one key, aka one child.
			if len(t.children[child].children) == 1 {
				for grandChild := range t.children[child].children {
					// Again only looping once.
					t.char = t.char + t.children[child].char
					t.children = map[string]*Trie{t.char: t.children[child].children[grandChild]}
				}
			}
		}
	}

	for child := range t.children {
		t.children[child].CompactTrie()
	}
}

func SaveTrie(path string, t *Trie) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(t)
}

func LoadTrie(path string) (*Trie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	decodedTrie := &Trie{}
	err = gob.NewDecoder(file).Decode(decodedTrie)
	return decodedTrie, err
}
