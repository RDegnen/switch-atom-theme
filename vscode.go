package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	termbox "github.com/nsf/termbox-go"
)

type themeBlob struct {
	Contributes contributions
}

type contributions struct {
	Themes []theme
}

type theme struct {
	ID    string
	Label string
}

func getThemeDirectories(inPath string) []string {
	files, err := filepath.Glob(fmt.Sprintf("%s/*theme*/package.json", inPath))
	if err != nil {
		fmt.Println(err)
	}
	return files
}

func getTotalLength(arr [][]string) int {
	length := 0
	for i := range arr {
		length += len(arr[i])
	}
	return length
}

func flattenString(arr [][]string) []string {
	newArr := make([]string, 0, getTotalLength(arr))
	for i := range arr {
		for _, v := range arr[i] {
			newArr = append(newArr, v)
		}
	}
	return newArr
}

func extractThemeBlob(filePath string) themeBlob {
	var blob themeBlob

	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(data, &blob); err != nil {
		fmt.Println(err)
	}
	return blob
}

func idOrLabel(t theme) string {
	if t.ID != "" {
		return t.ID
	}
	return t.Label
}

func extractLabels(blobs []themeBlob) []string {
	labels := make([]string, 0, len(blobs))
	for i := range blobs {
		for _, v := range blobs[i].Contributes.Themes {
			labels = append(labels, idOrLabel(v))
		}
	}
	return labels
}

func changeVscodeSettings(theme string, filepath string) map[string]interface{} {
	var settings interface{}
	data, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Println(err)
	}

	if err := json.Unmarshal(data, &settings); err != nil {
		fmt.Println(err)
	}

	settingsMap := settings.(map[string]interface{})
	settingsMap["workbench.colorTheme"] = theme
	return settingsMap
}

func mutateVscodeSettingsFile(newSettings map[string]interface{}, filePath string) {
	data, err := json.Marshal(newSettings)
	if err != nil {
		fmt.Println(err)
	}

	if err := ioutil.WriteFile(filePath, data, 0644); err != nil {
		fmt.Println(err)
	}
}

func vscode() {
	vscodeInstalledExtensions := fmt.Sprintf("%s/.vscode/extensions", os.Getenv("HOME"))
	vscodeDefaultExtensions := "/Applications/Visual Studio Code.app/Contents/Resources/app/extensions"
	vscodeSettingsFile := fmt.Sprintf("%s/Library/Application Support/Code/User/settings.json", os.Getenv("HOME"))
	packageFiles :=
		flattenString(
			[][]string{
				getThemeDirectories(vscodeDefaultExtensions),
				getThemeDirectories(vscodeInstalledExtensions),
			},
		)

	themeBlobs := make([]themeBlob, 0, len(packageFiles))
	for _, v := range packageFiles {
		themeBlobs = append(themeBlobs, extractThemeBlob(v))
	}

	themes := extractLabels(themeBlobs)
	options.data = themes

	redrawAll()

vscodeloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break vscodeloop
			case termbox.KeyArrowUp:
				options.moveCursorUp()
			case termbox.KeyArrowDown:
				options.moveCursorDown()
			case termbox.KeyEnter:
				newSettings := changeVscodeSettings(options.selectOption(), vscodeSettingsFile)
				mutateVscodeSettingsFile(newSettings, vscodeSettingsFile)
				break vscodeloop
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		redrawAll()
	}
}
