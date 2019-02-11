package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	termbox "github.com/nsf/termbox-go"
)

func contains(s [2]string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func changeAtomConfig(value string, themeOffset int, filePath string) string {
	file, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	lines := strings.Split(string(file), "\n")
	for i, line := range lines {
		if strings.Contains(line, "theme") {
			lines[i+themeOffset] = fmt.Sprintf(`      "%s"`, value)
		}
	}
	return strings.Join(lines, "\n")
}

func mutateAtomConfigFile(filePath string, data string) {
	if err := ioutil.WriteFile(filePath, []byte(data), 0644); err != nil {
		fmt.Println(err)
	}
}

func atom() {
	atomConfigDir := fmt.Sprintf("%s/.atom", os.Getenv("HOME"))
	configFile := fmt.Sprintf("%s/config.cson", atomConfigDir)
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

	options.data = syntax
	themeOffset := 2

	redrawAll()

atomloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break atomloop
			case termbox.KeyArrowUp:
				options.moveCursorUp()
			case termbox.KeyArrowDown:
				options.moveCursorDown()
			case termbox.KeyEnter:
				data := changeAtomConfig(options.selectOption(), themeOffset, configFile)
				mutateAtomConfigFile(configFile, data)
				if themeOffset == 1 {
					break atomloop
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
