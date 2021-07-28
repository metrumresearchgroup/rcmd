package rcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/metrumresearchgroup/environ"
	log "github.com/sirupsen/logrus"
)

var censoredVars = map[string]interface{}{
	"GITHUB_TOKEN":      nil,
	"GITHUB_PAT":        nil,
	"GHE_TOKEN":         nil,
	"GHE_PAT":           nil,
	"AWS_ACCESS_KEY_ID": nil,
	"AWS_SECRET_KEY":    nil,
}

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnv(sysEnvVars []string, rs *RSettings) ([]string, error) {
	if rs.AsUser {
		return configureEnvAsUser(sysEnvVars, rs)
	} else {
		return configureEnvAsNotUser(sysEnvVars, rs)
	}
}

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnvAsNotUser(sysEnvVars []string, rs *RSettings) ([]string, error) {
	evs := environ.New(sysEnvVars)
	_, err := evs.Drop("R_LIBS_USER", "R_LIBS_SITE")
	if err != nil {
		return nil, err
	}

	if path := prependPath(evs.Get("PATH"), rs.RPath); path != "" {
		evs.Set("PATH", path)
	}

	tmpdir, err := os.MkdirTemp("", "")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Warn("error making temporary directory while overriding R_LIBS_USER for install.")
	}

	// Force R_LIBS_USER to be an non-empty dir so that we can be sure it won't get overridden by default R paths.
	evs.Set("R_LIBS_USER", tmpdir)

	lp, ok := rs.LibPathsEnv()
	if ok {
		evs.Set("R_LIBS_SITE", lp)
	}

	return evs.AsSlice(), nil
}

func prependPath(path string, rPath string) string {
	rDir := filepath.Dir(rPath)

	if rDir != "" && rDir != "." && !strings.HasPrefix(path, rDir) {
		path = fmt.Sprintf("%s:%s", rDir, path)
		log.WithField("path", path).Debug("adding Rpath to front of system PATH")
	}
	return path
}

func configureEnvAsUser(sysEnvVars []string, rs *RSettings) ([]string, error) {
	evs := environ.New(sysEnvVars)

	if path := prependPath(evs.Get("PATH"), rs.RPath); path != "" {
		evs.Set("PATH", path)
	}

	lp, ok := rs.LibPathsEnv()
	if ok {
		evs.Set("R_LIBS_SITE", lp)
	}

	return evs.AsSlice(), nil
}

// Returns the environment variables passed as a slice of name=value env strings
func censorEnvVars(nvp []string, add ...string) []string {
	var es struct{}
	var censoredString []string
	addlCensoredVars := make(map[string]struct{})
	if len(add) != 0 {
		for _, v := range add {
			addlCensoredVars[strings.ToUpper(v)] = es
		}
	}
	for _, v := range nvp {
		evs := strings.SplitN(v, "=", 2)
		_, present := censoredVars[strings.ToUpper(evs[0])]
		_, presentAddl := addlCensoredVars[strings.ToUpper(evs[0])]
		if present || presentAddl {
			censoredString = append(censoredString, fmt.Sprintf("%s=%s", evs[0], "***HIDDEN***"))
		} else {
			censoredString = append(censoredString, v)
		}
	}
	return censoredString
}
