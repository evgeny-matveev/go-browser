package httpclient

import (
	"bufio"
	"fmt"
	"gobrowser/internal/urlparser"
	"net"
	"strings"
	"testing"
)

func TestRequestReturnsSuccessfulResponseBody(t *testing.T) {
	body, requestLines, err := requestWithResponse(t, "HTTP/1.0 200 OK\r\nContent-Type: text/html\r\n\r\nHello, browser!")
	if err != nil {
		t.Fatalf("Request returned error: %v", err)
	}

	if body != "Hello, browser!" {
		t.Fatalf("Request returned body %q, want %q", body, "Hello, browser!")
	}

	wantRequestLines := []string{
		"GET /index.html HTTP/1.0",
		"Host: example.com",
		"",
	}
	if fmt.Sprint(requestLines) != fmt.Sprint(wantRequestLines) {
		t.Fatalf("Request sent lines %q, want %q", requestLines, wantRequestLines)
	}
}

func TestRequestRejectsUnsupportedStatus(t *testing.T) {
	_, _, err := requestWithResponse(t, "HTTP/1.0 404 Not Found\r\n")
	if err == nil {
		t.Fatal("Request expected error for non-200 status, got nil")
	}
}

func TestRequestRejectsUnsupportedResponseEncodings(t *testing.T) {
	tests := []struct {
		name     string
		response string
	}{
		{
			name:     "transfer encoding",
			response: "HTTP/1.0 200 OK\r\nTransfer-Encoding: chunked\r\n\r\n",
		},
		{
			name:     "content encoding",
			response: "HTTP/1.0 200 OK\r\nContent-Encoding: gzip\r\n\r\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, _, err := requestWithResponse(t, tt.response)
			if err == nil {
				t.Fatal("Request expected error for unsupported response encoding, got nil")
			}
		})
	}
}

func requestWithResponse(t *testing.T, response string) (string, []string, error) {
	t.Helper()

	originalDial := dial
	t.Cleanup(func() {
		dial = originalDial
	})

	requestLines := make(chan []string, 1)
	dial = func(network, address string) (net.Conn, error) {
		if network != "tcp" {
			return nil, fmt.Errorf("network = %q, want %q", network, "tcp")
		}
		if address != "example.com:80" {
			return nil, fmt.Errorf("address = %q, want %q", address, "example.com:80")
		}

		client, server := net.Pipe()
		go serveResponse(server, response, requestLines)
		return client, nil
	}

	body, err := Request(urlparser.URL{
		Scheme: "http",
		Host:   "example.com",
		Path:   "/index.html",
	})

	return body, <-requestLines, err
}

func serveResponse(conn net.Conn, response string, requestLines chan<- []string) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	lines := []string{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			requestLines <- lines
			return
		}

		line = strings.TrimSuffix(line, "\r\n")
		lines = append(lines, line)
		if line == "" {
			break
		}
	}

	requestLines <- lines
	_, _ = conn.Write([]byte(response))
}
