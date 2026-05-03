package renderer

import "strings"

func Render(html string) (string, error) {
	var textBuilder strings.Builder
	inTag := false
	
	for _, r := range html {
		if r == '<' {
			inTag = true
		} else if r == '>' {
			inTag = false
		} else if !inTag {
			textBuilder.WriteRune(r)
		}
	}
	
	return textBuilder.String(), nil
}
