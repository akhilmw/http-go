package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseHeaders(t *testing.T) {
	t.Run("Valid single header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, 23, n)
		assert.False(t, done)
	})

	t.Run("Valid single header with extra whitespace", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host:           localhost:42069    \r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, len("Host:           localhost:42069    \r\n"), n)
		assert.False(t, done)
	})

	t.Run("Valid 2 headers with existing headers", func(t *testing.T) {
		headers := NewHeaders()
		headers["host"] = "localhost:42069"

		data := []byte("User-Agent: curl/8.0\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "localhost:42069", headers["host"])
		assert.Equal(t, "curl/8.0", headers["user-agent"])
		assert.Equal(t, len("User-Agent: curl/8.0\r\n"), n)
		assert.False(t, done)
	})

	t.Run("Valid done", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, 2, n)
		assert.True(t, done)
	})

	t.Run("Invalid spacing header", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("       Host: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Invalid space before colon", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Host : localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})

	t.Run("Valid single header with uppercase key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("Content-Length: 10\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "10", headers["content-length"])
		assert.Equal(t, len("Content-Length: 10\r\n"), n)
		assert.False(t, done)
	})

	t.Run("Valid single header with mixed case key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("X-Forwarded-For: 127.0.0.1\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.NoError(t, err)
		assert.Equal(t, "127.0.0.1", headers["x-forwarded-for"])
		assert.Equal(t, len("X-Forwarded-For: 127.0.0.1\r\n"), n)
		assert.False(t, done)
	})

	t.Run("Invalid character in header key", func(t *testing.T) {
		headers := NewHeaders()
		data := []byte("H©st: localhost:42069\r\n\r\n")

		n, done, err := headers.Parse(data)

		require.Error(t, err)
		assert.Equal(t, 0, n)
		assert.False(t, done)
	})
}
