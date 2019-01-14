package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/nsf/termbox-go"
)

const (
	listOffset = 2
)

var (
	current        string
	curev          termbox.Event
	options        selectOptions
	selectedOption string
	atomConfigDir  = fmt.Sprintf("%s/.atom", os.Getenv("HOME"))
)

func tbprint(x, y int, fg, bg termbox.Attribute, msg string) {
	for _, c := range msg {
		termbox.SetCell(x, y, c, fg, bg)
		x++
	}
}

func changeAtomConfig(value string, themeOffset int) {
	file, err := ioutil.ReadFile(fmt.Sprintf("%s/config.cson", atomConfigDir))
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		if strings.Contains(line, "theme") {
			lines[i+themeOffset] = fmt.Sprintf(`      "%s"`, value)
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(fmt.Sprintf("%s/config.cson", atomConfigDir), []byte(output), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func contains(s [2]string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

type selectOptions struct {
	data          []string
	cursorBoffset int
}

func (dl *selectOptions) render() {
	const coldef = termbox.ColorDefault
	for i, s := range options.data {
		tbprint(0, i+listOffset, coldef, coldef, s)
	}
}

func (dl *selectOptions) moveCursor(boffset int) {
	dl.cursorBoffset = boffset
	termbox.SetCursor(0, boffset)
}

func (dl *selectOptions) moveCursorUp() {
	if dl.cursorBoffset == listOffset {
		return
	}
	dl.moveCursor(dl.cursorBoffset - 1)
}

func (dl *selectOptions) moveCursorDown() {
	if dl.cursorBoffset == len(dl.data)+listOffset-1 {
		return
	}
	dl.moveCursor(dl.cursorBoffset + 1)
}

func (dl *selectOptions) selectOption() string {
	value := dl.data[dl.cursorBoffset-listOffset]
	return value
}

func redrawAll() {
	const coldef = termbox.ColorDefault
	termbox.Clear(coldef, coldef)
	tbprint(0, 0, termbox.ColorCyan, coldef, "Select a theme. Press 'ESC' to quit, use arrows to navigate")
	options.render()
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	files, err := ioutil.ReadDir(fmt.Sprintf("%s/packages", atomConfigDir))
	if err != nil {
		fmt.Println(err)
	}

	syntax := make([]string, 0, len(files))
	ui := make([]string, 0, len(files))
	builtInUI := [4]string{"atom-dark-ui", "atom-light-ui", "one-dark-ui", "one-light-ui"}
	UIToIgnore := [2]string{"atom-ide-ui", "linter-ui-default"}

	for _, f := range files {
		if strings.Contains(f.Name(), "syntax") {
			syntax = append(syntax, f.Name())
		}
		if strings.Contains(f.Name(), "ui") && !contains(UIToIgnore, f.Name()) {
			ui = append(ui, f.Name())
		}
	}

	for _, value := range builtInUI {
		ui = append(ui, value)
	}

	options.cursorBoffset = listOffset
	options.data = syntax
	termbox.SetCursor(0, options.cursorBoffset)

	redrawAll()

	themeOffset := 2
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
				changeAtomConfig(options.selectOption(), themeOffset)
				if themeOffset == 1 {
					break mainloop
				}
				themeOffset = 1
				options.data = ui
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
}
