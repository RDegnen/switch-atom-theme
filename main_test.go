package main

import "testing"

func TestExtractThemeBlob(t *testing.T) {
	blob := extractThemeBlob("./test_files/test_package.json")
	label := blob.Contributes.Themes[0].Label
	if label != "Abyss" {
		t.Errorf("Incorrect label, got %s, wanted 'Abyss'", label)
	}
}

func TestExtractLabels(t *testing.T) {
	testTheme1 := theme{"Theme One"}
	testContributions1 := contributions{[]theme{testTheme1}}
	testThemeBlob1 := themeBlob{testContributions1}

	testTheme2 := theme{"Theme Two"}
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
