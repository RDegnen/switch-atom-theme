package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func TestExtractThemeBlob(t *testing.T) {
	blob := extractThemeBlob("./test_files/test_package.json")
	label := blob.Contributes.Themes[0].Label
	if label != "Abyss" {
		t.Errorf("Incorrect label, got %s, wanted 'Abyss'", label)
	}
}

func TestExtractLabels(t *testing.T) {
	testTheme1 := theme{ID: "Theme One", Label: "Theme One"}
	testContributions1 := contributions{[]theme{testTheme1}}
	testThemeBlob1 := themeBlob{testContributions1}

	testTheme2 := theme{ID: "Theme Two", Label: "Theme Two"}
	testContributions2 := contributions{[]theme{testTheme2}}
	testThemeBlob2 := themeBlob{testContributions2}

	themes := []themeBlob{testThemeBlob1, testThemeBlob2}
	labels := extractLabels(themes)
	if labels[0] != "Theme One" {
		t.Errorf("Incorrect label, got %s, wanted Theme One", labels[0])
	}
	if labels[1] != "Theme Two" {
		t.Errorf("Incorrect label, got %s, wanted Theme Two", labels[1])
	}
}

func TestChangeVscodeSettings(t *testing.T) {
	testSettings := changeVscodeSettings("New Theme", "./test_files/test_vscode_settings.json")
	if testSettings["workbench.colorTheme"] != "New Theme" {
		t.Errorf("Incorrect theme, got %s, wanted Old Theme", testSettings["workbench.colorTheme"])
	}
}

func TestMutateVscodeSettingsFile(t *testing.T) {
	var expectedJSON interface{}
	var recievedJSON interface{}
	oldFilePath := "./test_files/test_vscode_settings.json"
	newFilePath := "./test_files/new_settings.json"
	testSettings := changeVscodeSettings("New Theme", oldFilePath)
	json.Unmarshal(
		[]byte(`{"editor.tabSize":2,"window.zoomLevel":0,"workbench.colorTheme":"New Theme","workbench.startupEditor":"newUntitledFile"}`),
		&expectedJSON,
	)

	mutateVscodeSettingsFile(testSettings, newFilePath)
	testData, err := ioutil.ReadFile(newFilePath)
	if err != nil {
		fmt.Println(err)
	}
	json.Unmarshal(testData, &recievedJSON)

	if !reflect.DeepEqual(recievedJSON, expectedJSON) {
		t.Errorf("Incorrect mutation of vscode test file,\n expected %s,\n to equal %s", recievedJSON, expectedJSON)
	}

	os.Remove(newFilePath)
}
