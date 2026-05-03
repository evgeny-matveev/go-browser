package main

import (
	"errors"
	"flag"
	"fmt"
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

	page, err := httpclient.Request(url)
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
