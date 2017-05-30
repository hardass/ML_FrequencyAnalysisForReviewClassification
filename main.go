package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"sort"
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

type Word struct {
	Name       string
	Frequency  int
	RateOfFreq float32
}

type SorterWordFrequency []Word

func (sort SorterWordFrequency) Len() int {
	return int(len(sort))
}
func (sorter SorterWordFrequency) Less(i, j int) bool {
	return sorter[i].Frequency < sorter[j].Frequency
}
func (sorter SorterWordFrequency) Swap(i, j int) {
	sorter[i], sorter[j] = sorter[j], sorter[i]
}

// var words_analysis_map map[string]Word

type GobShell struct {
	Words_map map[string][][2]int
	Sentences [][]string
}

var gobFileName string = "trainingdata.gob"

var debug bool = false

func main() {
	fmt.Println("Program Start...")

	// base data initialization
	InitializeBaseData()

	// analysis words requency
	WordsAnalysis()

	fmt.Println("Program End...")
}

func WordsAnalysis() {
	var totalWordsCount int = 0
	for _, value := range words_map {
		totalWordsCount = totalWordsCount + len(value)

	}
	fmt.Println(totalWordsCount)

	// words_analysis_map = make(map[string]Word)
	var words_sort_frequency SorterWordFrequency

	for key, value := range words_map {
		// words_analysis_map[key] = Word{Frequency: int32(len(value)), RateOfFreq: float32(len(value)) / float32(totalWordsCount)}
		words_sort_frequency = append(words_sort_frequency, Word{Name: key, Frequency: len(value), RateOfFreq: float32(len(value)) / float32(totalWordsCount)})
	}
	sort.Sort(sort.Reverse(words_sort_frequency))
	// sort.Reverse(words_sort_frequency)

	// for key, value := range words_analysis_map {
	// 	fmt.Printf("%s : %d : %.19f \n", key, value.Frequency, value.RateOfFreq)
	// }

	for i, word := range words_sort_frequency {
		fmt.Printf("No. %d: %s : %d : %.19f \n", i, word.Name, word.Frequency, word.RateOfFreq)
	}
}

func InitializeBaseData() {
	// if train data stuct file exists, load data directly; else read original train txt file and save to data struct file
	_, err := os.Stat(gobFileName)
	if err != nil && os.IsNotExist(err) {
		// read practice data from train.txt
		ReadPracticeData()

		// persistence of train data as structure
		gobShellSave := GobShell{Words_map: words_map, Sentences: sentences}
		SaveTrainingDataStruct(gobFileName, gobShellSave)
		fmt.Println("doesn't exist")
	} else {
		// load train data structure from file
		gobShellload := new(GobShell)
		LoadTrainDataStruct(gobFileName, gobShellload)
		words_map = gobShellload.Words_map
		sentences = gobShellload.Sentences
		fmt.Println("exists")
	}

	// fmt.Println(words_mapload)
	if debug {
		fmt.Println(sentences)

		for key, value := range words_map {
			fmt.Printf("%s : %d\n", key, value)
		}
	}
}

func LoadTrainDataStruct(file string, v interface{}) {
	loadFile, err := os.Open(file)
	Check(err, "cannot load saved file")
	defer loadFile.Close()
	decoder := gob.NewDecoder(loadFile)
	err = decoder.Decode(v)
	Check(err, "error in decoding from saved file")
}

func SaveTrainingDataStruct(file string, v interface{}) {
	saveFile, err := os.Create(file)
	Check(err, "save file cannot be created")
	defer saveFile.Close()
	encoder := gob.NewEncoder(saveFile)
	err = encoder.Encode(v)
	Check(err, "error in encoding to save file")
}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}

func ReadPracticeData() {
	// practiceFile, err := ioutil.ReadFile("train.txt")
	practiceFile, err := os.Open("train.full.txt")
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
		// fmt.Println("Para ", i, ": ", para)

		//read the paragragh's class
		class := strings.IndexFunc(para, unicode.IsSpace)
		// fmt.Println("class: ", para[:class])
		currentSentence := []string{para[:class]}

		for j, sen := range strings.Split(para, ",") {
			//read sentence
			// fmt.Println("Sen ", j+1, ": ", sen)
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

	// for key, value := range words_map {
	// 	fmt.Printf("%s : %d\n", key, value)
	// }
}
