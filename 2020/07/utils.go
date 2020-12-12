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

func addBagToDataMapIfNotExists(dataMapPtr *map[string]*Bag, patternColorArr []string, parentBag *Bag) *Bag {

	patternColor := strings.Join(patternColorArr, " ")
	nodeBag, ok := (*dataMapPtr)[patternColor]
	if VERBOSE {
		fmt.Printf("nodeBag %v, ok %t\n", nodeBag, ok)
	}
	if !ok { // if the nodeBag was not already in the dataMap, insert it now

		pattern := patternColorArr[0]
		color := patternColorArr[1]
		nodeBag = &Bag{
			pattern:  pattern,
			color:    color,
			children: []*Bag{},
		}
		if VERBOSE {
			fmt.Printf("adding '%s' to dataMap", nodeBag.getType())
		}
		if parentBag != nil {
			nodeBag.parents = []*Bag{parentBag}
			if VERBOSE {
				fmt.Printf(", child of %s", parentBag.getType())
			}
		}
		(*dataMapPtr)[nodeBag.getType()] = nodeBag
		if VERBOSE {
			fmt.Println()
		}
	} else {
		if parentBag != nil {
			nodeBag.parents = append(nodeBag.parents, parentBag)
		}
	}

	return nodeBag
}

func mungeData(dataStr string, dataMapPtr *map[string]*Bag) map[string]*Bag {
	nodeChildInfoArr := strings.Split(dataStr, " bags contain ")
	nodeType := nodeChildInfoArr[0]
	nodeTypeArr := strings.Split(nodeType, " ")
	if VERBOSE {
		fmt.Printf("about to check dataMap for nodeBag '%s'\n", strings.Join(nodeTypeArr, " "))
	}
	nodeBag := addBagToDataMapIfNotExists(dataMapPtr, nodeTypeArr, nil)

	childrenBagArr := make([]*Bag, 0)
	children := strings.Split(nodeChildInfoArr[1], ", ")

	// loop through all of nodeBag's children
	for _, child := range children {
		childInfoArr := strings.Split(child, " ")
		if childInfoArr[0] == "no" || childInfoArr[0] == "0" {
			continue
		}

		childNum, err := strconv.Atoi(childInfoArr[0])
		if err != nil {
			log.Fatalf("could not convert number of bags '%s' to integer\n", childInfoArr[0])
		}

		if VERBOSE {
			fmt.Printf("about to check dataMap for childBag '%s %s'\n", childInfoArr[1], childInfoArr[2])
		}
		childBag := addBagToDataMapIfNotExists(dataMapPtr, childInfoArr[1:3], nodeBag)

		for j:=0; j<childNum; j++ {
			childrenBagArr = append(childrenBagArr, childBag)
		}
	} // end of children loop

	nodeBag.children = childrenBagArr

	(*dataMapPtr)[fmt.Sprintf("%s %s", nodeBag.pattern, nodeBag.color)] = nodeBag
	return *dataMapPtr
}



func makeDataMap(s *bufio.Scanner) map[string]*Bag {
	dataMap := map[string]*Bag{}

	for s.Scan() {
		dataLineStr := s.Text()
		dataMap = mungeData(dataLineStr, &dataMap)
		//&dataMap = dataMapPtr
		if VERBOSE {
			fmt.Println("\n\nnew bagNode...")
			//fmt.Println(bagNode.print())
			//fmt.Println("\n")
		}

	}
	if err := s.Err(); err != nil {
		log.Fatal(err)
	}

	if VERBOSE {
		fmt.Printf("length of data map: %d\n", len(dataMap))
	}
	return dataMap
}

func getFile(fptr *string) *os.File {
	flag.Parse()

	f, err := os.Open(*fptr)
	if err != nil {
		log.Fatal(err)
	}
	return f
}