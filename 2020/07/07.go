package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang-collections/collections/set"
)

var VERBOSE = false

type Bag struct {
	pattern  string
	color    string
	children []*Bag
	parents  []*Bag
}

func (b *Bag) getType () string {
	return fmt.Sprintf("%s %s", b.pattern, b.color)
}

func (b *Bag) print () string {
	bagType := b.getType()
	childStr := ""
	for _, child := range b.children {
		childStr += child.getType()
		childStr += ", "
	}
	parentStr := ""
	for _, parent := range b.parents {
		parentStr += parent.getType()
		parentStr += ", "
	}

	return fmt.Sprintf("%s\n\tchildren: %s\n\tparents: %s",
		bagType, childStr, parentStr)
}

func printParents(b *Bag, i int, aset *set.Set) {
	for _, parent := range b.parents {
		aset.Insert(parent)
		if VERBOSE {
			fmt.Printf("\n%d|%s\n... is parent of %s\n", i, parent.print(), b.getType())
		}
		printParents(parent, i+1, aset)
	}
	i--
}


func printChildren(b *Bag, i int, childArrPtr *[]Bag) {

	for _, child := range b.children {
		*childArrPtr = append(*childArrPtr, *child)
		grandChildStr := ""
		for _, grandChild := range child.children {
			grandChildStr += grandChild.getType()
			grandChildStr += ", "
		}
		if VERBOSE {
			fmt.Printf("%d|%s, child of %s, has children: %s\n", i, child.getType(), b.getType(), grandChildStr)
		}
		printChildren(child, i+1, childArrPtr)
	}
	i--
}

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	dataMap := makeDataMap(bufio.NewScanner(f))

	fmt.Printf("*********** SHINY GOLD **************\n")
	shinyGold := dataMap["shiny gold"]

	aset := set.New()
	printParents(shinyGold, 0, aset)

	var childArr []Bag
	printChildren(shinyGold, 0, &childArr)

	fmt.Printf("set of parents has length %d\n", aset.Len())
	fmt.Printf("array of children has length %d\n", len(childArr))
}