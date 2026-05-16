package main

import (
	"errors"
	"flag"
	"fmt"
	"gobrowser/internal/filereader"
	"gobrowser/internal/httpclient"
	"gobrowser/internal/renderer"
	"gobrowser/internal/urlparser"
	"log"
)

func main() {
	uri, err := parseUriFromArgs()
	if err != nil {
		log.Fatal(err)
	}

	url, err := urlparser.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	var page string
	switch u := url.(type) {
	case *urlparser.WebURL:
		page, err = httpclient.Request(*u)
	case *urlparser.FileURL:
		page, err = filereader.Read(*u)
	default:
		log.Fatalf("unsupported URL type: %T", u)
	}
	if err != nil {
		log.Fatal(err)
	}

	text, err := renderer.Render(page)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(text)
}

func parseUriFromArgs() (string, error) {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return "", errors.New("missing URL")
	}
	return args[0], nil
}
