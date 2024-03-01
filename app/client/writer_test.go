package client

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWriterWriteInt64(t *testing.T) {
	input, output := int64(123), int64(123)
	buffer := &bytes.Buffer{}

	writer := NewWriter(buffer)
	err := writer.WriteInt64(input)
	require.NoError(t, err)

	reader := NewReader(buffer)
	result, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, output, result.Content())
}

func TestWriterWriteBulkString(t *testing.T) {
	input, output := []byte("foobar"), "foobar"
	buffer := &bytes.Buffer{}

	writer := NewWriter(buffer)
	err := writer.WriteBulkString(input)
	require.NoError(t, err)

	reader := NewReader(buffer)
	result, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, output, result.Content())
}

func TestWriterWriteNil(t *testing.T) {
	buffer := &bytes.Buffer{}

	writer := NewWriter(buffer)
	err := writer.WriteNil()
	require.NoError(t, err)

	reader := NewReader(buffer)
	result, err := reader.Read()
	require.NoError(t, err)
	assert.Equal(t, nil, result.Content())
}

func TestWriterWriteArray(t *testing.T) {
	tt := []struct {
		name   string
		input  []interface{}
		output interface{}
		err    string
	}{
		{
			name: "a simple array",
			input: []interface{}{
				"foo",
				"bar",
				nil,
				10,
				int64(65),
			},
			output: []interface{}{
				"foo",
				"bar",
				nil,
				int64(10),
				int64(65),
			},
		},
		{
			name: "an array of arrays",
			input: []interface{}{
				[]interface{}{
					"nope",
					"yup",
				},
				10,
				[]interface{}{
					"honey",
					"1000",
				},
			},
			output: []interface{}{
				[]interface{}{
					"nope",
					"yup",
				},
				int64(10),
				[]interface{}{
					"honey",
					"1000",
				},
			},
		},
		{
			name:   "a nil array",
			input:  nil,
			output: nil,
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			buffer := &bytes.Buffer{}

			writer := NewWriter(buffer)
			err := writer.WriteArray(test.input)
			if test.err != "" {
				assert.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)

				reader := NewReader(buffer)
				result, err := reader.Read()
				require.NoError(t, err)
				assert.Equal(t, test.output, result.Content())
			}
		})
	}
}
