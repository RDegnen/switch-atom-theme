package main

import termbox "github.com/nsf/termbox-go"

type selectionOptions struct {
	data          []string
	cursorBoffset int
}

func (so *selectionOptions) render() {
	const coldef = termbox.ColorDefault
	for i, s := range options.data {
		tbprint(0, i+listOffset, coldef, coldef, s)
	}
}

func (so *selectionOptions) moveCursor(boffset int) {
	so.cursorBoffset = boffset
	termbox.SetCursor(0, boffset)
}

func (so *selectionOptions) moveCursorUp() {
	if so.cursorBoffset == listOffset {
		return
	}
	so.moveCursor(so.cursorBoffset - 1)
}

func (so *selectionOptions) moveCursorDown() {
	if so.cursorBoffset == len(so.data)+listOffset-1 {
		return
	}
	so.moveCursor(so.cursorBoffset + 1)
}

func (so *selectionOptions) selectOption() string {
	value := so.data[so.cursorBoffset-listOffset]
	return value
}
