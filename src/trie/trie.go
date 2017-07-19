package trie

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"sort"
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
	Content   string `json:"word"`
	Frequency int    `json:"freq"`
	Distance  int    `json:"distance"`
}

// Answer implements sort.Interface for []Word
type Answer []Word

// Len implements the Len() method of the sort interface
func (ans Answer) Len() int {
	return len(ans)
}

// Swap implements the Swap() method of the sort interface
func (ans Answer) Swap(i, j int) {
	ans[i], ans[j] = ans[j], ans[i]
}

// Less implements the Less() method of the sort interface
// The comparison if based, first, on the distance (growing)
// then on the frequency (descending) and finally on the lexicographic order
func (ans Answer) Less(i, j int) bool {
	left, right := ans[i], ans[j]

	if left.Distance < right.Distance {
		return true
	} else if left.Distance == right.Distance {
		if left.Frequency > right.Frequency {
			return true
		} else if left.Frequency == right.Frequency {
			return lexicoOrder(left, right)
		}
	}

	return false
}

// lexicoOrder returns true if a is smaller than b
// using the lexicographic order
func lexicoOrder(a, b Word) bool {
	for i := 0; i < len(a.Content) && i < len(b.Content); i++ {
		if a.Content[i] != b.Content[i] {
			return a.Content[i] < b.Content[i]
		}
	}

	return true
}

// SearchCloseWords returns a list of words which
// Damereau-Levenstein distance from the word is at max equal to the distance parameter
func (t *Trie) SearchCloseWords(word string, distance int) []Word {
	wordList := []Word{}
	mdist := -1

	// Deletion
	if 1 <= distance {
		computeDistance(t, word[1:], 1, distance, "", &wordList, string(word[0]), "del", "del: "+strconv.Itoa(1))
	}

	for _, child := range t.Children {

		if child.Char == string(word[0]) {
			mdist = 0
		} else {
			mdist = 1
		}
		// Substitution
		computeDistance(child, word[1:], mdist, distance, "", &wordList, string(word[0]), "sub", "sub: "+strconv.Itoa(mdist))
		// Insertion
		computeDistance(child, word, 1, distance, "", &wordList, "", "insert", "insert: "+strconv.Itoa(1))
	}

	return wordList
}

// computeDistance calculates the distance from the query word to the word
// being constructed while visiting the trie
func computeDistance(node *Trie, word string, curDistance int, maxDistance int,
	currWord string, wordList *[]Word, deletedChar string, step string, path string) int {

	//fmt.Println("Step: " + step + " | word: " + word + " | node char " + node.Char + " | currentWord: " + currWord + " | dist: " + strconv.Itoa(curDistance))
	if curDistance > maxDistance {
		return curDistance
	}
	res, mdist, substitution, insertion, transposition := 10, -1, -1, -1, -1

	if node.IsWord {
		res = len(word)
	}

	if curDistance+1 <= maxDistance && len(word) > 0 {
		suppression := computeDistance(node, word[1:], curDistance+1, maxDistance,
			currWord, wordList, string(word[0]), "del", path+" del "+strconv.Itoa(curDistance+1))
		res = myMin(res, suppression)
	}

	for _, child := range node.Children {
		//fmt.Println("Word value "+ word)
		//fmt.Println("Child val: " + child.Char)

		if len(word) > 0 && child.Char == string(word[0]) {
			mdist = 0
		} else {
			mdist = 1
		}

		// Prevents useless recursive calls
		if len(word) > 0 && curDistance+mdist <= maxDistance {
			substitution = computeDistance(child, word[1:], curDistance+mdist, maxDistance,
				currWord+node.Char, wordList, string(word[0]), "sub", path+" sub: "+strconv.Itoa(curDistance+mdist))
		}

		// Prevents useless recursive calls
		if curDistance+1 <= maxDistance && (len(word) > 0 && child.Char != string(word[0]) || len(word) == 0) {
			insertion = computeDistance(child, word, curDistance+1, maxDistance,
				currWord+node.Char, wordList, "", "insert", path+" insert: "+strconv.Itoa(curDistance+1))
		}

		/*if len(word) > 0 && len(currWord) > 0 && deletedChar == child.Char && currWord[len(currWord)-1] == word[0] {
			transposition = computeDistance(child, word[1:], curDistance, maxDistance,
				currWord+node.Char, wordList, string(word[0]), "trans", path+" trans: "+strconv.Itoa(curDistance))
		}*/
		res = myMin(res, substitution, insertion, transposition)
	}

	if len(word) == 0 && node.IsWord && res <= maxDistance {
		newWord := Word{currWord + node.Char, node.Frequency, curDistance}
		//fmt.Println("Path: " + path + " | Inserted word: " + newWord.Content)
		*wordList = append(*wordList, newWord)
	}
	return res
}

// PrettyPrint displays the sorted words in JSON format
func PrettyPrint(words []Word) {
	var orderedArray = Answer(words)
	sort.Stable(orderedArray)
	jsonData, _ := json.Marshal(words)
	fmt.Println(string(jsonData))
}

// AddWord adds a word to the trie by creating a new node containing
// a character and indicating if it is a word or not
func (t *Trie) AddWord(letters string, frequency int) {
	char := string(letters[0]) // Bug with root

	if _, ok := t.Children[char]; !ok {
		// Creates new child.
		t.Children[char] = &Trie{false, make(map[string]*Trie), char, 0}
	}

	if len(letters) == 1 {
		t.Children[char].Frequency = frequency
		t.Children[char].IsWord = true
	} else {
		t.Children[char].AddWord(letters[1:], frequency)
	}
}

// CompactTrie merges a node with one child in order to gain some space
// The new node takes the old node's children as sons
func (t *Trie) CompactTrie() {
	if len(t.Children) == 1 {
		for child := range t.Children {
			t.Char = t.Char + child
			t.Children = t.Children[child].Children
		}
	}

	for child := range t.Children {
		t.Children[child].CompactTrie()
	}
}

// SaveTrie saves the Trie struct given a path. Uses the Gob serializer
func SaveTrie(path string, t *Trie) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	return gob.NewEncoder(file).Encode(t)
}

// LoadTrie loads a trie given a path
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

// CreateTrie creates the Trie structure given a text file's path
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

// myMin returns the minimum value in a slice
func myMin(args ...int) int {
	v := math.MaxInt64
	for _, arg := range args {
		if v > arg {
			v = arg
		}
	}

	return v
}
