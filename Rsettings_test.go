package rcmd

import (
	"runtime"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func TestLibPathsEnv(t *testing.T) {
	var tests = []struct {
		name     string
		in       RSettings
		expected string
	}{
		{
			name: "path found",
			in: RSettings{
				LibPaths: []string{
					// TODO: check if paths need to be checked for trailing /
					"path/to/folder1/",
					"path/to/folder2/",
				},
			},
			expected: "path/to/folder1/:path/to/folder2/",
		},
		{
			name: "empty paths",
			in: RSettings{
				LibPaths: []string{},
			},
			expected: "",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(tt *testing.T) {
			t := wrapt.WrapT(t)

			actual, ok := test.in.LibPathsEnv()

			if actual != "" && !ok {
				t.Errorf("LibPaths exist but ok is false")
			}

			t.A.Equal(test.expected, actual)
		})
	}
}

func TestParseVersionData(t *testing.T) {
	var tests = []struct {
		data     []byte
		version  *RVersion
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
			version: &RVersion{
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
			version: &RVersion{
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
			version: &RVersion{
				Major: 1,
				Minor: 2,
				Patch: 3,
			},
			platform: "x86_64-pc-linux-gnu",
			message:  "Manually built Ubuntu test",
		},
	}
	for _, test := range tests {
		t.Run(test.message, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			version, platform, err := parseVersionData(test.data)

			t.A.NoError(err)
			t.A.Equal(test.version, version)
			t.A.Equal(test.platform, platform)
		})
	}
}

func TestRMethod(t *testing.T) {
	var tests = []struct {
		name     string
		rpath    string
		platform string
		expected string
	}{
		{
			name:     "windows - empty Rpath",
			rpath:    "",
			platform: "windows",
			expected: "R.exe",
		},
		{
			name:     "windows - full Rpath",
			rpath:    `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
			platform: "windows",
			expected: `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
		},
		{
			name:     "windows - full Rpath, without exe extension",
			rpath:    `C:\Program Files\R\R-3.5.2\bin\i386\R`,
			platform: "windows",
			expected: `C:\Program Files\R\R-3.5.2\bin\i386\R.exe`,
		},
		{
			name:     "windows - R with exe extension",
			rpath:    `R.exe`,
			platform: "windows",
			expected: `R.exe`,
		},
		{
			name:     "windows - R without extension",
			rpath:    `R`,
			platform: "windows",
			expected: `R.exe`,
		},
		{
			name:     "darwin - empty Rpath",
			rpath:    "",
			platform: "darwin",
			expected: "R",
		},
		{
			name:     "darwin: full Rpath",
			rpath:    "/usr/local/bin/R",
			platform: "darwin",
			expected: "/usr/local/bin/R",
		},
		{
			name:     "darwin: full Rpath, trailing /",
			rpath:    "/usr/local/bin/R/",
			platform: "darwin",
			expected: "/usr/local/bin/R",
		},
	}
	for _, test := range tests {
		if test.platform == runtime.GOOS {
			t.Run(test.name, func(tt *testing.T) {
				t := wrapt.WrapT(tt)

				rs, err := NewRSettings(test.rpath)
				t.A.NoError(err)

				r := rs.R(test.platform, false)

				t.A.Equal(test.expected, r)
			})
		}
	}
}
