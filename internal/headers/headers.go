package headers

import (
	"bytes"
	"errors"
	"strings"
)

type Headers map[string]string

func NewHeaders() Headers {
	return make(Headers)
}

func (h Headers) Parse(data []byte) (n int, done bool, err error) {
	// ("Host: localhost:42069\r\n\r\n")

	idx := bytes.Index(data, []byte("\r\n"))
	if idx == -1 {
		return 0, false, nil
	}

	if idx == 0 {
		return 2, true, nil
	}

	line := string(data[:idx])

	parts := strings.SplitN(line, ":", 2)
	if len(parts) != 2 {
		return 0, false, errors.New("invalid header")
	}

	rawKey := parts[0]
	rawValue := parts[1]

	// No spaces allowed before the colon or before the field name
	if rawKey != strings.TrimSpace(rawKey) {
		return 0, false, errors.New("invalid header")
	}

	value := strings.TrimSpace(rawValue)

	h[rawKey] = value

	return idx + 2, false, nil

}
