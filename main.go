package main

import (
	"fmt"
	"io"
	"math/rand"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"time"
)

var speed = 3

func main() {
	setup()

	text := getInput()
	numLines := strings.Count(text, "\n")
	charMap := make(map[int]map[int]bool)
	odds := 100

	for {
		x := 0
		y := 0
		unseen := false

		for _, r := range text {
			// if we're at the end of the line, move to the next line
			if r == '\n' {
				x = 0

				y++

				fmt.Fprint(os.Stdout, "\n")

				continue
			}

			// if we've already seen this character, print it and move on
			if _, ok := charMap[y]; ok {
				if val, ok := charMap[y][x]; ok && val {
					x++

					fmt.Fprint(os.Stdout, string(r))

					continue
				}
			}

			unseen = true

			char := " "

			if oddsRandom(odds) {
				// minimum of 3 milliseconds required to make the animation look good
				time.Sleep(time.Duration(speed) * time.Millisecond)

				char = string(r)

				if _, ok := charMap[y]; !ok {
					charMap[y] = make(map[int]bool)
				}

				charMap[y][x] = true
			}

			x++

			fmt.Fprint(os.Stdout, char)
		}

		if !unseen {
			onDone(0)
		}

		odds--

		moveCursorUp(numLines)
	}
}

// setup is called at the beginning of the program.
func setup() {
	rand.Seed(time.Now().UnixNano())

	// get fist argument
	if len(os.Args) > 1 {
		speedArg := os.Args[1]

		num, err := strconv.Atoi(speedArg)
		if err == nil {
			speed = num
		}
	}

	hideCursor()

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	go func() {
		<-signalChan
		onDone(1)
	}()
}

// onDone is called when the program is done.
func onDone(code int) {
	showCursor()
	os.Exit(code)
}

// getInput reads from stdin and returns the input as a string.
func getInput() string {
	bytes, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintln(os.Stderr, "error reading from stdin:", err)
		onDone(1)
	}

	return string(bytes)
}

func oddsRandom(num int) bool {
	if num <= 0 {
		return true
	}

	return rand.Intn(num) == 1
}

func moveCursorUp(count int) {
	fmt.Printf("\033[%dA", count)
}

func hideCursor() {
	fmt.Print("\033[?25l")
}

func showCursor() {
	fmt.Print("\033[?25h")
}
