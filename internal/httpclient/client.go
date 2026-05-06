package httpclient

import (
	"bufio"
	"crypto/tls"
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

var dial = net.Dial

func Request(url urlparser.URL) (string, error) {
	conn, err := connect(url.Scheme, url.Host, url.Port)
	if err != nil {
		return "", err
	}
	defer conn.Close()

	request := prepareGetRequest(url.Path, url.Host)

	if _, err := conn.Write(request); err != nil {
		return "", err
	}

	reader := bufio.NewReader(conn)

	if err := handleStatus(reader); err != nil {
		return "", err
	}

	if err := handleHeaders(reader); err != nil {
		return "", err
	}

	content, err := io.ReadAll(reader)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

func connect(scheme, host, port string) (net.Conn, error) {
	var (
		conn net.Conn
		err  error
	)

	switch scheme {
	case "http":
		if port == "" {
			port = "80"
		}
		conn, err = dial("tcp", host+":"+port)
		if err != nil {
			return conn, err
		}
		return conn, nil
	case "https":
		if port == "" {
			port = "443"
		}
		conn, err = dial("tcp", host+":"+port)
		if err != nil {
			return conn, err
		}
		config := &tls.Config{ServerName: host}
		tlsConn := tls.Client(conn, config)
		if err := tlsConn.Handshake(); err != nil {
			return conn, err
		}
		return tlsConn, nil
	default:
		return conn, fmt.Errorf("unsupported scheme: %s", scheme)
	}
}

func prepareGetRequest(path string, host string) []byte {
	request := fmt.Sprintf("GET %s HTTP/1.1\r\n", path)
	request += fmt.Sprintf("Host: %s\r\n", host)
	request += fmt.Sprintf("Connection: close\r\n")
	request += fmt.Sprintf("User-Agent: Mosaic\r\n")
	request += "\r\n"
	return []byte(request)
}

func handleStatus(reader *bufio.Reader) error {
	statusLine, err := readLine(reader)
	if err != nil {
		return err
	}

	status, err := parseStatus(statusLine)
	if err != nil {
		return err
	}

	// TODO: handle statuses
	if status.Code != 200 {
		return fmt.Errorf("unsupported status: %d %s", status.Code, status.Explanation)
	}

	return nil
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

func handleHeaders(reader *bufio.Reader) error {
	headers, err := parseHeaders(reader)
	if err != nil {
		return err
	}

	if len(headers) > 0 {
		if _, ok := headers["transfer-encoding"]; ok {
			return errors.New("cannot handle response: transfer-encoding is not supported")
		}

		if _, ok := headers["content-encoding"]; ok {
			return errors.New("cannot handle response: content-encoding is not supported")
		}
	}

	return nil
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
