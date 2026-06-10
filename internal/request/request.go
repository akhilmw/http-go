package request

import (
	"errors"
	"io"
	"strings"

	"github.com/akhilmw/http-go/internal/headers"
)

const bufferSize = 8

type parserState int

const (
	initialized parserState = iota
	parsingHeaders
	done
)

type Request struct {
	RequestLine RequestLine
	Headers     headers.Headers
	state       parserState
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

func (r *Request) parse(data []byte) (int, error) {
	totalBytesParsed := 0

	for r.state != done {
		n, err := r.parseSingle(data[totalBytesParsed:])
		if err != nil {
			return 0, err
		}

		if n == 0 {
			break
		}

		totalBytesParsed += n
	}

	return totalBytesParsed, nil
}

func (r *Request) parseSingle(data []byte) (int, error) {
	switch r.state {
	case initialized:
		requestLine, consumed, err := parseRequestLine(string(data))
		if err != nil {
			return 0, err
		}

		if consumed == 0 {
			return 0, nil
		}

		r.RequestLine = *requestLine
		r.state = parsingHeaders

		return consumed, nil

	case parsingHeaders:
		n, headersDone, err := r.Headers.Parse(data)
		if err != nil {
			return 0, err
		}

		if n == 0 {
			return 0, nil
		}

		if headersDone {
			r.state = done
		}

		return n, nil

	case done:
		return 0, errors.New("trying to parse in done state")

	default:
		return 0, errors.New("unknown state")
	}
}

func parseRequestLine(str string) (*RequestLine, int, error) {
	// GET /coffee HTTP/1.1\r\nHost: localhost:42069\r\nUser-Agent: curl/7.81.0\r\nAccept: */*\r\n\r\n

	idx := strings.Index(str, "\r\n")
	if idx == -1 {
		return nil, 0, nil
	}

	requestLine := str[:idx]

	parts := strings.Split(requestLine, " ")
	if len(parts) != 3 {
		return nil, 0, errors.New("invalid request line")
	}

	method := parts[0]
	target := parts[1]
	version := parts[2]

	// check first is a method
	if !isValidMethod(method) {
		return nil, 0, errors.New("method not valid")
	}

	if version != "HTTP/1.1" {
		err := errors.New("only supports HTTP 1.1")
		return nil, 0, err
	}

	bytesConsumed := idx + len("\r\n")

	return &RequestLine{
		Method:        method,
		RequestTarget: target,
		HttpVersion:   "1.1",
	}, bytesConsumed, nil

}

func RequestFromReader(reader io.Reader) (*Request, error) {

	buf := make([]byte, bufferSize)
	readToIndex := 0

	request := &Request{
		state:   initialized,
		Headers: headers.NewHeaders(),
	}

	for request.state != done {

		if readToIndex == len(buf) {
			newBuf := make([]byte, len(buf)*2)
			copy(newBuf, buf)
			buf = newBuf
		}

		n, err := reader.Read(buf[readToIndex:])
		if err != nil {
			if err == io.EOF {
				if request.state != done {
					return nil, errors.New("incomplete request")
				}
				break
			}
			return nil, err
		}

		readToIndex += n
		consumed, err := request.parse(buf[:readToIndex])
		if err != nil {
			return nil, err
		}

		if consumed > 0 {
			copy(buf, buf[consumed:readToIndex])
			readToIndex -= consumed
		}

	}

	return request, nil

}
