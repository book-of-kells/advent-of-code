package main

import (
	"bufio"
	"flag"
	"fmt"
	"sort"
)

var VERBOSE = false


func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))
	dataArr = append(dataArr, 0)
	sort.Ints(dataArr)

	dataArr = append(dataArr, dataArr[len(dataArr)-1]+3)

	fmt.Printf("sorted: %v\n", dataArr)

	/*
	testinput2.txt
	one-drop: 22
	three-drop: 10
	sorted: [0 1 2 3 4 7 8 9 10 11 14 17 18 19 20 23 24 25 28 31 32 33 34 35 38 39 42 45 46 47 48 49 52]
	map[1:[0 1 2 3 5 6 7 8 11 12 13 15 16 19 20 21 22 24 27 28 29 30] 2:[] 3:[4 9 10 14 17 18 23 25 26 31]]
	19208 combinations
	*/

	joltDrops := getJoltDrops(dataArr)
	fmt.Println(joltDrops)
	answer := len(joltDrops[1])*len(joltDrops[3])
	fmt.Printf("answer: %d\n", answer)
}

func getJoltDrops(dataArr []int) map[int][]int {
	joltDropMap := map[int][]int{
		1: make([]int, 0),
		2: make([]int, 0),
		3: make([]int, 0),
	}
	for i, joltage := range dataArr {
		if i == len(dataArr) - 1 {
			return joltDropMap
		}
		key := dataArr[i+1] - joltage
		joltDropMap[key] = append(joltDropMap[key], i)
	}
	return joltDropMap
}

/*
 * Takes a joltDropMap and creates data structures called "substitutions" for now, for lack of a better word.
 * The plan was to figure out how many combinations would be possible by finding the number of possible
 * "substitution" combinations.
 * I got this far on my own, then went online and read that people were doing the same thing and then solving
 * it by using the tribonacci sequence on these "substitutions".
 */
func getSubstitutions(joltDropMap map[int][]int) {

	// 1. get list of one-drop > two-drop substitutions
	// todo data structures with substitutions only
	// 			> one-drop: [2, 3] 	> []
	//			> two-drop: []		> [2]

	// 2. get list of one-drop > three-drop substitutions
	// todo data structures with substitutions only
	// 			> one-drop: [2, 3, 4]	> []
	//			> three-drop: []		> [2]

}

func getOneDropToTwoDropSubstitutions(joltDropMap map[int][]int) []map[int][]int {
	/*
	 * where there are two one-jolt drops consecutively,
	 * drop them both and the lowest index goes to the two-drop array
	 */

	possibleJoltDropMapsArr := make([]map[int][]int, 0)
	oneJoltDropArr := joltDropMap[1]
	twoJoltDropArr := joltDropMap[2]
	for i := 1; i<len(oneJoltDropArr)-1; i++ {

		currOneDropIdx := oneJoltDropArr[i]
		nextOneDropIdx := oneJoltDropArr[i+1]
		if currOneDropIdx == nextOneDropIdx + 1{
			newJoltDropmap := map[int][]int{
				1: append(oneJoltDropArr[:i], oneJoltDropArr[i+2:]...),
				2: append(twoJoltDropArr, currOneDropIdx),
				3: joltDropMap[3],
			}
			possibleJoltDropMapsArr = append(possibleJoltDropMapsArr, newJoltDropmap)
		}

	}
	// todo alternatively, return a new data structure with substitutions only
	// 			> one-drop: [2, 3] 	> []
	//			> two-drop: []		> [2]
	return possibleJoltDropMapsArr
}

func getOneDropToThreeDropSubstitutions(joltDropMap map[int][]int) {
	/*
	 * where there are three one-jolt drops consecutively,
	 * drop them both and the lowest index goes to the three-drop array
	 */
}

func getTwoDropToThreeDropSubstitutions(joltDropMap map[int][]int) {
	/*
	 *
	 */
}