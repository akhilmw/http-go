package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {

	f, err := os.Open("messages.txt")
	if err != nil {
		fmt.Println("Some error : ", err)
	}

	defer f.Close()

	parts := make([]any, 0, 8)

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
				} else {
					actual := currentChunk[:lastLineIdx]
					parts = append(parts, actual)
					// traverse the current state of parts that will hold the current line
					var completeLine string
					for _, p := range parts {
						byteSlice := p.([]byte)
						completeLine += string(byteSlice)
					}

					fmt.Printf("read : %s\n", completeLine)
					parts = parts[:0]

					currentChunk = currentChunk[lastLineIdx+1:]
					if len(currentChunk) == 0 {
						break
					}
				}
			}
		}

		if err != nil {
			if errors.Is(err, io.EOF) {
				// Print "read: end" as required by your test assertion
				fmt.Println("read: end")
				fmt.Println("File reading complete")
				return
			}
			fmt.Println("Read error:", err)
			return
		}

	}

}
