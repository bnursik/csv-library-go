package main

import (
	"errors"
	"fmt"
	"io"
	"os"
)

type CSVParser interface {
	ReadLine(r io.Reader) (*string, error)
}

var ErrQuote = errors.New("excess or missing \" in quoted-field") // ErrFieldCount = errors.New("wrong number of fields")

type CSV struct{}

func (c CSV) ReadLine(r io.Reader) (*string, error) {
	buff := make([]byte, 1)
	line := ""

	for {
		_, err := r.Read(buff)
		char := string(buff[0])
		if char == "\n" || char == "\r" || char == "\r\n" || err == io.EOF {
			return &line, err
		}

		line += char
	}
}

func main() {
	file, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	var csvparser CSVParser = CSV{}

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				fmt.Println(*line)
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}
		fmt.Println(*line)
	}
}