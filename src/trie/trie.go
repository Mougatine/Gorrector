package trie

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math"
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
	Content   string
	Frequency int
	Distance  int
}

func (t *Trie) SearchCloseWords(word string, distance int) []Word {
	a := make([]Word, 1)

	fmt.Println("looking...", word, distance)
	for _, child := range t.Children {
		a = append(a, searchCloseWords(child, word, distance, 0, "")...)
	}
	fmt.Println(a)
	os.Exit(3)
	return a
}

func searchCloseWords(node *Trie, word string, maxDist int, curDist int, curr string) []Word {
	if word[0] != node.Char {
		curDist++
	}
	word = word[1:]

	if curDist > maxDist {
		return []Word{}
	}

	mdist := -1
	ans := make([]Word, 1)
	curr += node.Char

	if len(word)+curDist < maxDist

	if len(node.Children) == 0 {
		curDist += len(word)
	}

	if curDist <= maxDist && node.IsWord {
		ans = append(ans, Word{curr, node.Frequency, curDist})
	}

	if curDist+1 < maxDist && len(word) > 0 {
		suppression := searchCloseWords(node, word[1:], curDist+1, maxDist, curr)
		ans := append(ans, suppression...)

		for _, child := range node.Children {
			if len(word) > 0 && child.Char == string(word[0]) {
				mdist = 0
			} else {
				mdist = 1
			}

			substitution := searchCloseWords(child, word[1:], curDist+mdist, maxDist, curr)
			insertion := searchCloseWords(child, word, curDist+1, maxDist, curr)
			ans = append(ans, substitution...)
			ans = append(ans, insertion...)
		}
	}

	return ans
}

/*
Not working, from utard'slides.
*/
func computeDistance(node *Trie, word string, curDistance int, maxDistance int) int {
	if curDistance > maxDistance {
		return curDistance
	}
	res, mdist := -1, -1

	if len(node.Children) == 0 {
		res = len(word)
	}
	if curDistance+1 < maxDistance {
		suppression := computeDistance(node, word[1:], curDistance+1, maxDistance)
		res = myMin(res, suppression)
	}

	for _, child := range node.Children {
		if len(word) > 0 && child.Char == string(word[0]) {
			mdist = 0
		} else {
			mdist = 1
		}

		substitution := computeDistance(child, word[1:], curDistance+mdist, maxDistance)
		insertion := computeDistance(child, word, curDistance+1, maxDistance)
		res = myMin(res, substitution, insertion)
	}

	return res
}

func myMin(args ...int) int {
	v := math.MaxInt64
	for _, arg := range args {
		if v > arg {
			v = arg
		}
	}

	return v
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
