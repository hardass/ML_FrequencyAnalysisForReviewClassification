package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode"
)
var words_map map[string][][2]int
var sentences [][]string

func main() {
	words_map = map[string][][2]int{
		"a": {{0, 1}, {1, 1}, {2, 2}},
		"b": {{0, 2}, {1, 3}, {1, 4}},
	}

	fmt.Println(words_map)
}

func Check(err error, log string) {
	if err != nil {
		fmt.Println(log)
		panic(err)
	}
}
