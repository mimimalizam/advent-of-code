package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Vec struct {
	X, Y int
}

func up(pos Vec) Vec {
	return Vec{X: pos.X, Y: pos.Y - 1}
}

func down(pos Vec) Vec {
	return Vec{X: pos.X, Y: pos.Y + 1}
}

func left(pos Vec) Vec {
	return Vec{X: pos.X - 1, Y: pos.Y}
}

func right(pos Vec) Vec {
	return Vec{X: pos.X + 1, Y: pos.Y}
}

var m [][]byte

type Portal struct {
	name          string
	warpPositions []Vec
}

func is(pos Vec, v byte) bool {
	if pos.Y < 0 || pos.X < 0 {
		return false
	}

	if pos.Y >= len(m) || pos.X >= len(m[0]) {
		return false
	}

	return m[pos.Y][pos.X] == v
}

func at(pos Vec) byte {
	return m[pos.Y][pos.X]
}

func isPortal(pos Vec) bool {
	if pos.Y < 0 || pos.X < 0 {
		return false
	}

	if pos.Y >= len(m) || pos.X >= len(m[0]) {
		return false
	}

	return m[pos.Y][pos.X] >= 'A' && m[pos.Y][pos.X] <= 'Z'
}

var portals = map[string]Portal{}

func load(filename string) {
	inputFile, _ := os.OpenFile(filename, os.O_RDONLY, os.ModePerm)
	reader := bufio.NewReader(inputFile)

	for {
		line, _, err := reader.ReadLine()
		if err != nil {
			break // EOF
		}

		m = append(m, line)
	}

	for i, line := range m {
		for j, c := range line {
			if c >= 'A' && c <= 'Z' {
				pos := Vec{X: j, Y: i}

				if !isPortal(pos) {
					continue
				}

				name := ""
				warpPos := Vec{}

				pUp := up(pos)
				pDown := down(pos)
				pLeft := left(pos)
				pRight := right(pos)

				if !is(pDown, '.') && !is(pUp, '.') && !is(pLeft, '.') && !is(pRight, '.') {
					continue
				}

				if isPortal(pUp) {
					name = string(at(pUp)) + string(c)
					warpPos = pDown
				}

				if isPortal(pDown) {
					name = string(c) + string(at(pDown))
					warpPos = pUp
				}

				if isPortal(pLeft) {
					name = string(at(pLeft)) + string(c)
					warpPos = pRight
				}

				if isPortal(pRight) {
					name = string(c) + string(at(pRight))
					warpPos = pLeft
				}

				portal := portals[name]
				portal.warpPositions = append(portal.warpPositions, warpPos)

				portals[name] = portal
			}
		}
	}

	fmt.Println(portals)
}

func show(pos Vec) {
	res := ""

	for i, line := range m {
		for j, c := range line {
			if j == pos.X && i == pos.Y {
				res += fmt.Sprint("\033[31m@\033[0m")
			} else {
				res += fmt.Sprint(string(c))
			}
		}

		res += "\n"
	}

	time.Sleep(500 * time.Millisecond)

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()

	fmt.Println(res)
}

var visited = map[Vec]bool{}

func solve(pos Vec, steps int) (int, bool) {
	fmt.Println(pos)

	show(pos)

	visited[pos] = true

	minOk := false
	min := 1000000001

	for _, nextPos := range []Vec{up(pos), down(pos), left(pos), right(pos)} {
		if at(nextPos) == 'Z' {
			return steps, true
		}

		if isPortal(nextPos) {
			n1 := nextPos
			n2 := Vec{X: nextPos.X + nextPos.X - pos.X, Y: nextPos.Y + nextPos.Y - pos.Y}

			name := ""

			if n1.X == n2.X {
				if n1.Y < n2.Y {
					name = string(at(n1)) + string(at(n2))
				} else {
					name = string(at(n2)) + string(at(n1))
				}
			}

			if n1.Y == n2.Y {
				if n1.X < n2.X {
					name = string(at(n1)) + string(at(n2))
				} else {
					name = string(at(n2)) + string(at(n1))
				}
			}

			fmt.Println(n1, n2)
			fmt.Println(name)
			fmt.Println(portals[name])

			if name == "AA" {
				continue
			}

			if portals[name].warpPositions[0] == pos {
				nextPos = portals[name].warpPositions[1]
			}

			if portals[name].warpPositions[1] == pos {
				nextPos = portals[name].warpPositions[0]
			}
		}

		if visited[nextPos] || at(nextPos) == ' ' || at(nextPos) == '#' {
			continue
		}

		steps, ok := solve(nextPos, steps+1)
		if ok {
			if steps < min {
				min = steps
				minOk = true
			}
		}
	}

	visited[pos] = false

	return min, minOk
}

func main() {
	load("input2.txt")

	pos := portals["AA"].warpPositions[0]
	show(pos)

	fmt.Println(pos)

	steps, ok := solve(pos, 0)

	fmt.Println(steps, ok)
}