package httpclient

import (
	"bufio"
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	validTests := []struct {
		name    string
		raw     string
		want    string
	}{
		{
			name: "reads crlf line",
			raw:  "hello\r\n",
			want: "hello",
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLine(bufio.NewReader(strings.NewReader(tt.raw)))
			if err != nil {
				t.Fatalf("readLine(%q) returned error: %v", tt.raw, err)
			}
			if got != tt.want {
				t.Fatalf("readLine(%q) = %q, want %q", tt.raw, got, tt.want)
			}
		})
	}

	invalidTests := []struct {
		name string
		raw  string
	}{
		{
			name: "rejects non crlf line",
			raw:  "hello\n",
		},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := readLine(bufio.NewReader(strings.NewReader(tt.raw)))
			if err == nil {
				t.Fatalf("readLine(%q) expected error, got %q", tt.raw, got)
			}
		})
	}
}

func TestParseStatus(t *testing.T) {
	validTests := []struct {
		name    string
		raw     string
		want    Status
	}{
		{
			name: "parses valid status",
			raw:  "HTTP/1.0 200 OK",
			want: Status{
				Version:     "HTTP/1.0",
				Code:        200,
				Explanation: "OK",
			},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStatus(tt.raw)
			if err != nil {
				t.Fatalf("parseStatus(%q) returned error: %v", tt.raw, err)
			}
			if got != tt.want {
				t.Fatalf("parseStatus(%q) = %+v, want %+v", tt.raw, got, tt.want)
			}
		})
	}

	invalidTests := []struct {
		name string
		raw  string
	}{
		{
			name: "rejects malformed status",
			raw:  "HTTP/1.0 OK",
		},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseStatus(tt.raw)
			if err == nil {
				t.Fatalf("parseStatus(%q) expected error, got %+v", tt.raw, got)
			}
		})
	}
}

func TestParseHeaders(t *testing.T) {
	validTests := []struct {
		name    string
		raw     string
		want    Headers
	}{
		{
			name: "parses headers until empty line",
			raw:  "Content-Type: text/html\r\nX-Test:  hello \r\n\r\nrest of body",
			want: Headers{
				"content-type": "text/html",
				"x-test":       "hello",
			},
		},
	}

	for _, tt := range validTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaders(bufio.NewReader(strings.NewReader(tt.raw)))
			if err != nil {
				t.Fatalf("parseHeaders(%q) returned error: %v", tt.raw, err)
			}
			if len(got) != len(tt.want) {
				t.Fatalf("parseHeaders(%q) = %+v, want %+v", tt.raw, got, tt.want)
			}
			for key, wantValue := range tt.want {
				if got[key] != wantValue {
					t.Fatalf("parseHeaders(%q)[%q] = %q, want %q", tt.raw, key, got[key], wantValue)
				}
			}
		})
	}

	invalidTests := []struct {
		name string
		raw  string
	}{
		{
			name: "rejects header without colon",
			raw:  "Content-Type text/html\r\n\r\n",
		},
	}

	for _, tt := range invalidTests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseHeaders(bufio.NewReader(strings.NewReader(tt.raw)))
			if err == nil {
				t.Fatalf("parseHeaders(%q) expected error, got %+v", tt.raw, got)
			}
		})
	}
}
