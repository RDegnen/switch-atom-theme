package main

import (
	termbox "github.com/nsf/termbox-go"
)

const (
	listOffset = 2
)

var (
	current        string
	curev          termbox.Event
	options        selectionOptions
	selectedOption string
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func redrawAll() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	tbprint(0, 0, termbox.ColorCyan, coldef, "Select a theme. Press 'ESC' to quit, use arrows to navigate")
	options.render()
	termbox.Flush()
}

func selectEditor(value string) {
	switch value {
	case "atom":
		atom()
	case "vscode":
		vscode()
	}
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	options.cursorBoffset = listOffset
	options.data = []string{"vscode", "atom"}
	termbox.SetCursor(0, options.cursorBoffset)

	redrawAll()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyArrowUp:
				options.moveCursorUp()
			case termbox.KeyArrowDown:
				options.moveCursorDown()
			case termbox.KeyEnter:
				selectEditor(options.selectOption())
				break mainloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
}
