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

func DLDistance(str1 string, str2 string) int {
	len1, len2 := len(str1), len(str2)
	inf := len1 + len2
	d := make([][]int, len1+2)
	lastRow := make(map[string]int)

	for i := range d {
		d[i] = make([]int, len2+2)
	}

	// Fill array
	for i := 0; i < len1+2; i++ {
		d[i][0] = inf
		d[i][1] = i - 1
	}

	for j := 0; j < len2+2; j++ {
		d[0][j] = inf
		d[1][j] = j - 1
	}

	for row := 2; row < len1+2; row++ {

		// Current character in a
		chA := str1[row-2]
		lastMatchCol := 1

		for col := 2; col < len2+2; col++ {
			chB := str2[col-2]
			lastMatchRow := lastRow[string(chB)]
			if lastMatchRow == 0 {
				lastMatchRow = 1
			}
			var cost = 0
			if chA != chB {
				cost = 1
			}
			substitution := d[row-1][col-1] + cost
			addition := d[row][col-1] + 1
			deletion := d[row-1][col] + 1
			transposition := d[lastMatchRow-1][lastMatchCol-1] + (row - lastMatchRow - 1) + 1 + (col - lastMatchCol - 1)
			minDist := sliceMin([]int{substitution, addition, deletion, transposition})
			d[row][col] = minDist

			if cost == 0 {
				lastMatchCol = col
			}
		}

		lastRow[string(chA)] = row
	}

	return d[len1][len2]

}
