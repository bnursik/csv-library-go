package main

import (
	"fmt"
	"io"
	"os"
)

func Parser(r io.Reader) {
	buff := make([]byte, 1)
	var line string
	for {
		_, err := r.Read(buff)
		char := string(buff[0])

		if char == "\n" {
			fmt.Println(line)
			break
		}

		line += char

		if err == io.EOF {
			break
		}
	}
}

func main() {
	file, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	Parser(file)
}
