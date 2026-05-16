package urlparser

import (
	"errors"
	"strings"
)

type ParsedURL interface {
	Scheme() string
}

// Scheme://Host:Port/Path
type WebURL struct {
	SchemeValue string
	Host        string
	Port        string
	Path        string
}

func (u WebURL) Scheme() string {
	return u.SchemeValue
}

// Scheme:///Path
type FileURL struct {
	SchemeValue string
	Path        string
}

func (u FileURL) Scheme() string {
	return u.SchemeValue
}

type DataURL struct {
	SchemeValue string
	MediaType   string
	Data        string
}

func (u DataURL) Scheme() string {
	return u.SchemeValue
}

// func Parse[URLType URL](uri string) (URLType, error)
func Parse(raw string) (ParsedURL, error) {
	switch {
	case strings.HasPrefix(raw, "http://") ||
		strings.HasPrefix(raw, "https://"):
		return parseWebScheme(raw)
	case strings.HasPrefix(raw, "file://"):
		return parseFileScheme(raw)
	case strings.HasPrefix(raw, "data:"):
		return parseDataScheme(raw)
	case strings.HasPrefix(raw, "/"):
		return &FileURL{SchemeValue: "file", Path: raw}, nil
	default:
		return nil, errors.New("unsupported URL scheme")
	}
}

func parseWebScheme(raw string) (*WebURL, error) {
	scheme, urn, _ := strings.Cut(raw, "://")
	if !strings.Contains(urn, "/") {
		urn += "/"
	}
	hostname, path, _ := strings.Cut(urn, "/")
	if hostname == "" {
		return nil, errors.New("invalid URL: empty hostname")
	}
	hostname, port, _ := strings.Cut(hostname, ":")
	u := WebURL{
		SchemeValue: scheme,
		Host:        hostname,
		Port:        port,
		Path:        "/" + path,
	}
	return &u, nil
}

func parseFileScheme(raw string) (*FileURL, error) {
	scheme, path, _ := strings.Cut(raw, "://")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u := FileURL{
		SchemeValue: scheme,
		Path:        path,
	}
	return &u, nil
}

func parseDataScheme(raw string) (*DataURL, error) {
	scheme, opaque, _ := strings.Cut(raw, "data:")
	mediaType, data, _ := strings.Cut(opaque, ",")
	if mediaType != "text/html" {
		return nil, errors.New("usupported data type")
	}
	u := DataURL{
		SchemeValue: scheme,
		MediaType:   mediaType,
		Data:        data,
	}
	return &u, nil
}
