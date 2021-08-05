package rcmd

import (
	"fmt"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestRVersion(t *testing.T) {
	var tests = []struct {
		in             RVersion
		expectedString string
		expectedFull   string
	}{
		{
			RVersion{3, 5, 2},
			"3.5",
			"3.5.2",
		},
		{
			RVersion{2, 1, 4},
			"2.1",
			"2.1.4",
		},
	}
	for i, test := range tests {
		t.Run(test.expectedFull, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			actual := test.in.ToString()
			t.A.Equal(test.expectedString, actual, fmt.Sprintf("test num: %v", i+1))
			actual = test.in.ToFullString()
			t.A.Equal(test.expectedFull, actual, fmt.Sprintf("test num: %v", i+1))
		})
	}
}
