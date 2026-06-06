package request

import (
	"errors"
	"io"
	"strings"
)

type Request struct {
	RequestLine RequestLine
}

type RequestLine struct {
	HttpVersion   string
	RequestTarget string
	Method        string
}

func isValidMethod(method string) bool {
	if method == "" {
		return false
	}

	for _, ch := range method {
		if ch < 'A' || ch > 'Z' {
			return false
		}
	}

	return true
}

func parseRequestLine(str string) (*RequestLine, error) {
	// GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n

	lines := strings.Split(str, "\r\n")
	requestLine := lines[0]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, errors.New("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// check first is a method
	if !isValidMethod(method) {
		return nil, errors.New("method not valid")
	}

	if version != "HTTP/1.1" {
		err := errors.New("Only supports HTTP 1.1")
		return nil, err
	}

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	bytes, err := io.ReadAll(reader)
	if err != nil {
		return nil, err
	}

	str := string(bytes)

	requestLine, err := parseRequestLine(str)
	if err != nil {
		return nil, err
	}

	request := Request{RequestLine: *requestLine}
	return &request, nil

}
