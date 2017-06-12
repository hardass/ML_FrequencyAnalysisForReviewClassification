package main

import (
	"bufio"
	// "encoding/gob"
	"fmt"
	"io"
	// "math"
	"os"
	// "sort"
	"strconv"
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
// var words_map map[string][][2]int
var words_map map[string]Word

type Word struct {
	// Name       string
	Appearance [][2]int
	// {
	// 	"good": {{0, 1}, {1, 1}, {2, 2}}, <- word of "good" appears in #0 paragraph's #1 sentence, and #1p's #1s, and #2p's #2s;
	// 	"bad": {{0, 2}, {1, 3}, {1, 4}}, <- word of "bad" appears in #0p's #2s, and #1p's #3s, and #1p's #4s
	// }
	WordRankOnEmergeTimes [5]float64
	WordRankOnEmergeRatio [5]float64
	// Frequency      int
	// RateOfFreq     float64
	// Class_Average  float64
	// Class_Variance float64
	// Weight         float64
}

var trainingParagraph_array []Paragraph
var testingParagraph_array []Paragraph

type Paragraph struct {
	Rank           int
	Sentences      []string
	RankPrediction int
}

var trainingSentences [][]string

// {
// 	{"3","s1","s2"},
// 	{"1","s1","s2","s3"},
// 	...
// }
// Paragraph 0's class is 3, has sentences 1 and 2;
// Paragraph 1's class is 1, has sentences 1, 2 and 3;

var testingSentences [][]string

/*
type Words []Word

func (sort Words) Len() int {
	return int(len(sort))
}
func (sorter Words) Less(i, j int) bool {
	// return sorter[i].Frequency < sorter[j].Frequency
	// return sorter[i].Class_Variance < sorter[j].Class_Variance
	return sorter[i].Weight < sorter[j].Weight
}
func (sorter Words) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}


var words_post_analysis_map map[string]Word

type GobShell struct {
	Words_map map[string][][2]int
	Sentences [][]string
}

var gobFileName string = "trainingdata.gob"
*/

var testingFileName string = "dev.txt"
var verificationResultFileName string = "dev_result.txt"

type testingFileOutput struct {
	EstimatedWeight float64
	CalculatedClass float64
	ActualClass     float64
	Content         string
}

var output = []testingFileOutput{}

var debug bool = false

func main() {
	fmt.Println("Program Start...")

	// train data process
	Training()

	// result verification
	Verification()

	Score()

	fmt.Println("Program End...")
}

func Score() {
	j := 0
	for _, s := range testingParagraph_array {
		if s.Rank == s.RankPrediction {
			j++
		}
	}
	fmt.Println("match:", float64(j)/float64(len(testingParagraph_array)))
}

func Verification() {
	testingFile, err := os.Open(testingFileName)
	Check(err, "testing file could not be opened")
	defer testingFile.Close()

	notALetter := func(char rune) bool { return !unicode.IsLetter(char) }

	br := bufio.NewReader(testingFile)
	for i := 0; ; i++ {
		line, _, flag := br.ReadLine()
		if flag == io.EOF {
			break
		}

		para := string(line)

		t, _ := strconv.Atoi(para[:strings.IndexFunc(para, unicode.IsSpace)])
		paragraph := Paragraph{Rank: int(t), Sentences: []string{para}}

		var ratioOfParagraph [5]float64
		// sum each word's ratio as paragraph's ratio
		for _, word := range strings.FieldsFunc(para, notALetter) {
			word = strings.ToLower(word)

			if wordStatisic, ok := words_map[word]; ok {
				for i := 0; i < len(ratioOfParagraph); i++ {
					ratioOfParagraph[i] += wordStatisic.WordRankOnEmergeRatio[i]
				}
			}
		}

		// found the max rank posibility
		rankPredict := 0
		ratioPredict := 0.0
		for i, v := range ratioOfParagraph {
			if v > ratioPredict {
				ratioPredict = v
				rankPredict = i
			}
		}
		paragraph.RankPrediction = rankPredict
		testingParagraph_array = append(testingParagraph_array, paragraph)
	}

	// for _, v := range testingParagraph_array {
	// 	fmt.Println(v.Rank, v.RankPrediction)
	// }
}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}

func Training() {
	// practiceFile, err := ioutil.ReadFile("train.txt")
	practiceFile, err := os.Open("train.full.txt")
	Check(err, "practice file cannot be loaded")
	defer practiceFile.Close()

	notALetter := func(char rune) bool { return !unicode.IsLetter(char) }
	words_map = make(map[string]Word)

	br := bufio.NewReader(practiceFile)
	for i := 0; ; i++ {
		line, _, flag := br.ReadLine()
		if flag == io.EOF {
			break
		}

		//read a paragragh
		para := string(line)
		// fmt.Println("Para ", i, ": ", para)
		class_position := strings.IndexFunc(para, unicode.IsSpace)

		//read the paragragh's class
		class, _ := strconv.Atoi(para[:class_position])
		paragraph := Paragraph{Rank: class}
		// fmt.Println("class: ", para[:class])
		// e.g. {0 [0 Her fans walked out muttering words like `` horrible '' and `` terrible   '' but had so much fun dissing the film that they did n't mind the ticket cost .]}
		// e.g. {1 [1 In this case zero .]}
		var sentence []string

		for j, sen := range strings.Split(para, ",") {
			//read sentence
			// e.g. Sen  1 :  3 The Rock is destined to be the 21st Century 's new `` Conan '' and that he 's going to make a splash even greater than Arnold Schwarzenegger
			// e.g. Sen  2 :   Jean-Claud Van Damme or Steven Segal .
			//finish sentences construction
			sentence = append(sentence, sen)

			for _, wordname := range strings.FieldsFunc(sen, notALetter) {
				//read word
				//finish word_map construction
				//convert word to lower case
				wordname = strings.ToLower(wordname)
				word := words_map[wordname]

				word.Appearance = append(word.Appearance, [][2]int{{i, j}}...)

				word.WordRankOnEmergeTimes[class]++

				words_map[wordname] = word
			}
		}
		paragraph.Sentences = sentence
		trainingParagraph_array = append(trainingParagraph_array, paragraph)
	}

	// update WordRankOnEmergeRatio
	for name, word := range words_map {
		times := 0.0
		for _, v := range word.WordRankOnEmergeTimes {
			times += v
		}
		for i, _ := range word.WordRankOnEmergeRatio {
			word.WordRankOnEmergeRatio[i] = word.WordRankOnEmergeTimes[i] / times
		}
		words_map[name] = word
	}

	// z := 0
	// for key, value := range words_map {
	// 	fmt.Printf("%s : %.9f : %.9f \n", key, value.WordRankOnEmergeTimes, value.WordRankOnEmergeRatio)
	// 	if z > 100 {
	// 		break
	// 	}
	// 	z++
	// }

	// fmt.Println(trainingParagraph_array)
}
