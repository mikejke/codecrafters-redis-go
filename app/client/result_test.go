package client

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResultContent(t *testing.T) {
	tt := []struct {
		name   string
		input  interface{}
		result interface{}
	}{
		{
			name:   "a good int64",
			input:  int64(10),
			result: int64(10),
		},
		{
			name:   "a string instead of an int64",
			input:  "some string",
			result: "some string",
		},
	}

	for _, test := range tt {
		t.Run(test.name, func(t *testing.T) {
			r := &Result{
				content: test.input,
			}

			value := r.Content()
			assert.Equal(t, test.result, value)
		})
	}
}
