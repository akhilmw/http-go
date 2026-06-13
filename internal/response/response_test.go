package response

import (
	"bytes"
	"testing"

	"github.com/akhilmw/http-go/internal/headers"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriteChunkedBody(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	n, err := w.WriteChunkedBody([]byte("Hello"))

	require.NoError(t, err)
	assert.Equal(t, 5, n)
	assert.Equal(t, "5\r\nHello\r\n", buf.String())
}

func TestWriteChunkedBodyHexSize(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	body := []byte("I could go for a cup of coffee") // 30 bytes = 1E hex
	n, err := w.WriteChunkedBody(body)

	require.NoError(t, err)
	assert.Equal(t, len(body), n)
	assert.Equal(t, "1E\r\nI could go for a cup of coffee\r\n", buf.String())
}

func TestWriteChunkedBodyDone(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	n, err := w.WriteChunkedBodyDone()

	require.NoError(t, err)
	assert.Equal(t, len("0\r\n"), n)
	assert.Equal(t, "0\r\n", buf.String())
}

func TestWriteTrailers(t *testing.T) {
	var buf bytes.Buffer
	w := NewWriter(&buf)

	h := headers.NewHeaders()
	h.Set("X-Content-SHA256", "abc123")
	h.Set("X-Content-Length", "42")

	err := w.WriteTrailers(h)

	require.NoError(t, err)

	output := buf.String()
	assert.Contains(t, output, "x-content-sha256: abc123\r\n")
	assert.Contains(t, output, "x-content-length: 42\r\n")
	assert.True(t, bytes.HasSuffix(buf.Bytes(), []byte("\r\n")))
}
