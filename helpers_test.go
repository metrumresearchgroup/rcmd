package rcmd

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryName(t *testing.T) {
	tests := []struct {
		os       string
		platform string
		expected string
	}{
		{
			os:       "linux",
			platform: "x86_64-pc-linux-gnu",
			expected: "pkg_version_R_x86_64-pc-linux-gnu.tar.gz",
		},
		{
			os:       "darwin",
			expected: "pkg_version.tgz",
		},
		{
			os:       "windows",
			expected: "pkg_version.zip",
		},
	}
	for _, tt := range tests {
		name := binaryNameOs(tt.os, "pkg", "version", tt.platform)
		assert.Equal(t, tt.expected, name, fmt.Sprintf("Not equal: %s", tt.os))
	}
}

func TestBinaryExt(t *testing.T) {
	tests := []struct {
		os       string
		platform string
		path     string
		expected string
	}{
		{
			os:       "linux",
			platform: "PLATFORM",
			path:     "/var/tmp/package_1.2.3.tar.gz",
			expected: "package_1.2.3_R_PLATFORM.tar.gz",
		},
		{
			os:       "darwin",
			platform: "PLATFORM",
			path:     "/var/tmp/package_1.2.3.tar.gz",
			expected: "package_1.2.3.tgz",
		},
		{
			os:       "windows",
			platform: "PLATFORM",
			path:     "/var/tmp/package_1.2.3.tar.gz",
			expected: "package_1.2.3.zip",
		},
	}
	for _, tt := range tests {
		name := binaryExtOs(tt.os, tt.path, tt.platform)
		assert.Equal(t, tt.expected, name, fmt.Sprintf("Not equal: %s", tt.os))
	}
}
