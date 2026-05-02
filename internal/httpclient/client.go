package httpclient

import (
	"bufio"
	"errors"
	"fmt"
	"gobrowser/internal/urlparser"
	"io"
	"net"
	"strconv"
	"strings"
)

type Status struct {
	Version     string
	Code        int
	Explanation string
}

type Headers map[string]string

func Request(url urlparser.URL) (string, error) {
	conn, err := net.Dial("tcp", url.Host+":80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	request := fmt.Sprintf("GET %s HTTP/1.0\r\n", url.Path)
	request += fmt.Sprintf("Host: %s\r\n", url.Host)
	request += "\r\n"

	_, err = conn.Write([]byte(request))
	if err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)

	statusLine, err := readLine(reader)
	if err != nil {
		return "", err
	}
	status, err := parseStatus(statusLine)
	if err != nil {
		return "", err
	}
	// TODO: handle statuses
	if status.Code != 200 {
		return "", fmt.Errorf("unsupported status: %d %s", status.Code, status.Explanation)
	}

	headers, err := parseHeaders(reader)
	if err != nil {
		return "", err
	}
	if len(headers) > 0 {
		if _, ok := headers["transfer-encoding"]; ok {
			return "", errors.New("cannot handle response: transfer-encoding is not supported")
		}
		if _, ok := headers["content-encoding"]; ok {
			return "", errors.New("cannot handle response: content-encoding is not supported")
		}
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	if !strings.HasSuffix(line, "\r\n") {
		return "", errors.New("invalid response: expected CRLF")
	}
	line = strings.TrimRight(line, "\r\n")
	return line, nil
}

func parseStatus(line string) (Status, error) {
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return Status{},
			fmt.Errorf("invalid response status: %v", parts)
	}
	code, err := strconv.Atoi(parts[1])
	if err != nil {
		return Status{}, err
	}
	return Status{
		Version:     parts[0],
		Code:        code,
		Explanation: parts[2],
	}, nil
}

func parseHeaders(reader *bufio.Reader) (Headers, error) {
	headers := Headers{}
	for {
		headerLine, err := readLine(reader)
		if err != nil {
			return Headers{}, err
		}
		if headerLine == "" {
			break
		}
		name, value, found := strings.Cut(headerLine, ":")
		if !found {
			return Headers{}, fmt.Errorf("invalid response header: expected a colon, got: %q", headerLine)
		}
		name = strings.ToLower(strings.TrimSpace(name))
		headers[name] = strings.TrimSpace(value)
	}
	return headers, nil
}
