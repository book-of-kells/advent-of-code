package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func mungeData(dataStr string) (*Passport, error) {
	// 1. split by space delimiter, then trim()
	passportDataArr := strings.Split(strings.Trim(dataStr, " "), " ")

	// 2. loop through entries and add key, val to map
	passportDataMap := make(map[string]string, 0)
	for _, value := range passportDataArr {
		keyFieldArr := strings.Split(strings.Trim(value, " "), ":")
		key := strings.Trim(keyFieldArr[0], " ")
		passportDataMap[strings.Title(key)] = strings.Trim(keyFieldArr[1], " ")
	}
	return NewPassport(passportDataMap)
}

func makeDataArr(s *bufio.Scanner) []Passport {
	dataArr := make([]Passport, 0)
	totalPassports := 0
	validPassports := 0
	passportDataStr := ""
	for s.Scan() {
		if lineStr := s.Text(); lineStr != "" {
			passportDataStr += " "
			passportDataStr += lineStr
		} else {
			// blank line signals end of passport data
			if p, err := mungeData(passportDataStr); p != nil && err == nil {
				dataArr = append(dataArr, *p)
				validPassports++
			} else {
				fmt.Println(err)
			}
			passportDataStr = ""
			totalPassports++
		}
	}
	if passportDataStr != "" {
		if p, err := mungeData(passportDataStr); p != nil && err == nil{
			dataArr = append(dataArr, *p)
			validPassports++
		} else {
			fmt.Println(err)
		}
		passportDataStr = ""
		totalPassports++
	}

	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("(%d or %d)/%d valid passports\n", validPassports, len(dataArr), totalPassports)

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