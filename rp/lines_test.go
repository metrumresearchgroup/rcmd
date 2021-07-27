package rp

import (
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestLineScanning(t *testing.T) {
	var installArgsTests = []struct {
		in       []byte
		expected []string
		context  string
	}{
		{
			[]byte("[1] line 1 info"),
			[]string{"line 1 info"},
			"simplest",
		},
		{
			[]byte(`[1] line 1  
[2]	line 2  	`),
			[]string{
				"line 1",
				"line 2",
			},
			"two lines with whitespace",
		},
		{
			[]byte(`[1] line 1  
[2]	line 2  	
[3]
`),
			[]string{
				"line 1",
				"line 2",
			},
			"two lines with trailing new lines",
		},
	}
	for _, test := range installArgsTests {
		t.Run(test.context, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			actual, err := ScanLines(test.in)

			t.A.NoError(err)

			t.Run("expected", func(t *wrapt.T) {
				t.A.Equal(test.expected, actual)
			})
		})

	}
}
