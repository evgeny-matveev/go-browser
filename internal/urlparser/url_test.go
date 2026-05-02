package urlparser

import "testing"

func TestParseInvalidURL(t *testing.T) {
	tests := []string{
		"example.com",
		"http://",
		"://example.com",
		"",
	}
	for _, raw := range tests {
		t.Run(raw, func(t *testing.T) {
			_, err := Parse(raw)

			if err == nil {
				t.Fatalf("Parse(%q) expected error, got nil", raw)
			}
		})
	}
}

func TestParse(t *testing.T) {
	tests := []struct {
		name string
		raw  string
		want URL
	}{
		{
			name: "simple http URL",
			raw:  "http://example.com",
			want: URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/",
			},
		},
		{
			name: "url with path",
			raw:  "http://example.com/index.html",
			want: URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/index.html",
			},
		},
		{
			name: "url with nested path",
			raw:  "http://example.com/docs/index.html",
			want: URL{
				Scheme: "http",
				Host:   "example.com",
				Path:   "/docs/index.html",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Parse(tt.raw)
			if err != nil {
				t.Fatalf("Parse(%q) returned error: %v", tt.raw, err)
			}
			if got != tt.want {
				t.Errorf("Parse(%q) = %+v, want %+v", tt.raw, got, tt.want)
			}
		})
	}
}
