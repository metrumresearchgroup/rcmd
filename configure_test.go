package rcmd_test

import (
	"io/ioutil"
	"runtime"
	"strings"
	"testing"

	"github.com/metrumresearchgroup/environ"
	"github.com/metrumresearchgroup/wrapt"
	"github.com/stretchr/testify/assert"

	"github.com/metrumresearchgroup/rcmd/v2"
)

type configureArgsTestCase struct {
	context string
	// mocked system environment variables per os.Environ()
	input    []string
	expected []string
}

func TestConfigureArgs1(t *testing.T) {
	defaultRS, err := rcmd.NewRSettings("")
	assert.NoError(t, err)

	// there should always be at least one libpath
	defaultRS.LibPaths = []string{"path/to/install/lib"}
	var tests = []configureArgsTestCase{
		{
			context:  "minimal",
			input:    []string{},
			expected: []string{"R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "non-impactful system env set",
			input:    []string{"MISC_ENV=foo", "MISC2=bar"},
			expected: []string{"MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			context:  "non-impactful system env set with known package",
			input:    []string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=true"},
			expected: []string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "impactful system env set on separate package",
			input:    []string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
			expected: []string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "impactful system env set with known package",
			input:    []string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=true"},
			expected: []string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_SITE env set",
			input:    []string{"R_LIBS_SITE=original/path", "MISC2=bar"},
			expected: []string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_SITE env set with known package",
			input:    []string{"R_LIBS_SITE=original/path", "MISC2=bar", "DPLYR_ENV=true"},
			expected: []string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_USER env set",
			input:    []string{"R_LIBS_USER=original/path", "MISC2=bar"},
			expected: []string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_USER env set with known package",
			input:    []string{"R_LIBS_USER=original/path", "MISC2=bar", "DPLYR_ENV=true"},
			expected: []string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_SITE and R_LIBS_USER env set",
			input:    []string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
			expected: []string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			context:  "R_LIBS_SITE and R_LIBS_USER env set",
			input:    []string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar", "DPLYR_ENV=true"},
			expected: []string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
	}
	for _, test := range tests {
		t.Run(test.context, func(tt *testing.T) {
			t := wrapt.WrapT(tt)

			actual, err := rcmd.ConfigureEnv(test.input, defaultRS)
			t.A.NoError(err)

			// Make sure that all environment variables are present
			// Also make sure that R_LIBS_USER is set.
			checkEnvVarsValid(t, test.expected, actual)

			// assert.Equal(tt.expected, actual, fmt.Sprintf("%s, test num: %v", tt.context, i+1))
		})
	}
}

func Test_configureEnv(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := rcmd.NewRSettings("R")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cleanEnv := environ.FromOS()
	_, err = cleanEnv.Keep("PATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want, err := rcmd.ConfigureEnv(cleanEnv.AsSlice(), rs)
	t.A.NoError(err)
	got, err := rcmd.ConfigureEnv(cleanEnv.AsSlice(), rs)
	t.A.NoError(err)

	wantenv := environ.New(want)
	gotenv := environ.New(got)

	t.Run("Env", func(t *wrapt.T) {
		for _, key := range wantenv.Keys() {
			if wantenv.Get(key) == "" {
				continue
			}
			t.Run(key, func(t *wrapt.T) {
				if key != "R_LIBS_USER" {
					t.Run("match", func(t *wrapt.T) {
						if wantenv.Get(key) != gotenv.Get(key) {
							t.Errorf("wantenv.Get(`%s`): %v, gotenv.Get(`%s`): %v", key, wantenv.Get(key), key, gotenv.Get(key))
						}
					})
				} else {
					t.Run("should not match", func(t *wrapt.T) {
						if wantenv.Get("R_LIBS_USER") == gotenv.Get("R_LIBS_USER") {
							t.Errorf("")
						}
					})
				}
			})
		}
	})
}

// Utility functions.
func checkEnvVarsValid(t *wrapt.T, expected []string, actualResults []string) {
	t.Helper()

	rLibsUserFound := false
	for _, envVar := range actualResults {
		t.Run(envVar, func(t *wrapt.T) {
			if strings.HasPrefix(envVar, "R_LIBS_USER") {
				rLibsUserFound = true
				tmpDir := strings.Split(envVar, "=")[1]
				checkIsTempDir(t, tmpDir)
				assert.DirExists(t, tmpDir)
				dirEntries, err := ioutil.ReadDir(tmpDir)
				assert.Nil(t, err)
				assert.Empty(t, dirEntries, "failure: R_LIBS_USER was not set to an EMPTY temp directory")
			} else {
				assert.Contains(t, expected, envVar, "excess environment vars found")
				// assert.Equal(testCase.expected[index], envVar) // We are no longer claiming that order matters.
			}
		})
	}
	assert.True(t, rLibsUserFound, "R_LIBS_USER was not set -- we expect it to always be set")
	// Make sure we're not missing any expected vars. A little redundant, but the easiest way to do this.
	for _, envVar := range expected {
		if strings.HasPrefix(envVar, "R_LIBS_USER") {
			continue
		} else {
			assert.Contains(t, actualResults, envVar, "missing expected environment var")
		}
	}
}

func checkIsTempDir(t *wrapt.T, tmpDir string) {
	t.Helper()

	switch runtime.GOOS {
	case "darwin":
		assert.True(t, strings.Contains(tmpDir, "var/folders"), "R_LIBS_USER not set to temp directory: Dir found: %s", tmpDir)
	case "linux":
		t.Skip("tmp dir check not implemented for linux")
	case "windows":
		t.Skip("tmp dir check not implemented for linux")
	default:
		t.Skip("tmp dir check not implemented for detected os")
	}
}

// end Utility functions
