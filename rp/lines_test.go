package rp_test

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"

	. "github.com/metrumresearchgroup/rcmd/v2/rp"
)

func TestLineScanning(tt *testing.T) {
	var installArgsTests = []struct {
		name     string
		in       []byte
		fn       func([]byte) []byte
		expected []byte
	}{
		{
			in:       []byte("[1] line 1 info"),
			expected: []byte("line 1 info"),
			name:     "simplest",
		},
		{
			in: []byte(`
> 2+2
[1] 4
> q("no")
`),
			fn:       OutputOnly,
			expected: []byte("4\n"),
			name:     "simplest",
		},
		{
			in: []byte(`[1] line 1  
[2]	line 2  	`),
			expected: []byte("line 1\nline 2"),
			name:     "two lines with whitespace",
		},
		{
			in: []byte(`[1] line 1  
[2]	line 2  	
[3]
`),
			expected: []byte("line 1\nline 2\n"),
			name:     "two lines with trailing new lines",
		},
	}
	for _, test := range installArgsTests {
		tt.Run(test.name, func(tt *testing.T) {
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
