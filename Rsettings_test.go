package rcmd

import (
	"fmt"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLibPathsEnv(t *testing.T) {
	var libPathTests = []struct {
		in       RSettings
		expected string
	}{
		{
			RSettings{
				LibPaths: []string{
					// TODO: check if paths need to be checked for trailing /
					"path/to/folder1/",
					"path/to/folder2/",
				},
			},
			"R_LIBS_SITE=path/to/folder1/:path/to/folder2/",
		},
		{
			RSettings{
				LibPaths: []string{},
			},
			"",
		},
	}
	for _, tt := range libPathTests {
		ok, actual := tt.in.LibPathsEnv()
		if actual != "" && !ok {
			t.Errorf("LibPaths present, should be ok")
		}
		if actual != tt.expected {
			t.Errorf("GOT: %s, WANT: %s", actual, tt.expected)
		}
	}
}

func TestParseVersionData(t *testing.T) {
	var rVersionTests = []struct {
		data     []byte
		version  RVersion
		platform string
		message  string
	}{
		{
			data: []byte(`R version 3.6.0 (2019-04-26) -- "Planting of a Tree"
			Copyright (C) 2019 The R Foundation for Statistical Computing
			Platform: x86_64-apple-darwin15.6.0 (64-bit)
			
			R is free software and comes with ABSOLUTELY NO WARRANTY.
			You are welcome to redistribute it under the terms of the
			GNU General Public License versions 2 or 3.
			For more information about these matters see
			https://www.gnu.org/licenses/.

`),
			version: RVersion{
				Major: 3,
				Minor: 6,
				Patch: 0,
			},
			platform: "x86_64-apple-darwin15.6.0",
			message:  "darwin test",
		},
		{
			data: []byte(`R version 3.5.2 (2018-12-20) -- "Eggshell Igloo"
Copyright (C) 2018 The R Foundation for Statistical Computing
Platform: x86_64-w64-mingw32/x64 (64-bit)
			
R is free software and comes with ABSOLUTELY NO WARRANTY.
You are welcome to redistribute it under the terms of the
GNU General Public License versions 2 or 3.
For more information about these matters see
http://www.gnu.org/licenses/.

`),
			version: RVersion{
				Major: 3,
				Minor: 5,
				Patch: 2,
			},
			platform: "x86_64-w64-mingw32/x64",
			message:  "windows test",
		},
		{
			data: []byte(`
			R version 1.2.3 (2018-12-20) -- "name for Ubuntu"            
			Platform: x86_64-pc-linux-gnu (64-bit)
			`),
			version: RVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			platform: "x86_64-pc-linux-gnu",
			message:  "Manually built Ubuntu test",
		},
	}
	for _, tt := range rVersionTests {
		version, platform := parseVersionData(tt.data)
		assert.Equal(t, tt.version, version, fmt.Sprintf("Version not equal: %s", tt.message))
		assert.Equal(t, tt.platform, platform, fmt.Sprintf("Platform not equal: %s", tt.message))
	}
}

func TestRMethod(t *testing.T) {
	var rTests = []struct {
		rpath    string
		platform string
		expected string
		message  string
	}{
		{
			rpath:    "",
			platform: "windows",
			expected: "R.exe",
			message:  "windows - empty Rpath",
		},
		{
			rpath:    `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
			platform: "windows",
			expected: `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
			message:  "windows - full Rpath",
		},
		{
			rpath:    `C:\Program Files\R\R-3.5.2\bin\i386\R`,
			platform: "windows",
			expected: `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
			message:  "windows - full Rpath, without exe extension",
		},
		{
			rpath:    `R.exe`,
			platform: "windows",
			expected: `R.exe`,
			message:  "windows - R with exe extension",
		},
		{
			rpath:    `R`,
			platform: "windows",
			expected: `R.exe`,
			message:  "windows - R without extension",
		},
		{
			rpath:    "",
			platform: "darwin",
			expected: "R",
			message:  "darwin - empty Rpath",
		},
		{
			rpath:    "/usr/local/bin/R",
			platform: "darwin",
			expected: "/usr/local/bin/R",
			message:  "darwin: full Rpath",
		},
		{
			rpath:    "/usr/local/bin/R/",
			platform: "darwin",
			expected: "/usr/local/bin/R",
			message:  "darwin: full Rpath, trailing /",
		},
		// {
		// 	rpath:    "/R",
		// 	platform: "darwin",
		// 	expected: "/R",
		// 	message:  "darwin: full Rpath, root R /",
		// },
		// {
		// 	rpath:    "/R/",
		// 	platform: "darwin",
		// 	expected: "/R",
		// 	message:  "darwin: full Rpath, root R / with trailing /",
		// },
		// TODO: linux tests
	}
	for _, tt := range rTests {
		if tt.platform == runtime.GOOS {
			rs := NewRSettings(tt.rpath)
			r := rs.R(tt.platform, false)
			assert.Equal(t, tt.expected, r, fmt.Sprintf("R not equal to <%s>. %s", tt.expected, tt.message))
		}
	}
}
