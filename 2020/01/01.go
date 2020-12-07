package main

import (
    "bufio"
    "flag"
    "fmt"
    "syscall"
)

var NUMBER_OF_INTEGERS_TO_FIND = 3
var SUM_OF_INTS = 2020
var VERBOSE = false

var TOTAL_ITERATIONS = 0

// find the $NUMBER_OF_INTEGERS_TO_FIND integers in the input file that sum to $SUM_OF_INTS
func main() {
    vptr := flag.Bool("v", false, "verbose")
    flag.Parse()
    VERBOSE = *vptr

    f := getFile()
    defer f.Close()
    dataArr := makeDataArr(bufio.NewScanner(f))
    mainDataArr := &dataArr

    // level starts at 1
    n := 1
    // mthInt is 0
    mthInt := 0
    // sumOfInts is 2020
    sumOfInts := SUM_OF_INTS

    // 1. determine neededSum
    neededSumOrInt := sumOfInts - mthInt

    if solutionArr := getIntArrayRecursively(neededSumOrInt, mthInt, mainDataArr, n); solutionArr != nil {
        product := 1
        for _, num := range *solutionArr {
            product *= num
        }
        fmt.Printf("%d integers that sum to %d and multiply to %d: %v\n",
            NUMBER_OF_INTEGERS_TO_FIND, SUM_OF_INTS, product, *solutionArr)
        fmt.Printf("TOTAL_ITERATIONS: %d\n", TOTAL_ITERATIONS)
        syscall.Exit(0)
    }
}

func getIntArrayRecursively(sumOfInts int, mthInt int, dataArr *[]int, n int) *[]int {
    // just for debugging
    tabs, ordinal := getFormatters(n)

    // 1. determine neededSum
    neededSumOrInt := sumOfInts - mthInt

    for i, nthInt := range (*dataArr)[:len(*dataArr)-1] {
        TOTAL_ITERATIONS++

        // 2. don't run in these cases
        if nthInt > neededSumOrInt {
            continue
        }
        if VERBOSE {
            fmt.Printf("%d|%sAt index %d, checking if %d could be the %d%s int ",
                n, tabs, i, nthInt, n, ordinal)
            //fmt.Printf("that makes %d + %d + ... = %d\n", mthInt, nthInt, sumOfInts)
            sumStr := getSumString(n, mthInt, nthInt, sumOfInts)
            fmt.Printf("that makes %s\n", sumStr)
        }



        if n < NUMBER_OF_INTEGERS_TO_FIND {
            // 3. if there are ${NUMBER_OF_INTEGERS_TO_FIND - n} more integers
            // that, with mthInt and nthInt, sum to $SUM_OF_INTS, add to array and return
            dataArrSlice := (*dataArr)[i+1:]

            if intArrPtr := getIntArrayRecursively(neededSumOrInt, nthInt, &dataArrSlice, n+1); intArrPtr != nil {
                *intArrPtr = append(*intArrPtr, nthInt)
                if VERBOSE {
                   fmt.Printf("%d|%sreturning %d integers: %v\n", n, tabs, len(*intArrPtr), *intArrPtr)
                }
                return intArrPtr
            }
        } else { // checking for final integer
            // 3. if nthInt is a match, return nthInt as the only element in an array of integers
            if nthInt == neededSumOrInt {
                intArrPtr := &[]int{nthInt}

                if VERBOSE {
                   fmt.Printf("%d|%sreturning %d integers: %v\n", n, tabs, len(*intArrPtr), *intArrPtr)
                }
                return intArrPtr
            }
        }
    }
    return nil
}
