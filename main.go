package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type reader struct {
	words  *[]record
	// rating *[]int -->  remove it
}

type record struct {
	word    []byte
	counter int
	// checked bool   -->  remove it
}

func (r *reader) contains(element []byte) (bool, int) {
	for index, v := range *r.words {
		if bytes.Equal(v.word, element) {
			return true, index
		}
		index = index + 1
	}
	return false, 0
}

func (r *reader) read_from_chan(ch chan []byte) {
	for node := range ch {
		state, index := r.contains(node)
		if state {
			(*r.words)[index].counter++
		} else {
			record := record{node, 1}
			*r.words = append(*r.words, record)
		}
	}
}


func (r *reader) get20mostfrequentwords() {
	sort.Slice(*r.words, func(i, j int) bool {
		return (*r.words)[i].counter > (*r.words)[j].counter
	})
	for i := 0; i < 20; i++ {
		fmt.Println((*r.words)[i].counter, " ", string((*r.words)[i].word))
	}
}

func main() {

	file, err := os.Open("mobydick.txt") //open file
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = file.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	readingBuf := make([]byte, 1) //read file by one letter only

	words := make([]record, 0)
	reader := reader{words: &words} //creating reader

	writingBuf := make([]byte, 0)
	// 1: writer := writer{&writingBuf}  -->  we remove in order to work with buf directly.

	ch := make(chan []byte)
	//channel that we will use to pass slices of bytes from writer to reader

	// 2:  use Reader instead of reading from file directly.
	r := bufio.NewReader(file)
	go func() {
		for {
			//reading file's letters one by one
			n, err := r.Read(readingBuf)

			if n > 0 {
				byteVal := readingBuf[0]
				if byteVal >= 65 && byteVal <= 90 { //if symbol is uppercase letter

					byteVal = byteVal + 32
					writingBuf = append(writingBuf, byteVal) //writing to temporary buffer

				} else if byteVal >= 97 && byteVal <= 122 { //if symbol is lowercase letter

					writingBuf = append(writingBuf, byteVal) //writing to temporary buffer

				} else if byteVal == 32 && len(writingBuf) != 0 { //if symbol is [space], and we have letters in our buffer

					ch <- writingBuf
					writingBuf = nil

				} else if ((byteVal > 122 || byteVal < 65) || (byteVal > 90 && byteVal < 97)) && len(writingBuf) != 0 {

					ch <- writingBuf
					writingBuf = nil

				} else {
					continue
				}
			}

			if err == io.EOF {
				ch <- writingBuf
				writingBuf = nil  //send temporary buffer content to channel, empty the temporary buffer
				break
			}
		}
		close(ch) //close channel, so our that our reader will stop working after there are no elements left, in other case reader will cause deadlock
	}()

	reader.read_from_chan(ch)

	reader.get20mostfrequentwords() //getting 20 most frequent words, and write it to rating slice
}
