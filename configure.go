package rcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/metrumresearchgroup/environ"
	log "github.com/sirupsen/logrus"
)

// ConfigureEnv contains the default environment variables usually from
// os.Environ().
func ConfigureEnv(env []string, rs *RSettings) ([]string, error) {
	if rs.AsUser {
		return configureEnvAsUser(env, rs)
	} else {
		return configureEnvAsNotUser(env, rs)
	}
}

// sysEnvVars contains the default environment variables usually from
// os.Environ().
func configureEnvAsNotUser(env []string, rs *RSettings) ([]string, error) {
	evs := environ.New(env)
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
