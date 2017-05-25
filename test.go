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

	br := bufio.NewReader(practiceFile)
	for i := 0; ; i++ {
		line, _, flag := br.ReadLine()
		if flag == io.EOF {
			break
		}
		sen := string(line)
		fmt.Println("Line ", i, ": ", sen)
		notALetter := func(char rune) bool { return !unicode.IsLetter(char) }
		for j, word := range strings.FieldsFunc(sen, notALetter) {
			fmt.Println("word #", j, ": ", word)
		}
	}

}

func test() {
	fmt.Println("test start")
	////append 2 slice
	words_map = map[string][][2]int{
		"a": {{0, 1}, {1, 1}, {2, 2}},
		"b": {{0, 2}, {1, 3}, {1, 4}},
	}
	words_map["a"] = append(words_map["a"], [][2]int{{3, 3}}...)

	fmt.Println(words_map["a"])
	for key, value := range words_map {
		fmt.Printf("%s : %d\n", key, value)
	}
	fmt.Println("test end")
	////

	////read one line from txt file
	practiceFile, err := os.Open("train.txt")
	Check(err, "practice file cannot be loaded")
	defer practiceFile.Close()

	scanner := bufio.NewScanner(practiceFile)
	scanner.Split(bufio.ScanLines)

	success := scanner.Scan()
	if success == false {
		err = scanner.Err()
		if err == nil {
			fmt.Println("scan to the end of the file")
		} else {
			panic(err)
		}
	}
	fmt.Println("frist line is: ", scanner.Text())
	////
}
