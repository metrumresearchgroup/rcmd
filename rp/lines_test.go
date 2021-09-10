package rp

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestLineScanning(t *testing.T) {
	var installArgsTests = []struct {
		in       []byte
		fn       func([]byte) []byte
		expected []byte
		context  string
	}{
		{
			in:       []byte("[1] line 1 info"),
			expected: []byte("line 1 info\n"),
			context:  "simplest",
		},
		{
			in: []byte(`
> 2+2
[1] 4
> q("no")
`),
			fn:       OutputOnly,
			expected: []byte("4\n"),
			context:  "simplest",
		},
		{
			in: []byte(`[1] line 1  
[2]	line 2  	`),
			expected: []byte("line 1\nline 2\n"),
			context:  "two lines with whitespace",
		},
		{
			in: []byte(`[1] line 1  
[2]	line 2  	
[3]
`),
			expected: []byte("line 1\nline 2\n"),
			context:  "two lines with trailing new lines",
		},
	}
	for _, test := range installArgsTests {
		t.Run(test.context, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			if test.fn == nil {
				test.fn = ScanLines
			}
			actual := test.fn(test.in)

			t.Run("expected", func(t *wrapt.T) {
				t.A.Equal(test.expected, actual)
			})
		})
	}
}
