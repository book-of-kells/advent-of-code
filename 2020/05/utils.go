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

// converts a string like 'FBBFBFFLLR' into two integers
//	example:
//- `BFFFBBFRRR`: row 70, column 7, seat ID 567.
//- `FFFBBBFRRR`: row 14, column 7, seat ID 119.
//- `BBFFBBFRLL`: row 102, column 4, seat ID 820.
func mungeData(dataStr string) (*Seat, error) {
	// do stuff here like
	rowStr := dataStr[:7]
	binaryRowStr := strings.ReplaceAll(
		strings.ReplaceAll(rowStr, "F", "0"), "B", "1")
	colStr := dataStr[7:]
	binaryColStr := strings.ReplaceAll(
		strings.ReplaceAll(colStr, "L", "0"), "R", "1")

	s := Seat{}
	if decRowNum, err := strconv.ParseInt(binaryRowStr, 2, 8); err != nil {
		return nil, err
	} else {
		s.row = int(decRowNum)
	}
	if decColNum, err := strconv.ParseInt(binaryColStr, 2, 8); err != nil {
		return nil, err
	} else {
		s.col = int(decColNum)
	}

	id := 8*s.row + s.col
	s.id = &id
	return &s, nil
}

func makeDataArr(s *bufio.Scanner) []*Seat {
	dataArr := make([]*Seat, 1024)
	for s.Scan() {
		dataLineStr := s.Text()
		if seat, err := mungeData(dataLineStr); err == nil {
			dataArr[seat.getId()] = seat
		} else {
			fmt.Println(err)
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