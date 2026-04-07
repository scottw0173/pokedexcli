package main

import (
	"strings"
)

func cleanInput(text string) []string {
	cleanedText := strings.Fields(strings.ToLower(text))
	return cleanedText
}
