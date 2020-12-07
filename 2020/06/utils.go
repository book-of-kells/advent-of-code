package main

import (
	"bufio"
	"flag"
	"github.com/golang-collections/collections/set"
	"log"
	"os"
)

func mungeData(groupDataArr []string) set.Set {

	// initialize groupCharSet with the first person in the group.
	groupCharSet := set.New()
	for _, char := range groupDataArr[0] {
		groupCharSet.Insert(string(char))
	}

	for _, yesAnswerStr := range groupDataArr {
		personCharSet := set.New()

		// get all letters in yesAnswerStr
		for _, char := range yesAnswerStr {
			personCharSet.Insert(string(char))
		}

		groupCharSet = groupCharSet.Intersection(personCharSet)
	}

	return *groupCharSet
}

func makeDataArr(s *bufio.Scanner) []set.Set {
	dataArr := make([]set.Set, 0)
	groupDataArr := make([]string, 0)
	for s.Scan() {
		if dataLineStr := s.Text(); dataLineStr != "" {
			groupDataArr = append(groupDataArr, dataLineStr)
		} else {
			// blank line signals end of data
			dataArr = append(dataArr, mungeData(groupDataArr))
			groupDataArr = make([]string, 0)
		}
	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}
	return dataArr
}

func getFile() *os.File {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}