package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"syscall"
)

var VERBOSE = false

func main() {
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile()
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))

	totalNop := 0 // 61 nop
	nopMap := make(map[int]LineCommand)
	totalJmp := 0 // 229 jmp
	jmpMap := make(map[int]LineCommand)

	for i, elem := range dataArr {
		if elem.command == "jmp" {
			totalJmp++
			jmpMap[i] = *elem
		}
		if elem.command == "nop" {
			totalNop++
			nopMap[i] = *elem
		}
		fmt.Printf("\nelement %d\n%v\n", i, elem)
	}


	// first try replacing a nop command with a jmp command
	for nopIdx, nopCmd := range nopMap {
		storedNopCmd := dataArr[nopIdx]
		newJmpCmd := &LineCommand{nopIdx, "jmp", nopCmd.direction, nopCmd.num, false}
		dataArr[nopIdx] = newJmpCmd
		if isSuccess := executeProgram(dataArr); *isSuccess == true {
			fmt.Printf("Success! Changed nop %v for jmp %v\n", storedNopCmd, newJmpCmd)
			syscall.Exit(0)
		}

		// restore the swapped nop command and try another
		dataArr[nopIdx] = storedNopCmd
		resetAll(dataArr)
	}


	// next try replacing a jmp command with a nop command
	for jmpIdx, jmpCmd := range jmpMap {
		storedJmpCmd := dataArr[jmpIdx]
		newNopCmd := &LineCommand{jmpIdx, "nop", jmpCmd.direction, jmpCmd.num, false}
		dataArr[jmpIdx] = newNopCmd
		if isSuccess := executeProgram(dataArr); *isSuccess == true {
			fmt.Printf("Success! Changed jmp %v for nop %v\n", storedJmpCmd, newNopCmd)
			syscall.Exit(0)
		}

		// restore the swapped jmp command and try another
		dataArr[jmpIdx] = storedJmpCmd
		resetAll(dataArr)
	}
}

func executeProgram(cmdArr []*LineCommand) *bool {
	isSuccess := false
	accValue := 0
	command := cmdArr[0]

	for command.executed == false {
		fmt.Printf("accValue=%d | %v\n", accValue, command)
		if command.executed == true {
			fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
			return &isSuccess
		}
		switch command.command {
		case "acc":
			if command.direction == "+" {
				accValue += command.num
			} else if command.direction == "-" {
				accValue -= command.num
			} else {
				log.Fatalf("direction should be '+' or '-' but instead was %s\n", command.direction)
			}
			command.executed = true
			if command.index+1 == len(cmdArr) {
				fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
				return &isSuccess
			}
			command = cmdArr[command.index+1]
		case "jmp":
			command.executed = true
			if command.direction == "+" {
				command = cmdArr[command.index + command.num]
			} else if command.direction == "-" {
				command = cmdArr[command.index - command.num]
			} else {
				log.Fatalf("direction should be '+' or '-' but instead was %s\n", command.direction)
			}

		case "nop":
			command.executed = true
			if command.index+1 == len(cmdArr) {
				fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
				return &isSuccess
			}
			command = cmdArr[command.index+1]
		}
		if command.index == len(cmdArr) - 1 { // executing this at the end because we know the final element isn't `acc`
			fmt.Println("************************************************************************")
			fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
			fmt.Println("************************************************************************")
			isSuccess = true
			return &isSuccess
		}
	}
	return &isSuccess
}

/*
func executeProgramBackwards(cmdArr []*LineCommand) {
	accValue := 0
	command := cmdArr[len(cmdArr)-1]

	for command.executed == false {
		fmt.Printf("accValue=%d | %v\n", accValue, command)
		if command.executed == true {
			fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
			syscall.Exit(0)
		}
		switch command.command {
		case "acc":
			if command.direction == "+" {
				accValue -= command.num
			} else if command.direction == "-" {
				accValue += command.num
			} else {
				log.Fatalf("direction should be '+' or '-' but instead was %s\n", command.direction)
			}
			command.executed = true
			if command.index-1 == 0 {
				fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
				syscall.Exit(0)
			}
			command = cmdArr[command.index-1]
		case "jmp":
			command.executed = true
			if command.direction == "+" {
				if command.index < command.num {
					_ = fmt.Errorf("About to jump to a negative index with command %v. Changing this jmp to a nop\n", command)
					cmdArr[command.index] = &LineCommand{command.index, "nop", "+", 0, true}
					command = cmdArr[command.index - 1]
				} else {
					command = cmdArr[command.index - command.num]
				}

			} else if command.direction == "-" {
				if (command.index + command.num) > (len(cmdArr) - 1) {
					_ = fmt.Errorf("About to jump out of bounds with command %v. Changing this jmp to a nop\n", command)
					cmdArr[command.index] = &LineCommand{command.index, "nop", "+", 0, true}
					command = cmdArr[command.index - 1]
				} else {
					command = cmdArr[command.index+command.num]
				}
			} else {
				log.Fatalf("direction should be '+' or '-' but instead was %s\n", command.direction)
			}

		case "nop":
			command.executed = true
			if command.index-1 == 0 {
				fmt.Printf("****** DONE!! cmd %v and accumulated value is %d ******\n", command, accValue)
				syscall.Exit(0)
			}
			command = cmdArr[command.index-1]
		}
	}
}
 */
