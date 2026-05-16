package filereader

import (
	"fmt"
	"gobrowser/internal/urlparser"
	"net/http"
	"os"
	"strings"
)

func Read(url urlparser.FileURL) (string, error) {
	data, err := os.ReadFile(url.Path)
	if err != nil {
		return "", err
	}
	mimeType := http.DetectContentType(data)
	if !strings.HasPrefix(mimeType, "text/") {
		return "", fmt.Errorf("unsupported file type: %s", mimeType)
	}
	return string(data), nil
}
