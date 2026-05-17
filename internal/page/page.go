package page

import (
	"regexp"
	"strings"
)

// matches either a full tag (<...>) or a run of non-tag text
var re = regexp.MustCompile(`<[^>]*>|[^<]+`)

type Token interface{}

type Text struct {
	Data string
}

type Tag struct {
	Name string
}

func Lex(html string) []Token {
	var tokens []Token
	for _, chunk := range re.FindAllString(html, -1) {
		if strings.HasPrefix(chunk, "<") {
			tokens = append(tokens, Tag{Name: chunk[1 : len(chunk)-1]})
		} else {
			tokens = append(tokens, Text{Data: chunk})
		}
	}
	return tokens
}
