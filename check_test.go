package rcmd

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCheckArgs(t *testing.T) {
	tests := []struct {
		input    CheckArgs
		expected []string
	}{
		{
			CheckArgs{},
			[]string{},
		},
		{
			CheckArgs{NoInstall: true, NoVignettes: true},
			[]string{"--no-install", "--no-vignettes"},
		},
		// TODO: consider what happens if dir with space in it
		// would currently return --output=some/dir with space/in
		// which probably would be interpretted wrong and should be string escaped
		{
			CheckArgs{Output: "some/dir"},
			[]string{"--output=some/dir"},
		},
		{
			CheckArgs{Output: "some/dir", Library: "some/other/dir"},
			[]string{"--output=some/dir", "--library=some/other/dir"},
		},
		{
			CheckArgs{AsCran: true},
			[]string{"--as-cran"},
		},

	}
	for _, tt := range tests {
		assert.Equal(t, tt.expected, tt.input.CliArgs())
	}
}
