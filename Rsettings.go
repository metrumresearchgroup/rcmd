package rcmd

import (
	"context"
	"errors"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/metrumresearchgroup/command"

	"github.com/metrumresearchgroup/rcmd/rp"
)

// NewRSettings initializes RSettings.
func NewRSettings(rPath string) (*RSettings, error) {
	rs := RSettings{
		RPath: rPath,
	}
	// since we have the path in the constructor, we might as well get the
	// R version now too

	rv, err := getRVersion(rPath)
	if err != nil {
		return nil, err
	}

	rs.Version = *rv

	return &rs, nil
}

// R provides a cleaned path to the R executable.
func (rs RSettings) R(os string, script bool) string {
	r := "R"

	if rs.RPath != "" {
		// Need to trim trailing slash as will form the R CMD syntax
		// eg /path/to/R CMD, so can't have /path/to/R/ CMD
		r = strings.TrimSuffix(rs.RPath, "/")
	}
	// TODO: check if this could have problems with trailing slash on windows
	// TODO: better to use something like filepath.clean? would that sanitize better?
	// filepath.Clean does not remove trailing \ on mac. maybe it works on windows?
	r = filepath.Clean(r)

	if script {
		r += "script"
	}

	if os == "windows" && !strings.HasSuffix(r, ".exe") {
		r = r + ".exe"
	}

	return r
}

var rversion RVersion
var rplatform string

func getRVersion(rPath string) (*RVersion, error) {
	if rversion.ToString() != "0.0" {
		version := rversion

		return &version, nil
	}

	if rPath == "" {
		rPath = "R"
	}

	cmd := command.NewWithContext(context.Background(), rPath, "--quiet", "--version", "--vanilla")
	co, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var version *RVersion
	var platform string

	version, platform, err = parseVersionData(co)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("parseVersionData returned nil version")
	}

	rversion = *version
	rplatform = platform

	return version, nil
}

// getRVersion returns the R version, and sets R Version and R platform in RSettings
// unlike the other methods, this one is a pointer, as RVersion mutates the known R Version,
// as if it is not defined, it will shell out to R to determine the version, and mutate itself
// to set that value, while also returning the RVersion.
// This will keep any program using rs from needing to shell out multiple times.
func (rs *RSettings) getRVersion() (*RVersion, error) {
	if rs.Version.ToString() != "0.0" {
		version := rs.Version

		return &version, nil
	}

	cmd := command.NewWithContext(context.Background(), "R", "--quiet", "--version", "--vanilla")
	co, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	var version *RVersion
	var platform string

	version, platform, err = parseVersionData(co)
	if err != nil {
		return nil, err
	}
	if version == nil {
		return nil, errors.New("parseVersionData returned nil version")
	}

	rs.Version = *version
	rs.Platform = platform

	return version, nil
}

func parseVersionData(data []byte) (version *RVersion, platform string, err error) {
	lines, err := rp.ScanLines(data)
	if err != nil {
		return nil, "", err
	}

	for _, line := range lines {
		if strings.HasPrefix(line, "R version") {
			spl := strings.Split(line, " ")
			if len(spl) < 3 {
				return nil, "", errors.New("error getting R version")
			}
			rsp := strings.Split(spl[2], ".")
			if len(rsp) == 3 {
				maj, _ := strconv.Atoi(rsp[0])
				min, _ := strconv.Atoi(rsp[1])
				pat, _ := strconv.Atoi(rsp[2])
				// this should now make it so in the future it will be set so should only need to shell out to R once
				version = &RVersion{
					Major: maj,
					Minor: min,
					Patch: pat,
				}
			} else {
				return nil, "", errors.New("error getting R version")
			}
		} else if strings.HasPrefix(line, "Platform:") {
			rsp := strings.Split(line, " ")
			if len(rsp) > 0 {
				platform = strings.Trim(rsp[1], " ")
			} else {
				return nil, "", errors.New("error getting R Platform")
			}
		}
	}

	return version, platform, nil
}

// LibPathsEnv returns the library formatted in the style to be set as an environment variable.
func (rs RSettings) LibPathsEnv() (string, bool) {
	if len(rs.LibPaths) == 0 {
		return "", false
	}
	if len(rs.LibPaths) == 1 && rs.LibPaths[0] == "" {
		return "", false
	}

	return strings.Join(rs.LibPaths, ":"), true
}
