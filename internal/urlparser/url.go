package urlparser

import (
	"errors"
	"strings"
)

type URL struct {
	Scheme string
	Host   string
	Path   string
}

func Parse(uri string) (URL, error) {
	scheme, urn, found := strings.Cut(uri, "://")
	if !found {
		return URL{}, errors.New("invalid URL: missing scheme")
	}
	if scheme == "" {
		return URL{}, errors.New("invalid URL: empty scheme")
	}
	if !strings.Contains(urn, "/") {
		urn += "/"
	}
	hostname, path, _ := strings.Cut(urn, "/")
	if hostname == "" {
		return URL{}, errors.New("invalid URL: empty hostname")
	}
	u := URL{
		Scheme: scheme,
		Host:   hostname,
		Path:   "/" + path,
	}
	return u, nil
}
