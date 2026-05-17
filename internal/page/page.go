package page

import "strings"

type Token interface{}

type Text struct {
	Data string
}

func Lex(html string) []Token {
	var buf strings.Builder
	inTag := false

	for _, r := range html {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			buf.WriteRune(r)
		}
	}
	return []Token{Text{Data: buf.String()}}
}
