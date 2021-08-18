package util

import (
	"strings"
)

// Constructor until
type Constructor struct {
	FTSearch string
}

// FTSQuery for construct search in full text search tsquery
func (c *Constructor) FTSQuery() string {
	var newText strings.Builder
	text := strings.TrimSpace(c.FTSearch)
	s := strings.Split(text, " ")

	for i := 0; i < len(s); i++ {
		word := s[i]
		newText.WriteString(word + ":*")

		if i != len(s)-1 {
			newText.WriteString(" & ")
		}
	}
	return newText.String()
}
