package trie

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"sort"
	"strconv"
)

// Trie struct Represents a trie, Frequency is the word frequency
// Value is nil for the internal nodes,
type Trie struct {
	Value     []byte
	Children  []*Trie
	Frequency uint32
}

type WordList struct {
	words   []Word
	wordMap map[string]bool
}

type Word struct {
	Content   string `json:"word"`
	Frequency uint32 `json:"freq"`
	Distance  uint8  `json:"distance"`
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

// ExactSearch is used to find a word when the distance is equal to zero
/*func (t *Trie) ExactSearch(word string, distance uint8) []Word {

	var res []Word
	var match Word
	var child Trie
	node := *t

	byteWord := []byte(word)

	for _, val := range byteWord {
		_, prs := node.Children[val]

		if !prs {
			return []Word{}
		}

		node = child
	}

	match = Word{string(node.Value), node.Frequency, distance}

	return append(res, match)
}*/

// SearchCloseWords returns a list of words which
// Damereau-Levenstein distance from the word is at max equal to the distance parameter
/*func (t *Trie) SearchCloseWords(word string, distance int) []Word {
	wordList := WordList{[]Word{}, make(map[string]bool)}
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

	return wordList.words
}*/

// computeDistance calculates the distance from the query word to the word
// being constructed while visiting the trie
/*func computeDistance(node *Trie, word string, curDistance int, maxDistance int,
	currWord string, wordList *WordList, deletedChar string, step string, path string) int {

	//fmt.Println("Step: " + step + " | word: " + word + " | node char " + node.Char + " | currentWord: " + currWord + " | deleted: " + deletedChar + " | dist: " + strconv.Itoa(curDistance))

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

		// If a substitution with mdist = 0 is possible, don't try an insertion; Prevents duplication
		if curDistance+1 <= maxDistance && (len(word) > 0 && child.Char != string(word[0]) || len(word) == 0) {
			insertion = computeDistance(child, word, curDistance+1, maxDistance,
				currWord+node.Char, wordList, "", "insert", path+" insert: "+strconv.Itoa(curDistance+1))
		}

		if len(word) > 0 && len(currWord) > 0 && node.Char == string(word[0]) && child.Char == deletedChar &&
			node.Char != child.Char {
			transposition = computeDistance(child, word[1:], curDistance, maxDistance,
				currWord+node.Char, wordList, string(word[0]), "trans", path+" trans: "+strconv.Itoa(curDistance))
		}
		res = myMin(res, substitution, insertion, transposition)
	}

	if len(word) == 0 && node.IsWord && res <= maxDistance {
		newWord := Word{currWord + node.Char, node.Frequency, curDistance}
		//fmt.Println("Path: " + path + " | Inserted word: " + newWord.Content)
		_, prs := (*wordList).wordMap[newWord.Content]
		if !prs {
			(*wordList).words = append((*wordList).words, newWord)
			(*wordList).wordMap[newWord.Content] = true
		}
	}
	return res
}*/

// PrettyPrint displays the sorted words in JSON format
func PrettyPrint(words []Word) {
	var orderedArray = Answer(words)
	sort.Stable(orderedArray)
	jsonData, _ := json.Marshal(words)
	fmt.Println(string(jsonData))
}

// CompactTrie merges a node with one child in order to gain some space
// The new node takes the old node's children as sons
/*func (t *Trie) CompactTrie() {
	if len(t.Children) == 1 {
		for child := range t.Children {
			t.Char = t.Char + child
			t.Children = t.Children[child].Children
		}
	}

	for child := range t.Children {
		t.Children[child].CompactTrie()
	}
}*/

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
	var line [][]byte
	finished := false

	data, err := ioutil.ReadFile(path)

	if err != nil {
		panic("Error while opening the source file")
	}

	root := &Trie{nil, nil, 0}
	delim := []byte("	")

	for !finished {
		n, raw := readLine(data)
		if n > 1 {
			data = data[n:]
			line = bytes.Split(raw, delim)
			freq, _ := strconv.ParseUint(string(line[1]), 10, 32)
			root.AddWord(line[0], uint32(freq))
		} else {
			finished = true
		}
	}

	return root, nil
}

// AddWord adds a word to the trie by creating a new node.
// A node has:
// 	* `Value`: A suffix value, to get the whole word we have to sum the precedent value.
//  * `Children`: A list of *Trie.
//  * `Frequency`: A frequence value. If the frequence is equal to 0, the node doesn't contains a word.
func (t *Trie) AddWord(word []byte, frequency uint32) {
	node := t
	hasInserted := false

	for {
		for i, child := range node.Children {
			hasInserted = false

			prefix := getCommonPrefix(child.Value, word)
			if prefix == 0 { // No prefix in common.
				continue
			}

			// Insertion of an intermediary node called 'newChild'.
			child.Value = child.Value[prefix:]
			newChild := &Trie{word[0:prefix], []*Trie{child}, 0}
			node.Children[i] = newChild

			node = newChild
			word = word[prefix:]
			hasInserted = true
			break
		}

		// No prefix nodes have been found, thus we are creating a final node.
		if !hasInserted {
			child := &Trie{word, nil, frequency}
			node.Children = append(node.Children, child)
			break
		}
	}
}

func getCommonPrefix(w1, w2 []byte) int {
	var minLength int
	if len(w1) > len(w2) {
		minLength = len(w2)
	} else {
		minLength = len(w1)
	}

	for i := 0; i < minLength; i++ {
		if w1[i] != w2[i] {
			return i
		}
	}

	return minLength
}

func readLine(buf []byte) (int, []byte) {
	var res []byte
	var i int
	sep := byte('\n')

	for i = 0; i < len(buf); i++ {
		if buf[i] == sep {
			res = buf[:i]
			return i + 1, res
		}
	}
	return i + 1, res
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
