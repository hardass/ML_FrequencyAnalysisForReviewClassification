package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)

// type Word struct {
// 	key string
// 	senIdx string
// 	paraIdx string
// 	cls string
// }

// Map "word" as key, array of
// map[string][p][s]int, [p] is paragraph, 0->n; [s] is sentence, 1->n, and 0 is [p]'s class
var words_map map[string][][2]int

// {
// 	"good": {{0, 1}, {1, 1}, {2, 2}}, <- word of "good" appears in #0 paragraph's #1 sentence, and #1p's #1s, and #2p's #2s;
// 	"bad": {{0, 2}, {1, 3}, {1, 4}}, <- word of "bad" appears in #0p's #2s, and #1p's #3s, and #1p's #4s
// }

var sentences [][]string

// {
// 	{"3","s1","s2"},
// 	{"1","s1","s2","s3"},
// 	...
// }
// Paragraph 0's class is 3, has sentences 1 and 2;
// Paragraph 1's class is 1, has sentences 1, 2 and 3;

func main() {
	fmt.Println("Program Start...")

	// test
	// test()

	// read practice data from train.txt
	ReadPracticeData()

	fmt.Println("Program End...")
}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}

func ReadPracticeData() {
	// practiceFile, err := ioutil.ReadFile("train.txt")
	practiceFile, err := os.Open("train.txt")
	Check(err, "practice file cannot be loaded")
	defer practiceFile.Close()

	notALetter := func(char rune) bool { return !unicode.IsLetter(char) }
	words_map = make(map[string][][2]int)

	br := bufio.NewReader(practiceFile)
	for i := 0; ; i++ {
		line, _, flag := br.ReadLine()
		if flag == io.EOF {
			break
		}

		//read a paragragh
		para := string(line)
		fmt.Println("Para ", i, ": ", para)

		//read the paragragh's class
		class := strings.IndexFunc(para, unicode.IsSpace)
		fmt.Println("class: ", para[:class])
		currentSentence := []string{para[:class]}

		for j, sen := range strings.Split(para, ",") {
			//read sentence
			fmt.Println("Sen ", j+1, ": ", sen)
			//finish sentences construction
			// sentences[i][j+1] = sen
			currentSentence = append(currentSentence, sen)

			for _, word := range strings.FieldsFunc(sen, notALetter) {
				//read word
				//finish word_map construction
				words_map[word] = append(words_map[word], [][2]int{{i, j + 1}}...)
			}
		}
		sentences = append(sentences, currentSentence)
	}

	for key, value := range words_map {
		fmt.Printf("%s : %d\n", key, value)
	}

}