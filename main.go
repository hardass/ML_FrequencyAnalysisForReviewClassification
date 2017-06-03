package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
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
	Name           string
	Frequency      int
	RateOfFreq     float32
	Class_Average  float32
	Class_Variance float32
	Weight         float32
}

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
var verificationFileName string = "dev.txt"
var verificationResultFileName string = "dev_result.txt"

var debug bool = false

func main() {
	fmt.Println("Program Start...")

	// base data initialization
	InitializeBaseData()

	// analysis words requency
	WordsAnalysis()

	// result verification
	Verification()

	fmt.Println("Program End...")
}

func Verification() {
	verificationFile, err := os.Open(verificationFileName)
	Check(err, "verification file could not be opened")
	defer verificationFile.Close()

	notALetter := func(char rune) bool { return !unicode.IsLetter(char) }
	// var output []string
	type Output struct {
		Class   float64
		Content string
	}
	output := []Output{}
	var verificationResultClass []float64

	br := bufio.NewReader(verificationFile)
	for i := 0; ; i++ {
		line, _, flag := br.ReadLine()
		if flag == io.EOF {
			break
		}

		para := string(line)

		var class float32 = 0
		var match int = 0

		for _, word := range strings.FieldsFunc(para, notALetter) {
			word = strings.ToLower(word)

			if word_post_analysis, ok := words_post_analysis_map[word]; ok {
				class += word_post_analysis.Class_Average
				match++
			}
		}
		// output = append(output, strconv.FormatFloat(float64(class)/float64(match), 'f', 2, 32)+" "+para)
		verificationResultClass = append(verificationResultClass, float64(class)/float64(match))
		output = append(output, Output{Class: float64(class) / float64(match), Content: para})
	}

	// analysis training data's class distribution
	var trainingDataClassSatatic [5]int
	for i, _ := range sentences {
		j, _ := strconv.Atoi(sentences[i][0])
		switch j {
		case 0:
			trainingDataClassSatatic[0]++
		case 1:
			trainingDataClassSatatic[1]++
		case 2:
			trainingDataClassSatatic[2]++
		case 3:
			trainingDataClassSatatic[3]++
		case 4:
			trainingDataClassSatatic[4]++
		}
	}
	fmt.Println(trainingDataClassSatatic)
	// e.g. [1082 2228 1621 2328 1285]
	var trainingDataClassPercentage [5]float64
	for i, _ := range trainingDataClassSatatic {
		trainingDataClassPercentage[i] = float64(trainingDataClassSatatic[i]) / float64(len(sentences))
	}
	fmt.Println(trainingDataClassPercentage)
	// e.g. [0.1266385767790262 0.2607677902621723 0.1897237827715356 0.27247191011235955 0.15039794007490637]
	// calculate class thresholds
	sort.Float64s(verificationResultClass)
	// fmt.Println(verificationResultClass)
	var classThreshods [5]float64
	var p float64
	for i, _ := range trainingDataClassPercentage {
		p += trainingDataClassPercentage[i]
		classThreshods[i] = verificationResultClass[int(float64(len(verificationResultClass))*p)]
	}
	fmt.Println(classThreshods)
	// e.g. [1.851025644938151 2.0055873870849608 2.0936654743395353 2.272305443173363 2.8736173629760744]
	for i, _ := range output {
		switch {
		case output[i].Class < classThreshods[0]:
			output[i].Content = "0 " + output[i].Content
		case output[i].Class < classThreshods[1]:
			output[i].Content = "1 " + output[i].Content
		case output[i].Class < classThreshods[2]:
			output[i].Content = "2 " + output[i].Content
		case output[i].Class < classThreshods[3]:
			output[i].Content = "3 " + output[i].Content
		default:
			output[i].Content = "4 " + output[i].Content

		}
	}
	for _, s := range output {
		fmt.Println(s.Content)
	}
}

func WordsAnalysis() {
	var totalWordsCount int = 0
	for _, value := range words_map {
		totalWordsCount = totalWordsCount + len(value)

	}
	fmt.Println("total words:", totalWordsCount)

	var overallAverageClass float32 = 0
	for i, _ := range sentences {
		class, _ := strconv.Atoi(sentences[i][0])
		overallAverageClass = overallAverageClass + float32(class)
	}
	overallAverageClass = overallAverageClass / float32(len(sentences))
	fmt.Println("overall average class:", overallAverageClass)

	var words Words
	words_post_analysis_map = make(map[string]Word)

	// update how many times a word appears in whole training data set, and its percentage
	for key, value := range words_map {
		words = append(words, Word{Name: key, Frequency: len(value), RateOfFreq: float32(len(value)) / float32(totalWordsCount)})
	}

	// update the average of paragraph's classes containing the word
	for i, word := range words {
		// get [][2]int of each word
		v, _ := words_map[word.Name]

		var sum int = 0
		for _, p := range v {
			class, _ := strconv.Atoi(sentences[p[0]][0])
			sum = sum + class
		}
		words[i].Class_Average = float32(sum) / float32(len(v))
		// update variance of each word
		words[i].Class_Variance = float32(math.Pow(float64(words[i].Class_Average-overallAverageClass), 2))

		words[i].Weight = words[i].Class_Variance * words[i].RateOfFreq

		// write to post analysis word map
		words_post_analysis_map[words[i].Name] = words[i]
	}

	sort.Sort(sort.Reverse(words))
	// for i, word := range words {
	// 	fmt.Printf("#.%d: %s , %d , %.9f , %.9f , %.9f , %.9f\n", i, word.Name, word.Frequency, word.RateOfFreq, word.Class_Average, word.Class_Variance, word.Weight)
	// }
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
		fmt.Println("pre-loaded file doesn't exist")
	} else {
		// load train data structure from file
		gobShellload := new(GobShell)
		LoadTrainDataStruct(gobFileName, gobShellload)
		words_map = gobShellload.Words_map
		sentences = gobShellload.Sentences
		fmt.Println("pre-loaded file exists")
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
				//convert word to lower case
				word = strings.ToLower(word)
				words_map[word] = append(words_map[word], [][2]int{{i, j + 1}}...)
			}
		}
		sentences = append(sentences, currentSentence)
	}

	// for key, value := range words_map {
	// 	fmt.Printf("%s : %d\n", key, value)
	// }
}
