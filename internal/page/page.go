package page

import (
	"html"
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

func Lex(body string) []Token {
	var tokens []Token
	for _, chunk := range re.FindAllString(body, -1) {
		if strings.HasPrefix(chunk, "<") {
			tokens = append(tokens, Tag{Name: chunk[1 : len(chunk)-1]})
		} else {
			tokens = append(tokens, Text{Data: html.UnescapeString(chunk)})
		}
	}
	return tokens
}
