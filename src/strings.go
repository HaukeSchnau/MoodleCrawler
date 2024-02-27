package main

import (
	"regexp"
	"strings"
)

func matchFolder(name string) (folder string, ok bool) {
	lectureRegex := regexp.MustCompile(`(?i)Vorlesung|Folien|Kapitel|Aufzeichnung|LE[0-9]{1,2}|Woche [0-9]{1,2}|Slides`)
	exerciseRegex := regexp.MustCompile(`(?i)Übung|Aufgabe|Blatt|Exercise`)
	solutionRegex := regexp.MustCompile(`(?i)Lösung`)

	if lectureRegex.MatchString(name) {
		return "Vorlesung", true
	}

	if solutionRegex.MatchString(name) {
		return "Lösungen", true
	}

	if exerciseRegex.MatchString(name) {
		return "Aufgaben", true
	}

	return "nil", false
}

func makePath(outDir string, name string) string {
	normalizedTitle := normalizeUmlauts(name)
	folder, hasFolder := matchFolder(normalizedTitle)

	escapedTitle := strings.ReplaceAll(normalizedTitle, "/", "_")
	pathedTitle := strings.ReplaceAll(escapedTitle, " - ", "/")

	if hasFolder {
		return outDir + "/" + folder + "/" + pathedTitle
	}

	return outDir + "/" + pathedTitle
}

func normalizeUmlauts(str string) string {
	umlautMap := map[string]string{
		"ae":    "ä",
		"oe":    "ö",
		"ue":    "ü",
		"Ae":    "Ä",
		"Oe":    "Ö",
		"Ue":    "Ü",
		"AE":    "Ä",
		"OE":    "Ö",
		"UE":    "Ü",
		"&amp;": "&",
	}

	for key, value := range umlautMap {
		str = strings.ReplaceAll(str, key, value)
	}

	return str
}
