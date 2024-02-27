package main

import (
	"fmt"
	"io"
	"strconv"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{
		writer: w,
	}
}

// write writes RESP formatted content
func (w *Writer) write(messageType byte, content ...[]byte) error {
	if _, err := w.writer.Write([]byte{messageType}); err != nil {
		return fmt.Errorf("failed to write message type: %v", messageType)
	}

	for _, b := range content {
		if _, err := w.writer.Write(b); err != nil {
			return fmt.Errorf("failed to write bytes")
		}
	}

	return nil
}

// WriteNil writes a nil bulk string
func (w *Writer) WriteNil() error {
	return w.write(BulkStringType, []byte("-1"), Separator)
}

// WriteBulkString writes a bulk string
func (w *Writer) WriteBulkString(value []byte) error {
	return w.write(
		BulkStringType,
		[]byte(strconv.FormatInt(int64(len(value)), 10)),
		Separator,
		value,
		Separator,
	)
}

// WriteInt64 writes a int64
func (w *Writer) WriteInt64(v int64) error {
	return w.write(
		IntegerType,
		[]byte(strconv.FormatInt(v, 10)),
		Separator)
}

// WriteArray writes an array.
// Any unsupported value types of the array will cause an error.
func (w *Writer) WriteArray(values []interface{}) error {
	if values == nil {
		return w.write(ArrayType, []byte("-1"), Separator)
	}

	if err := w.write(
		ArrayType,
		[]byte(strconv.FormatInt(int64(len(values)), 10)),
		Separator,
	); err != nil {
		return err
	}

	for _, v := range values {
		switch t := v.(type) {
		case int:
		case int8:
		case int16:
		case int32:
		case int64:
			if err := w.WriteInt64(int64(t)); err != nil {
				return err
			}
		case string:
			if err := w.WriteBulkString([]byte(t)); err != nil {
				return err
			}
		case []byte:
			if err := w.WriteBulkString(t); err != nil {
				return err
			}
		case []interface{}:
			if err := w.WriteArray(t); err != nil {
				return err
			}
		case nil:
			if err := w.WriteNil(); err != nil {
				return err
			}
		default:
			return fmt.Errorf("unsupported type: the value [%#v] is not supported by this client, supported types are int8 to int64, strings, []byte, nil, and []interface{} of these same types", v)
		}
	}

	return nil
}
