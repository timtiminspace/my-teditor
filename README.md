# Terminal Text Editor (Go + Termbox)

A lightweight terminal-based text editor written in Go using termbox-go.
It supports basic text editing, file saving, cursor navigation, and a simple search bar (so far)

## Features

- Open and edit text files
- Insert and delete characters
- Multi-line eiditing
- Arrow key navigation
- Scrollable viewport
- Status bar with file name and cursor position
- File saving (Ctrl + S)
- Search mode (Ctrl + LsqBrckt)
- Terminal resize handling
- Unicode character width support

## Keybindings

Arrow keys -> move cursor
Enter -> new line
Backspace -> delete character/merge lines
Space/typing -> Insert character

Ctrl + S -> save file
Ctrl + LsqBrckt -> enter search mode (Ctrl + F functionality)

Esc -> quit editor/cancel search

## Running the Editor

1. Install dependencies

- go get github.com/nsf/termbox-go
- go get github.com/mattn/go-runewidth

2. Run with a file

- go run main.go filename.txt

If no file is provided, it opens out.txt

## How It Works

- Text is stored as a 2D slice of runes ([][]rune)
- The screen renders a viewport using scroll offsets
- Cursor position is tracked separately from viewport
- A status bar shows editor info
- Search mode temporarily captures keyboard input

## Limitations

- No syntax highlighting
- No undo/redo
- No mouse support
- No file explorer

## Requirements

- Go 1.18+
- Terminal that supports termbox

Feel free to use, modify and learn from this (: