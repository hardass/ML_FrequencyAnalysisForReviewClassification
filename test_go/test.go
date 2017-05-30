package main

import (
	// "encoding/gob"
	"fmt"
	// "os"
)

var words_map map[string][][2]int

// var sentences [][]string

func main() {
	words_map := map[string][][2]int{
		"a": {{0, 1}, {1, 1}, {2, 2}, {3, 3}},
		"b": {{0, 2}, {1, 3}, {1, 4}},
	}

	for key, value := range words_map {
		fmt.Println(key, len(value))
	}

}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}
