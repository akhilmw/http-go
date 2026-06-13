package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/akhilmw/http-go/internal/headers"
	"github.com/akhilmw/http-go/internal/request"
	"github.com/akhilmw/http-go/internal/response"
	"github.com/akhilmw/http-go/internal/server"
)

const port = 42069

func handler(w *response.Writer, req *request.Request) {

	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin/") {
		path := strings.TrimPrefix(req.RequestLine.RequestTarget, "/httpbin")
		url := "https://httpbin.org" + path

		resp, err := http.Get(url)
		if err != nil {
			body := []byte("failed to fetch httpbin\n")

			w.WriteStatusLine(response.StatusInternalServerError)

			responseHeaders := response.GetDefaultHeaders(len(body))
			responseHeaders.Set("Content-Type", "text/plain")

			w.WriteHeaders(responseHeaders)
			w.WriteBody(body)
			return
		}
		defer resp.Body.Close()

		w.WriteStatusLine(response.StatusOK)

		responseHeaders := response.GetDefaultHeaders(0)
		responseHeaders.Del("Content-Length")
		responseHeaders.Set("Transfer-Encoding", "chunked")
		responseHeaders.Set("Trailer", "X-Content-SHA256, X-Content-Length")
		responseHeaders.Set("Content-Type", "text/plain")

		w.WriteHeaders(responseHeaders)

		buf := make([]byte, 1024)
		var fullBody []byte
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				chunk := buf[:n]

				fullBody = append(fullBody, chunk...)

				if _, writeErr := w.WriteChunkedBody(chunk); writeErr != nil {
					return
				}
			}

			if err == io.EOF {
				break
			}

			if err != nil {
				return
			}
		}
		hash := sha256.Sum256(fullBody)

		trailers := headers.NewHeaders()
		trailers.Set("X-Content-SHA256", fmt.Sprintf("%x", hash))
		trailers.Set("X-Content-Length", fmt.Sprint(len(fullBody)))

		w.WriteChunkedBodyDone()
		w.WriteTrailers(trailers)
		return
	}

	var statusCode response.StatusCode
	var body []byte

	switch req.RequestLine.RequestTarget {
	case "/yourproblem":
		statusCode = response.StatusBadRequest
		body = []byte(`<html>
<head><title>400 Bad Request</title></head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>`)

	case "/myproblem":
		statusCode = response.StatusInternalServerError
		body = []byte(`<html>
<head><title>500 Internal Server Error</title></head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>`)

	default:
		statusCode = response.StatusOK
		body = []byte(`<html>
<head><title>200 OK</title></head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>`)
	}

	w.WriteStatusLine(statusCode)

	headers := response.GetDefaultHeaders(len(body))
	headers.Set("Content-Type", "text/html")

	w.WriteHeaders(headers)
	w.WriteBody(body)
}

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
	defer server.Close()
	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan
	log.Println("Server gracefully stopped")
}
