package main

import (
	"encoding/gob"
	"fmt"
	"os"
)

var words_map map[string][][2]int

type Shell struct {
	Words_map map[string][][2]int
}

// var sentences [][]string

func main() {
	words_mapout := map[string][][2]int{
		"a": {{0, 1}, {1, 1}, {2, 2}},
		"b": {{0, 2}, {1, 3}, {1, 4}},
	}
	shellout := Shell{Words_map: words_mapout}

	fmt.Println(words_mapout)

	SaveGob("words_map.gob", shellout)
	var shellin = new(Shell)
	LoadGob("words_map.gob", shellin)

	words_mapin := shellin.Words_map

	fmt.Println(words_mapin)
}

func SaveGob(file string, v interface{}) {
	saveFile, err := os.Create(file)
	Check(err, "save file cannot be created")
	defer saveFile.Close()
	encoder := gob.NewEncoder(saveFile)
	err = encoder.Encode(v)
	Check(err, "error in encoding to save file")
}

func LoadGob(file string, v interface{}) {
	loadFile, err := os.Open(file)
	Check(err, "cannot load saved file")
	defer loadFile.Close()
	decoder := gob.NewDecoder(loadFile)
	err = decoder.Decode(v)
	Check(err, "error in decoding from saved file")
}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}
