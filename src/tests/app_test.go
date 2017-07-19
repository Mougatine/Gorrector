package test

import (
	"bytes"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"testing"

	trie "../trie"
)

var testInputs = []struct {
	word string
	dist int
}{
	{"test", 0},
	{"cecci", 0},
	{"bonjour", 1},
	{"chouquette", 1},
}

func TestOutputs(t *testing.T) {
	// Creates the two binary files
	exec.Command("bash", "-c", "cd ../../ && make").Run()
	exec.Command("bash", "-c", "./../../TextMiningCompiler ../../words.txt ../../dict.bin").Run()

	t.Run("Output", func(t *testing.T) {
		appPath := "./../../TextMiningApp ../../dict.bin"
		refPath := "./../../ref/TextMiningApp ../../ref/dict.bin"

		for _, tt := range testInputs {
			cmdString := "echo \"approx " + strconv.Itoa(tt.dist) + " " + tt.word + "\"" + " | "
			out, _ := exec.Command("bash", "-c", cmdString+refPath).Output()
			myOut, _ := exec.Command("bash", "-c", cmdString+appPath).Output()
			if !bytes.Equal(myOut, out) {
				t.Errorf("\n Ref: %s \n Got: %s", out, myOut)
			}
		}

	})
}

// FIXME
/*func BenchmarkSearchWords(b *testing.B) {
	root, _ := trie.CreateTrie("../../words.txt")
	root.CompactTrie()
	trie.SaveTrie("dict.bin", root)
	//dict, _ := trie.LoadTrie("dict.bin")

	b.ResetTimer()
	//dict.SearchCloseWords("test", 0)
}*/

func readPrettyPrint(arr []trie.Word) []byte {
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	trie.PrettyPrint(arr)
	w.Close()
	myOut, _ := ioutil.ReadAll(r)
	os.Stdout = rescueStdout

	return myOut
}
