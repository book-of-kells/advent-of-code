package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type pwdData struct {
	password string
	ruleLetter string
	firstInt int
	secondInt int
}

func getPwdData(pwdLine string) pwdData {
	pwdArr := strings.Split(pwdLine, ":")
	pwd := strings.Trim(pwdArr[1], " ")
	ruleArr := strings.Split(pwdArr[0], " ")
	ruleLetter := strings.Trim(ruleArr[1], " ")

	// parse lower and upper bounds of rule
	ruleStr := strings.Trim(ruleArr[0], " ")
	boundsArr := strings.Split(ruleStr, "-")
	firstInt, err := strconv.Atoi(boundsArr[0])
	if err != nil{
		log.Fatalf("error converting %s to integer: %v\n", boundsArr[0], err)
	}
	secondInt, err := strconv.Atoi(boundsArr[1])
	if err != nil {
		log.Fatalf("error converting %s to integer: %v\n", boundsArr[1], err)
	}


	return pwdData{pwd, ruleLetter, firstInt, secondInt}
}

func isPasswordValidPartOne(pwdLine string) bool {
	pwdData := getPwdData(pwdLine)

	// get count of ruleLetter occurrences in pwd
	numOccurrences := strings.Count(pwdData.password, pwdData.ruleLetter)

	lowerInt := pwdData.firstInt
	upperInt := pwdData.secondInt

	if lowerInt > upperInt {
		tmp := upperInt
		upperInt = lowerInt
		lowerInt = tmp
	}

	// determine if number of occurrences of the letter is within the bounds
	return numOccurrences >= lowerInt && numOccurrences <= upperInt
}

func isFirstMatch(pwdData pwdData) bool {
	return string(pwdData.password[pwdData.firstInt-1]) == pwdData.ruleLetter
}

func isSecondMatch(pwdData pwdData) bool {
	return string(pwdData.password[pwdData.secondInt-1]) == pwdData.ruleLetter
}

func isOnlyOneTrue(first bool, second bool) bool {
	return (first && !second) || (!first && second)
}

func isPasswordValidPartTwo(pwdLine string) bool {
	pwdData := getPwdData(pwdLine)
	return isOnlyOneTrue(isFirstMatch(pwdData), isSecondMatch(pwdData))
}


func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err = f.Close(); err != nil {
			log.Fatal(err)
		}
	}()
	s := bufio.NewScanner(f)

	currIdx := 0
	numValid := 0
	for s.Scan() {
		pwdLine := s.Text()
		if isPasswordValidPartTwo(pwdLine) {
			numValid++
			fmt.Printf("%s is valid!\n", pwdLine)
		}
		currIdx++
	}
	// 329 is too low
	// changed rules then I got 314
	fmt.Printf("there are %d passwords total and %d are valid!\n", currIdx, numValid)
}