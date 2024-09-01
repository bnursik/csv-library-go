package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type CSVParser interface {
	ReadLine(r io.Reader) (*string, error)
	GetField(n int) (string, error)
}

var (
	ErrQuote      = errors.New("excess or missing \" in quoted-field")
	ErrFieldCount = errors.New("wrong number of fields")
)

type CSV struct {
	line string
}

func (c *CSV) ReadLine(r io.Reader) (*string, error) {
	buff := make([]byte, 1)
	lastchar := ""
	inQuotes := false
	c.line = ""

	for {
		_, err := r.Read(buff)
		char := string(buff[0])
		if char == "\n" || char == "\r" || char == "\r\n" || err == io.EOF {
			if inQuotes {
				return nil, ErrQuote
			}
			return &c.line, err
		}

		if char == `"` {
			if inQuotes && lastchar == `"` { // single quote inside the field
				inQuotes = !inQuotes
				continue
			}

			inQuotes = !inQuotes
		}

		c.line += char
		lastchar = char
	}
}

func (c *CSV) GetField(n int) (string, error) {
	if n < 0 || n > (1+strings.Count(c.line, ",")) {
		return "", ErrFieldCount
	}

	commaN := 0
	lastComma := -1

	for i := 0; i < len(c.line); i++ {
		if string(c.line[i]) == "," {
			commaN += 1
			if commaN == n {
				s := string(c.line[lastComma+1 : i])
				if len(s) > 0 && s[0] == '"' {
					s = s[1:]
				}
				if len(s) > 0 && s[len(s)-1] == '"' {
					s = s[:len(s)-1]
				}
				return s, nil
			}
			lastComma = i
		}
	}

	s := string(c.line[lastComma+1:])
	if len(s) > 0 && s[0] == '"' {
		s = s[1:]
	}
	if len(s) > 0 && s[len(s)-1] == '"' {
		s = s[:len(s)-1]
	}
	return s, nil
}

func NewCsv() *CSV {
	return new(CSV)
}

func main() {
	file, err := os.Open("sample.csv")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	CSVParser := NewCsv()

	for {
		_, err := CSVParser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				field, _ := CSVParser.GetField(3)
				fmt.Println(field)
				break
			}
			fmt.Println("Error reading line:", err)
			return
		}

		field, _ := CSVParser.GetField(3)
		fmt.Println(field)
	}
}
