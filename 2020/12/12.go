package main

import (
	"bufio"
	"flag"
	"fmt"
	"math"
)

var VERBOSE = false


type Ship struct {
	facing string
	x int
	y int
}

func (s *Ship) rotate (command Command) {
	dir := command.cmd
	degrees := command.num
	dirMap := map[string]int{
		"E": 0,
		"N": 90,
		"W": 180,
		"S": 270,
	}

	currDirection := s.facing
	currDirectionInt := dirMap[currDirection]
	switch dir {
	case "L":
		currDirectionInt += degrees
	case "R":
		currDirectionInt -= degrees
	}

	if currDirectionInt < 0 {
		currDirectionInt += 360
	}
	if currDirectionInt > 270 {
		currDirectionInt -= 360
	}

	intMap := map[int]string{
		0: "E",
		90: "N",
		180: "W",
		270: "S",
	}

	s.facing = intMap[currDirectionInt]
}

func (s *Ship) move (command Command) {

	moveIn := command.cmd
	if moveIn == "F" {
		moveIn = s.facing
	}
	switch moveIn{
	case "E":
		s.x += command.num
	case "N":
		s.y += command.num
	case "W":
		s.x -= command.num
	case "S":
		s.y -= command.num
	}
}


func main() {
	fptr := flag.String("file", "input.txt", "file path to read from")
	vptr := flag.Bool("v", false, "verbose")
	flag.Parse()
	VERBOSE = *vptr

	f := getFile(fptr)
	defer f.Close()
	dataArr := makeDataArr(bufio.NewScanner(f))

	ship := Ship{facing: "E", x: 0, y: 0}
	waypoint := WayPoint{x: 10, y:1, ship: &ship}
	fmt.Printf("ship is %v\n", ship)

	for i, command := range dataArr {
		fmt.Printf("\ncommand %d\n%v\n", i, command)

		switch command.cmd{
		case "L":
			waypoint.rotate(command)
		case "R":
			waypoint.rotate(command)
		default:
			waypoint.move(command)
		}
		fmt.Printf("Now ship is %v and waypoint is %v\n", ship, waypoint)
	}

	fmt.Printf("Manhattan distance: %d\n", int(math.Abs(float64(ship.x)) + math.Abs(float64(ship.y))))

}