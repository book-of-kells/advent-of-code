package main

import (
	"log"
	"math"
)

type WayPoint struct {
	x int
	y int
	ship *Ship
}

func (way *WayPoint) getPythagoreanDistanceFromShip () float64 {
	squares := (way.x - way.ship.x)^2 + (way.y - way.ship.y)^2
	ret := math.Sqrt(float64(squares))
	return ret
}

func (way *WayPoint) rotate90(dir string) {

	diffY := way.y - way.ship.y
	diffX := way.x - way.ship.x
	switch dir{
	case "L":
		way.x = way.ship.x - diffY
		way.y = way.ship.y + diffX
	case "R":
		way.x = way.ship.x + diffY
		way.y = way.ship.y - diffX

	default:
		log.Fatalf("could not rotate with direction %s\n", dir)
	}
}

func (way *WayPoint) rotate (command Command) {
	dir := command.cmd
	degrees := command.num

	switch degrees {
	case 180:
		way.x = way.ship.x - (way.x - way.ship.x)
		way.y = way.ship.y - (way.y - way.ship.y)
	case 90:
		way.rotate90(dir)
	case 270:
		switch dir {
		case "L":
			way.rotate90("R")
		case "R":
			way.rotate90("L")
		}
	}
}

// todo fix
func (way *WayPoint) move (command Command) {

	moveIn := command.cmd
	if moveIn == "F" {
		/*
		Action `F` means to move forward to the waypoint a number of times
		equal to the given value.

		The waypoint starts 10 units east and 1 unit north relative to the ship.
		The waypoint is relative to the ship;
		that is, if the ship moves, the waypoint moves with it.

		`F10` moves the ship to the waypoint 10 times
		(a total of 100 units east and 10 units north),
		leaving the ship at east 100, north 10.
		The waypoint stays 10 units east and 1 unit north of the ship.
		*/
		diffX := way.x - way.ship.x  // 110 - 100 = 10

		moveX := diffX * command.num // 70
		if moveX < 0 {
			way.ship.move(Command{cmd: "W", num: int(math.Abs(float64(moveX)))})
		} else {
			way.ship.move(Command{cmd: "E", num: moveX})
		}

		diffY := way.y - way.ship.y
		moveY := diffY * command.num
		if moveY < 0 {
			way.ship.move(Command{cmd: "S", num: int(math.Abs(float64(moveY)))})
		} else {
			way.ship.move(Command{cmd: "N", num: moveY})
		}

		way.x = way.ship.x + diffX
		way.y = way.ship.y + diffY
	}

	switch moveIn{
	case "E":
		way.x += command.num
	case "N":
		way.y += command.num
	case "W":
		way.x -= command.num
	case "S":
		way.y -= command.num
	}
}