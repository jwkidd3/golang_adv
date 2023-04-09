package main

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
)

var pl = fmt.Println

func main() {
	// Create a file
	fileName := "data.txt"
	f, err := os.Create(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// Make sure closing the file after the program ends
	// defer f.Close()
	defer func() {
		f.Close()
		pl("Close file after creating")
	}()

	iPrimes := []int{2, 3, 5, 7, 11}
	var sPrimes []string
	for _, val := range iPrimes {
		sPrimes = append(sPrimes, strconv.Itoa(val))
	}

	// Write to files
	for _, num := range sPrimes {
		_, err := f.WriteString(num + "\n")
		if err != nil {
			log.Fatal(err)
		}
	}

	// Open file for read only (O_RDONLY)
	f, err = os.Open(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// defer f.Close()
	defer func() {
		f.Close()
		pl("Close file after reading")
	}()

	// Read the file contents
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		pl("Prime :", scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	// Append to file
	/*
		Exactly one of O_RDONLY, O_WRONLY, or O_RDWR must be specified
		O_RDONLY : open the file read-only
		O_WRONLY : open the file write-only
		O_RDWR   : open the file read-write

		These can be or'ed
		O_APPEND : append data to the file when writing
		O_CREATE : create a new file if none exists
		O_EXCL   : used with O_CREATE, file must not exist
		O_SYNC   : open for synchronous I/O
		O_TRUNC  : truncate regular writable file when opened
	*/
	_, err = os.Stat(fileName)
	if errors.Is(err, os.ErrNotExist) {
		pl("File doesn't exist")
	} else {
		f, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			log.Fatal(err)
		}

		// defer f.Close()
		defer func() {
			f.Close()
			pl("Close file after writing 13")
		}()

		if _, err := f.WriteString("13\n"); err != nil {
			log.Fatal(err)
		}
	}
}
