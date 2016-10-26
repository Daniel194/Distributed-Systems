package main

import (
	"io"
	"os"
	"bytes"
)

func main() {
	// open input file
	fi, err := os.Open("/Users/daniellungu/Documents/Workspace/Distributed-Systems/resources/input.txt")
	if err != nil {
		panic(err)
	}
	// close fi on exit and check for its returned error
	defer func() {
		if err := fi.Close(); err != nil {
			panic(err)
		}
	}()

	// open output file
	fo, err := os.Create("/Users/daniellungu/Documents/Workspace/Distributed-Systems/resources/output.txt")
	if err != nil {
		panic(err)
	}
	// close fo on exit and check for its returned error
	defer func() {
		if err := fo.Close(); err != nil {
			panic(err)
		}
	}()

	// make a buffer to keep chunks that are read
	buf := make([]byte, 1024)
	for {
		// read a chunk
		n, err := fi.Read(buf)
		if err != nil && err != io.EOF {
			panic(err)
		}
		if n == 0 {
			break
		}

		oldString := string(buf)
		newString := takeOutVowels(oldString)

		newBuffer := []byte(newString)

		// write a chunk
		if _, err := fo.Write(newBuffer[:n]); err != nil {
			panic(err)
		}
	}
}

func takeOutVowels(intput string) string {
	var output bytes.Buffer

	for i := 0; i < len(intput); i++ {
		value := intput[i]

		if value != 97 && value != 101 && value != 105 && value != 111 && value != 117 {
			output.WriteString(string(value))
		}
	}

	return output.String()
}
