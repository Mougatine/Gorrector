package trie

import (
	"bufio"
	"encoding/gob"
	"os"
	"strconv"
	"strings"
)

type Trie struct {
	IsWord    bool
	Children  map[string]*Trie
	Char      string
	Frequency int
}

type Word struct {
	content   string
	Frequency int
}

func CreateRootTrie() *Trie {
	return &Trie{false, make(map[string]*Trie), "", 0}
}

func (t *Trie) AddWord(letters string, frequency int) {
	char := string(letters[0]) // Bug with root

	if _, ok := t.Children[char]; !ok {
		// Creates new child.
		t.Children[char] = &Trie{false, make(map[string]*Trie), char, 0}
	}

	if len(letters) == 1 {
		t.Children[char].IsWord = true
		t.Children[char].Frequency = frequency
	} else {
		t.Children[char].AddWord(letters[1:], frequency)
	}
}

func (t *Trie) FetchAllWords(prefix string) []string {
	words := make([]string, 1)
	current := prefix + t.Char

	if t.IsWord {
		words = append(words, current)
	}

	for _, child := range t.Children {
		words = append(words, child.FetchAllWords(current)...)
	}

	return words
}

func (t *Trie) CompactTrie() {
	if len(t.Children) == 1 {
		for child := range t.Children {
			// Only looping once: one key, aka one child.
			if len(t.Children[child].Children) == 1 {
				for grandChild := range t.Children[child].Children {
					// Again only looping once.
					t.Char = t.Char + t.Children[child].Char
					t.Children = map[string]*Trie{t.Char: t.Children[child].Children[grandChild]}
				}
			}
		}
	}

	for child := range t.Children {
		t.Children[child].CompactTrie()
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

func CreateTrie(path string) (*Trie, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	root := &Trie{false, make(map[string]*Trie), "", 0}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.Fields(scanner.Text())
		freq, err := strconv.Atoi(line[1])
		if err != nil {
			continue
		}
		root.AddWord(line[0], freq)
	}

	return root, nil
}
