package main

import (
	"flag"
	"gobrowser/internal/filereader"
	"gobrowser/internal/gui"
	"gobrowser/internal/httpclient"
	"gobrowser/internal/page"
	"gobrowser/internal/urlparser"
	"log"
)

func main() {
	uri := parseUriFromArgs()
	b := gui.NewBrowser()

	if len(uri) == 0 {
		b.Render(page.Lex("Hello, World!"))
		return
	}

	url, err := urlparser.Parse(uri)
	if err != nil {
		log.Fatal(err)
	}

	var body string
	switch u := url.(type) {
	case *urlparser.WebURL:
		body, err = httpclient.Request(*u)
	case *urlparser.FileURL:
		body, err = filereader.Read(*u)
	default:
		log.Fatalf("unsupported URL type: %T", u)
	}
	if err != nil {
		log.Fatal(err)
	}

	tokens := page.Lex(body)

	b.Render(tokens)
}

func parseUriFromArgs() string {
	flag.Parse()
	args := flag.Args()
	if len(args) < 1 {
		return ""
	}
	return args[0]
}
