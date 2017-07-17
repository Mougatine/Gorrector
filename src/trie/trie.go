package trie

import (
	"bufio"
	"encoding/gob"
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
	Content   string
	Frequency int
	Distance  int
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
// then on the frequency (descending) and on the lexicographic order
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
		if a.Content[i] < b.Content[i] {
			return true
		}
	}

	return false
}

func (t *Trie) SearchCloseWords(word string, distance int) []Word {
	wordList := []Word{}
	mdist := -1

	fmt.Println("looking...", word, distance)
	for _, child := range t.Children {
		if 1 < distance {
			computeDistance(child, word, 1, distance, "", &wordList, "del")
		}
		if child.Char == string(word[0]) {
			mdist = 0
		} else {
			mdist = 1
		}
		computeDistance(child, word[1:], mdist, distance, "", &wordList, "sub")
		computeDistance(child, word, 1, distance, "", &wordList, "insert")
	}

	return wordList
}

/*func searchCloseWords(node *Trie, word string, maxDist int, curDist int, curr string) []Word {
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
}*/

/*
Not working, from utard'slides.
*/
func computeDistance(node *Trie, word string, curDistance int, maxDistance int,
	currWord string, wordList *[]Word, step string) int {

	currWord = currWord + node.Char

	//fmt.Println("Step: " + step + " | word: " + word + " | currentWord: " + currWord + " | dist: " + strconv.Itoa(curDistance))
	if curDistance > maxDistance {
		return curDistance
	}
	res, mdist, substitution, insertion := 10, -1, -1, -1

	if node.IsWord {
		res = len(word)
	}

	if curDistance+1 <= maxDistance {
		wordVal := word
		if len(word) > 0 {
			wordVal = word[1:]
		}
		suppression := computeDistance(node, wordVal, curDistance+1, maxDistance,
			currWord, wordList, "del")
		res = myMin(res, suppression)
	}

	for _, child := range node.Children {
		//fmt.Println("Word value " + word)
		//if len(word) > 0 {
		//	fmt.Println("Child val: " + child.Char + " word val: " + string(word[0]))
		//	}

		if len(word) > 0 && child.Char == string(word[0]) {
			mdist = 0
		} else {
			mdist = 1
		}

		if len(word) > 0 && curDistance+mdist <= maxDistance {
			substitution = computeDistance(child, word[1:], curDistance+mdist, maxDistance,
				currWord, wordList, "sub")
		}

		if curDistance+1 <= maxDistance {
			insertion = computeDistance(child, word, curDistance+1, maxDistance,
				currWord, wordList, "insert")
		}
		res = myMin(res, substitution, insertion)
	}

	if len(word) == 0 && node.IsWord && res <= maxDistance {
		newWord := Word{currWord, node.Frequency, curDistance}
		//fmt.Println("Inserted word: " + newWord.Content)
		*wordList = append(*wordList, newWord)
	}
	return res
}

// PrettyPrint displays the sorted words in JSON format
func PrettyPrint(words []Word) {
	var orderedArray = Answer(words)
	sort.Stable(orderedArray)
	fmt.Print("[")
	for i := range words {
		fmt.Print("{\"word\":\"" + words[i].Content + "\",")
		fmt.Print("\"freq\":" + strconv.Itoa(words[i].Frequency) + ",")
		fmt.Print("\"distance\":" + strconv.Itoa(words[i].Distance) + "}")

		if i != len(words)-1 {
			fmt.Print(",")
		}
	}

	fmt.Print("]")
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

func myMin(args ...int) int {
	v := math.MaxInt64
	for _, arg := range args {
		if v > arg {
			v = arg
		}
	}

	return v
}
