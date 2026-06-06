package main

import (
	"bytes"
	"fmt"
	"io"
	"net"

	//"os"
	"strings"
)

// this returns channel of string
func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer f.Close()
		defer close(ch)

		parts := make([][]byte, 0)

		for {

			buffer := make([]byte, 8)
			n, err := f.Read(buffer)

			if n > 0 {
				currentChunk := buffer[:n]
				for {
					lastLineIdx := bytes.IndexByte(currentChunk, '\n')
					if lastLineIdx == -1 {
						parts = append(parts, currentChunk)
						break
					}
					actual := currentChunk[:lastLineIdx]
					parts = append(parts, actual)
					// traverse the current state of parts that will hold the current line
					var completeLine strings.Builder
					for _, part := range parts {
						completeLine.Write(part)
					}

					// fmt.Printf("read : %s\n", completeLine.String())
					ch <- completeLine.String()
					parts = parts[:0]

					currentChunk = currentChunk[lastLineIdx+1:]
					if len(currentChunk) == 0 {
						break
					}

				}
			}

			if err == io.EOF {
				return
			}

			if err != nil {
				return
			}
		}

	}()

	return ch

}

func main() {

	ln, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Println("Error while creating server : ", err)
		return
	}

	defer ln.Close()

	for {
		conn, err := ln.Accept()
		if err != nil {
			fmt.Println("Error while accepting connection : ", err)
			return

		}
		fmt.Println("Connection Accepted ")
		for line := range getLinesChannel(conn) {
			fmt.Printf("read: %s\n", line)
		}

		conn.Close()

	}

}
