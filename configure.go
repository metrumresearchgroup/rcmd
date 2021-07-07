package rcmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

var censoredVars = map[string]string{
	"GITHUB_TOKEN":      "GITHUB_TOKEN",
	"GITHUB_PAT":        "GITHUB_PAT",
	"GHE_TOKEN":         "GHE_TOKEN",
	"GHE_PAT":           "GHE_PAT",
	"AWS_ACCESS_KEY_ID": "AWS_ACCESS_KEY_ID",
	"AWS_SECRET_KEY":    "AWS_SECRET_KEY",
}

// sysEnvVars contains the default environment variables usually from
// os.Environ()
func configureEnv(sysEnvVars []string, rs RSettings) []string {
	envList := NvpList{}
	envVars := []string{}

	for _, p := range rs.EnvVars.Pairs {
		_, exists := envList.Get(p.Name)
		if !exists {
			envList = NvpAppend(envList, p.Name, p.Value)
		}
	}
	// system env vars generally
	for _, ev := range sysEnvVars {
		evs := strings.SplitN(ev, "=", 2)
		if len(evs) > 1 && evs[1] != "" {

			// we don't want to track the order of these anyway since they should take priority in the end
			// R_LIBS_USER takes precedence over R_LIBS_SITE
			// so will cause the loading characteristics to
			// not be representative of the hierarchy specified
			// in Library/Libpaths in the pkgr configuration
			// we only want R_LIBS_SITE set to control all relevant library paths for the user to
			if !rs.AsUser && evs[0] == "R_LIBS_USER" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_USER")
				continue
			}
			if !rs.AsUser && evs[0] == "R_LIBS_SITE" {
				log.WithField("path", evs[1]).Debug("overriding system R_LIBS_SITE")
				continue
			}
			if evs[0] == "PATH" {
				rDir := filepath.Dir(rs.RPath)
				if rDir != "" && rDir != "." && !strings.HasPrefix(evs[1], rDir) {
					evs[1] = fmt.Sprintf("%s:%s", rDir, evs[1])
					log.WithField("path", evs[1]).Debug("adding Rpath to front of system PATH")
				}
			}
			// if exists would be custom to the package hence should not accept the system env
			envList = NvpAppend(envList, evs[0], evs[1])
		}
	}

	if !rs.AsUser {
		tmpdir, err := os.MkdirTemp("", "")
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
			}).Warn("error making temporary directory while overriding R_LIBS_USER for install.")
		}
		// Force R_LIBS_USER to be an non-empty dir so that we can be sure it won't get overridden by default R paths.
		envList = NvpAppend(envList, "R_LIBS_USER", tmpdir)

		ok, lp := rs.LibPathsEnv()
		if ok {
			envList = NvpAppendPair(envList, lp)
		}
	}

	for _, p := range envList.Pairs {
		// the one and only place to append name=value strings to envVars
		envVars = append(envVars, p.GetString(p.Name))
	}

	return envVars
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
