package client

import (
	"errors"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReader_Read(t *testing.T) {
	tt := []struct {
		input  string
		result interface{}
		err    string
	}{
		{
			input:  "+OK\r\n",
			result: "OK",
		},
		{
			input:  ":1000\r\n",
			result: int64(1000),
		},
		{
			input:  "$6\r\nfoobar\r\n",
			result: "foobar",
		},
		{
			input:  "$0\r\n\r\n",
			result: "",
		},
		{
			input:  "$-1\r\n",
			result: nil,
		},
		{
			input: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			result: []interface{}{
				"foo",
				"bar",
			},
		},
		{
			input:  "*0\r\n",
			result: []interface{}{},
		},
		{
			input:  "*-1\r\n",
			result: nil,
		},
		{
			input: "*3\r\n:1\r\n:2\r\n:3\r\n",
			result: []interface{}{
				int64(1),
				int64(2),
				int64(3),
			},
		},
		{
			input: "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Foo\r\n-Bar\r\n",
			result: []interface{}{
				[]interface{}{
					int64(1),
					int64(2),
					int64(3),
				},
				[]interface{}{
					"Foo",
					errors.New("Bar"),
				},
			},
		},
		{
			input: "+\r",
			err:   "unexpected end of stream",
		},
		{
			input: "+BROKEN\r",
			err:   "unexpected end of stream",
		},
	}

	for _, test := range tt {
		t.Run(test.input, func(t *testing.T) {
			r := NewReader(strings.NewReader(test.input))
			result, err := r.Read()
			if test.err != "" {
				assert.EqualError(t, err, test.err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, test.result, result.Content())
			}
		})
	}

}
