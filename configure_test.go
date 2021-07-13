package rcmd

//
import (
	"testing"

	"github.com/metrumresearchgroup/environ"
)

/*
type configureArgsTestCase struct {
	context string
	// mocked system environment variables per os.Environ()
	input    []string
	expected []string
}
*/

/*
func TestConfigureArgs(t *testing.T) {
	defaultRS := NewRSettings("")
	// there should always be at least one libpath
	defaultRS.LibPaths = []string{"path/to/install/lib"}
	var installArgsTests = []configureArgsTestCase{
		{
			"minimal",
			"",
			[]string{},
			[]string{"R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"non-impactful system env set",
			"",
			[]string{"MISC_ENV=foo", "MISC2=bar"},
			[]string{"MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib"},
		},
		{
			"non-impactful system env set with known package",
			"dplyr",
			[]string{"MISC_ENV=foo", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"impactful system env set on separate package",
			"",
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"impactful system env set with known package",
			"dplyr",
			[]string{"MISC_ENV=foo", "MISC2=bar", "DPLYR_ENV=false"},
			[]string{"DPLYR_ENV=true", "MISC_ENV=foo", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_SITE env set",
			"",
			[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_SITE env set with known package",
			"dplyr",
			[]string{"R_LIBS_SITE=original/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_USER env set",
			"",
			[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_USER env set with known package",
			"dplyr",
			[]string{"R_LIBS_USER=original/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_SITE and R_LIBS_USER env set",
			"",
			[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
			[]string{"MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"R_LIBS_SITE and R_LIBS_USER env set",
			"dplyr",
			[]string{"R_LIBS_USER=original/path", "R_LIBS_SITE=original/site/path", "MISC2=bar"},
			[]string{"DPLYR_ENV=true", "MISC2=bar", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
		{
			"System contains sensitive information",
			"",
			[]string{"R_LIBS_USER=original/path", "GITHUB_PAT=should_get_hidden1", "ghe_token=should_get_hidden2", "ghe_PAT=should_get_hidden3", "github_token=should_get_hidden4"},
			[]string{"GITHUB_PAT=**HIDDEN**", "ghe_token=**HIDDEN**", "ghe_PAT=**HIDDEN**", "github_token=**HIDDEN**", "R_LIBS_SITE=path/to/install/lib", "R_LIBS_USER=SHOULD_BE_TMP_DIR"},
		},
	}
	for _, tt := range installArgsTests {
		t.Run(tt.context, func(t *testing.T) {
			actual := configureEnv(tt.sysEnv, defaultRS)

			// Make sure that all environment variables are present
			// Also make sure that R_LIBS_USER is set.
			checkEnvVarsValid(t, tt, actual)

			// assert.Equal(tt.expected, actual, fmt.Sprintf("%s, test num: %v", tt.context, i+1))
		})
	}
}
*/

/*
func TestConfigureArgs(t *testing.T) {
	var installArgsTests = []configureArgsTestCase{
		{
			"variety of spellings",
			[]string{
				"R_LIBS_USER=some/path",
				"GITHUB_PAT=should_get_hidden1",
				"ghe_token=should_get_hidden2",
				"ghe_PAT=should_get_hidden3",
				"github_token=should_get_hidden4",
				"AWS_ACCESS_KEY_ID=should_get_hidden5",
				"AWS_SECRET_KEY=should_get_hidden6",
				"ADDL_ARG=could-be-secret",
			},
			[]string{
				"R_LIBS_USER=some/path",
				"GITHUB_PAT=***HIDDEN***",
				"ghe_token=***HIDDEN***",
				"ghe_PAT=***HIDDEN***",
				"github_token=***HIDDEN***",
				"AWS_ACCESS_KEY_ID=***HIDDEN***",
				"AWS_SECRET_KEY=***HIDDEN***",
				"ADDL_ARG=could-be-secret",
			},
		},
	}
	for _, tt := range installArgsTests {
		actual := censorEnvVars(tt.input)
		assert.Equal(t, tt.expected, actual, tt.context)
	}
}
*/
/*
func TestConfigureArgsAddl(t *testing.T) {
	var installArgsTests = []configureArgsTestCase{
		{
			"additional hidden",
			[]string{
				"R_LIBS_USER=some/path",
				"GITHUB_PAT=should_get_hidden1",
				"ADDL_ARG=could-be-secret",
			},
			[]string{
				"R_LIBS_USER=some/path",
				"GITHUB_PAT=***HIDDEN***",
				"ADDL_ARG=***HIDDEN***",
			},
		},
	}
	for _, tt := range installArgsTests {
		actual := censorEnvVars(tt.input, "ADDL_ARG")
		assert.Equal(t, tt.expected, actual, tt.context)
	}
}
*/

func Test_configureEnv(t *testing.T) {
	rs, err := NewRSettings("R")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	cleanEnv := environ.FromOS()
	_, err = cleanEnv.Keep("PATH")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	want := configureEnv(cleanEnv.AsSlice(), rs)
	got := configureEnv(cleanEnv.AsSlice(), rs)

	wantenv := environ.New(want)
	gotenv := environ.New(got)

	t.Run("Env", func(t *testing.T) {
		for _, key := range wantenv.Keys() {
			if wantenv.Get(key) == "" {
				continue
			}
			t.Run(key, func(t *testing.T) {
				if key != "R_LIBS_USER" {
					t.Run("match", func(t *testing.T) {
						if wantenv.Get(key) != gotenv.Get(key) {
							t.Errorf("wantenv.Get(`%s`): %v, gotenv.Get(`%s`): %v", key, wantenv.Get(key), key, gotenv.Get(key))
						}
					})
				} else {
					t.Run("should not match", func(t *testing.T) {

						if wantenv.Get("R_LIBS_USER") == gotenv.Get("R_LIBS_USER") {
							t.Errorf("")
						}
					})
				}
			})
		}
	})
}
