package main

import (
	// "fmt"
	"bufio"
	"os"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

var ROWS, COLS int
var source_file string
var offsetX, offsetY int
var text_buffer = [][]rune{} //rune = go alias for int32 characters
var currentCol, currentRow int

func insertRune(event termbox.Event) { //creates a new slice 1 rune larger to facilitate the new character.
	// Copies what's before the cursor in the previous slice into the new slice, inserts the character, copies what's after the cursor in the old slice into the new one.
	newLine := make([]rune, len(text_buffer[currentRow])+1)
	copy(newLine[:currentCol], text_buffer[currentRow][:currentCol])
	if event.Key == termbox.KeySpace {
		newLine[currentCol] = ' '
	} else if event.Key == termbox.KeyTab {
		newLine[currentCol] = ' '
	} else {
		newLine[currentCol] = event.Ch
	}
	copy(newLine[currentCol+1:], text_buffer[currentRow][currentCol:])
	text_buffer[currentRow] = newLine
	currentCol++
}

func saveFile() {
	f, err := os.Create(source_file)
	if err != nil {
		return
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	for _, line := range text_buffer {
		w.WriteString(string(line))
		w.WriteRune('\n')
	}
	w.Flush()
	text_buffer = [][]rune{{}}
	currentRow, currentCol = 0, 0
	offsetX, offsetY = 0, 0
}

func processKeypress(event termbox.Event) {
	if event.Key == termbox.KeyCtrlS {
		saveFile()
		return
	}

	switch event.Type {
	case termbox.EventKey:
		if event.Key == termbox.KeyEsc {
			termbox.Close()
			os.Exit(0)
		}
		if event.Key == termbox.KeyBackspace || event.Key == termbox.KeyBackspace2 {
			if currentCol > 0 {
				line := text_buffer[currentRow]
				newLine := append(line[:currentCol-1], line[currentCol:]...)
				text_buffer[currentRow] = newLine
				currentCol--
			} else if currentRow > 0 {
				prev := text_buffer[currentRow-1]
				line := text_buffer[currentRow]
				text_buffer[currentRow-1] = append(prev, line...)
				text_buffer = append(text_buffer[:currentRow], text_buffer[currentRow+1:]...)
				currentRow--
				currentCol = len(prev)
			}
		} else if event.Key == termbox.KeySpace || event.Key == termbox.KeyTab || event.Ch != 0 {
			insertRune(event)
		} else {
			switch event.Key {
			case termbox.KeyEnter:
				line := text_buffer[currentRow]
				before := line[:currentCol]
				after := line[currentCol:]
				text_buffer[currentRow] = before
				text_buffer = append(
					text_buffer[:currentRow+1],
					append([][]rune{after}, text_buffer[currentRow+1:]...)...,
				)
				currentRow++
				currentCol = 0

			case termbox.KeyArrowUp:
				if currentRow > 0 {
					currentRow--
					if currentCol > len(text_buffer[currentRow]) {
						currentCol = len(text_buffer[currentRow])
					}
				}
			case termbox.KeyArrowDown:
				if currentRow < len(text_buffer)-1 {
					currentRow++
					if currentCol > len(text_buffer[currentRow]) {
						currentCol = len(text_buffer[currentRow])
					}
				}
			case termbox.KeyArrowLeft:
				if currentCol > 0 {
					currentCol--
				} else if currentRow > 0 {
					currentRow--
					currentCol = len(text_buffer[currentRow])
				}
			case termbox.KeyArrowRight:
				if currentCol < len(text_buffer[currentRow]) {
					currentCol++
				} else if currentRow < len(text_buffer)-1 {
					currentRow++
					currentCol = 0
				}
			}
		}
		if currentCol < offsetX {
			offsetX = currentCol
		}
		if currentCol >= offsetX+COLS {
			offsetX = currentCol - COLS + 1
		}
		if currentRow < offsetY {
			offsetY = currentRow
		}
		if currentRow >= offsetY+ROWS {
			offsetY = currentRow - ROWS + 1
		}
	}
}

func displayTextBuffer() {
	for row := 0; row < ROWS; row++ {
		bufferRow := row + offsetY
		for col := 0; col < COLS; col++ {
			bufferCol := col + offsetX
			ch := ' '
			if bufferRow < len(text_buffer) && bufferCol < len(text_buffer[bufferRow]) {
				ch = text_buffer[bufferRow][bufferCol]
			}
			termbox.SetCell(col, row, ch, termbox.ColorGreen, termbox.ColorDefault)
		}
	}
}

func msg(col int, row int, fg termbox.Attribute, bg termbox.Attribute, message string) { //msg is just so i can display a message.
	for _, ch := range message { //loops through a set of characters
		termbox.SetCell(col, row, ch, fg, bg) //for each character, gives them a cell
		col += runewidth.RuneWidth(ch)        //gives the character the amount of cells it needs
	} //runewidth is used to see how many cells are needed for the character, as some unicode characters need more than 1.
}

func read_file(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		source_file = filename
		text_buffer = append(text_buffer, []rune{})
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		row := []rune(line)
		text_buffer = append(text_buffer, row)
	}
}

func run() {
	err := termbox.Init() // initialises termbox
	if err != nil {
		panic(err) // catches error
	}
	if len(os.Args) > 1 {
		read_file(os.Args[1])
	} else {
		source_file = "out.txt"
		text_buffer = append(text_buffer, []rune{})
	}
	w, h := termbox.Size()
	COLS = w
	ROWS = h // sets rows and columns to terminal width and height (MAY NOT BE ACCURATE)
	if COLS < 80 {
		COLS = 78
	}
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	// displayTextBuffer()
}

func main() {
	run()
	defer termbox.Close()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
	for {
		displayTextBuffer()
		termbox.SetCursor(currentCol-offsetX, currentRow-offsetY)
		termbox.Flush()           // synchronises the internal back buffer with the terminal, AKA, displays the stuff you drew up in the background
		ev := termbox.PollEvent() // variable to record keystroke
		processKeypress(ev)
	}
}
