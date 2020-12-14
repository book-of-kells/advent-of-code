package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/golang-collections/collections/set"
	"sort"
)

var VERBOSE = false

type Substitution struct {
	first int
	second int
	third *int
}

func (s *Substitution) String() string {
	if s.third != nil {
		third := s.third
		return fmt.Sprintf("{%d %d %d}\n", s.first, s.second, *third)
	}
	return fmt.Sprintf("{%d %d}\n", s.first, s.second)
}

func (s *Substitution) getLastInt() int {
	third := s.third
	if third != nil {
		return *third
	} else {
		return s.second
	}
}

func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()

	// make data array with a 0-jolt at the beginning and a +3 jolt at the end
	dataArr := makeDataArr(bufio.NewScanner(f))
	dataArr = append(dataArr, 0)
	sort.Ints(dataArr)
	dataArr = append(dataArr, dataArr[len(dataArr)-1]+3)
	fmt.Printf("sorted: %v\n", dataArr)

	// get a map of the indexes of one-, two-, and three-jolt drops
	// todo clarify this comment
	joltDrops := getJoltDrops(dataArr)
	fmt.Println(joltDrops)

	/*answer := len(joltDrops[1])*len(joltDrops[3])
	fmt.Printf("answer for part one: %d\n", answer)*/

	// get an an array of Substitution objects, which each represent one of all the possible
	// ways a one-jolt drop could be changed into a two- or three-volt drop
	subArr := getSubstitutions(joltDrops[1])
	fmt.Println("\nsubArr")
	for _, sub := range subArr {
		fmt.Printf("%s ", sub.String())
	}

	// cluster the substitutions together where they are overlapping
	clusterArr := clusterSubstitutions(subArr)

	// clusters of a single substitution have two possible combinations -- they are 1. substituted or 2. not
	// clusters of three overlapping substitutions have four possible combinations
	// clusters of five overlapping substitutions have seven possible combinations
	clusterLenToComboNumMap := map[int]int {
		1: 2,
		3: 4,
		5: 7,
	}

	// calculate the # of possible jolt adapter combinations by multiplying each substitution cluster's # of possible combinations
	totalNumCombinations := 1
	fmt.Println("\nclusterArr")
	for _, cluster := range clusterArr {
		fmt.Println(cluster)
		clusterLen := cluster.Len()
		numCombosForCluster := clusterLenToComboNumMap[clusterLen]
		totalNumCombinations *= numCombosForCluster
	}

	fmt.Printf("totalNumCombinations: %d\n", totalNumCombinations)
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
 * I don't know yet how to solve this problem with the tribonacci algorithm, so I'm going with this
 * permutation algorithm for now.
 *
 * @param oneDropArr: the array of indexes where a one-drop jolt adapter is used
 * @return
 */
func getSubstitutions(oneDropArr []int) []Substitution{
	/* Group together consecutive numbers in arrays of 2 or 3 In this example, 0-3 and 5-6 are consecutive
	 * [0 1 2 3 5 6....]
	 * so there are several substitutions from 1-jolt to 2- or 3-jolt adapters,
	 * represented by arrays of 2 to 3 consecutive integers
	 *	[
	 *		[0, 1],
	 *		[0, 1, 2],
	 *		[1, 2],
	 *		[1, 2, 3],
	 *		[2, 3],
	 * 		[5, 6],
	 *		...
	 *	]
	 */

	var subArr []Substitution
	oneDropArrLen := len(oneDropArr)
	for i, x := range oneDropArr {
		if i == oneDropArrLen - 1 {
			continue
		}
		y := oneDropArr[i+1]
		if sub := reduceTwo(x, y); sub != nil {
			subArr = append(subArr, *sub)
		}
		if i == oneDropArrLen - 2 {
			continue
		}
		z := oneDropArr[i+2]
		if sub := reduceThree(x, y, z); sub != nil {
			subArr = append(subArr, *sub)
		}
	}

	return subArr
}

// helper function for getSubstitutions()
func reduceTwo(x int, y int) *Substitution {
	if y == x + 1 {
		substitution := Substitution{x, y, nil}
		return &substitution
	}
	return nil
}
// helper function for getSubstitutions()
func reduceThree(x int, y int, z int) *Substitution {
	if y == x + 1 && z == y + 1 {
		substitution := Substitution{x, y, &z}
		return &substitution
	}
	return nil
}

func clusterSubstitutions(substitutionsArr []Substitution) []set.Set {
	/* clusters together consecutive numbers. In this example, 0-3 are consecutive
	 * [0 1 2 3 5 6....]
	 *
	 * and the first five substitutions all overlap at least one other of the five
	 *	[
	 *		[0, 1],
	 *		[0, 1, 2],
	 *		[1, 2],
	 *		[1, 2, 3],
	 *		[2, 3],
	 *
	 *		[5, 6],
	 *		...
	 *	]
	 * so they would be grouped into a cluster like this:
	 *	[
	 *		[
	 *			[0, 1],
	 *			[0, 1, 2],
	 *			[1, 2],
	 *			[1, 2, 3],
	 *			[2, 3],
	 *		],
	 *		...
	 *	]
	 */

	subClustersArr := []set.Set{}

	subArrLen := len(substitutionsArr)
	tmpClusterSet := set.New()
	for i, subOne := range substitutionsArr {
		if i == subArrLen - 1 {
			tmpClusterSet.Insert(subOne)
			subClustersArr = append(subClustersArr, *tmpClusterSet)
			continue
		}

		// if this element and the next element are overlapping, add them to the current cluster (set)
		subTwo := substitutionsArr[i+1]
		if areOverlapping(subOne, subTwo) {
			tmpClusterSet.Insert(subOne)
			tmpClusterSet.Insert(subTwo)
			continue
		}

		// if this element and the next element are not overlapping,
		// and if tmpClusterSet is not empty
		if tmpClusterSet.Len() > 0 {
			// then add this clusterSet to the array of clustersets
			subClustersArr = append(subClustersArr, *tmpClusterSet)
			// initialize the next empty clusterset
			tmpClusterSet = set.New()
			// and move onto the next substitution element
			continue
		}

		// if this element and the next element are not overlapping
		// and if the tmpClusterSet *is* empty, then this element should be in a set on its own
		tmpClusterSet.Insert(subOne)
		subClustersArr = append(subClustersArr, *tmpClusterSet)
		tmpClusterSet = set.New()
	}

	return subClustersArr
}


// helper function for getSubstitutions()
func areOverlapping(subOne Substitution, subTwo Substitution) bool {
	// subOne and subTwo are arrays of either 2 or 3 consecutive integers
	// subOne precedes subTwo, so if their integers overlap,
	// then the last integer in subOne will be found in subTwo and vice versa
	if subTwo.first == subOne.getLastInt() || subTwo.second == subOne.getLastInt() {
		return true
	}
	return false
}
