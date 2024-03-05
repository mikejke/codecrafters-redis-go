package client

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	SimpleStringType = '+'
	ErrorType        = '-'
	BulkStringType   = '$'
	ArrayType        = '*'
	IntegerType      = ':'
)

var (
	Separator = []byte("\r\n")
)

type Reader struct {
	scanner *bufio.Scanner
}

func NewReader(r io.Reader) *Reader {
	scanner := bufio.NewScanner(bufio.NewReaderSize(r, 1024))
	scanner.Split(splitter)

	return &Reader{
		scanner: scanner,
	}
}

func (r *Reader) Read() (*Result, error) {
	return readRESP(r.scanner)
}

func splitter(data []byte, atEOF bool) (int, []byte, error) {
	if len(data) < 3 {
		if atEOF {
			return 0, nil, fmt.Errorf("unexpected end of stream")
		}

		return 0, nil, nil
	}

	found := bytes.Index(data, Separator)

	if found == -1 && atEOF {
		return 0, nil, fmt.Errorf("unexpected end of stream")
	}

	if found == -1 {
		return 0, nil, nil
	}

	if data[0] == BulkStringType {
		length, err := strconv.ParseInt(string(data[1:found]), 10, 64)
		if err != nil {
			return 0, nil, err
		}

		if length == -1 {
			return 5, []byte{BulkStringType}, nil
		}

		if length == 0 {
			return 6, []byte{SimpleStringType}, nil
		}

		expectedEnding := found + int(length) + 4
		if len(data) >= expectedEnding {
			start := found + 1
			data[start] = '+'
			return expectedEnding, data[start : expectedEnding-2], nil
		}

		if atEOF {
			return 0, nil, fmt.Errorf("unexpected end of stream")
		}

		return 0, nil, err
	}

	return found + 2, data[:found], nil
}

func readRESP(r *bufio.Scanner) (*Result, error) {
	for r.Scan() {
		line := r.Text()
		switch line[0] {
		case SimpleStringType:
			return &Result{
				content: line[1:],
			}, nil
		case BulkStringType:
			return &Result{
				content: nil,
			}, nil
		case ErrorType:
			// if an error just wrap the error and return it
			return &Result{
				content: errors.New(line[1:]),
			}, nil
		case IntegerType:
			content, err := strconv.ParseInt(line[1:], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse returned integer: %v (value: %v)", err, line)
			}
			return &Result{
				content: content,
			}, nil
		case ArrayType:
			length, err := strconv.ParseInt(line[1:], 10, 64)
			if err != nil {
				return nil, fmt.Errorf("failed to parse array length: %v (value: %v)", err, line)
			}

			if length == -1 {
				return &Result{content: nil}, nil
			}

			contents := make([]interface{}, 0, length)

			for x := int64(0); x < length; x++ {
				result, err := readRESP(r)
				if err != nil {
					return nil, fmt.Errorf("failed to read item %v from array", x)
				}

				contents = append(contents, result.content)
			}

			return &Result{
				content: contents,
			}, nil
		}
	}

	if r.Err() == nil {
		return nil, errors.New("scanner was empty")
	}

	return nil, r.Err()
}
