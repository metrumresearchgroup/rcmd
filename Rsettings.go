package rcmd

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/metrumresearchgroup/command"
	"github.com/spf13/afero"

	"github.com/metrumresearchgroup/rcmd/v2/rp"
)

// NewRSettings initializes RSettings.
func NewRSettings(rPath string) (*RSettings, error) {
	rs := RSettings{
		RPath: rPath,
	}
	// since we have the path in the constructor, we might as well get the
	// R version now too

	rv, rpl, rpa, err := GetRVersionPlatformPath(rPath)
	if err != nil {
		return nil, err
	}

	rs.Version = *rv
	rs.Platform = string(rpl)
	rs.RPath = string(rpa)

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

type RPath string

func (rp RPath) Exists(fs afero.Fs) (bool, error) {
	return afero.Exists(fs, string(rp))
}

// GetRVersionPlatformPath returns the R version, and sets R Version and R platform in global state.
// unlike the other methods, this one is a pointer, as RVersion mutates the known R Version,
// as if it is not defined, it will shell out to R to determine the version, and mutate itself
// to set that value, while also returning the RVersion.
// This will keep any program using rs from needing to shell out multiple times.
func GetRVersionPlatformPath(rPath string) (*RVersion, Platform, RPath, error) {
	if rPath == "" {
		rPath = "R"
	}

	cmd := command.NewWithContext(context.Background(), rPath, "--quiet", "--version", "--vanilla")
	co, err := cmd.CombinedOutput()
	if err != nil {
		return nil, "", "", err
	}

	var version *RVersion
	var platform Platform

	version, platform, err = parseVersionData(co)
	if err != nil {
		return nil, "", "", err
	}

	return version, platform, RPath(rPath), nil
}

type Platform string

func parseVersionData(data []byte) (*RVersion, Platform, error) {
	lines := rp.ScanLines(data)

	var plat []byte
	var version *RVersion
	var err error

	scanner := bufio.NewScanner(bytes.NewReader(lines))
	for scanner.Scan() {
		line := scanner.Bytes()
		if bytes.HasPrefix(line, []byte("R version")) {
			version, err = VersionLine(line)
			if err != nil {
				return nil, "", err
			}
			// bail if platform is already set
			if len(plat) > 0 {
				break
			}
		} else if bytes.HasPrefix(line, []byte("Platform:")) {
			rsp := bytes.Split(line, []byte(" "))
			if len(rsp) > 0 {
				plat = bytes.Trim(rsp[1], " ")
				// bail if version is already set
				if version != nil {
					break
				}
			} else {
				return nil, "", errors.New("error getting R Platform")
			}
		}
	}

	return version, Platform(plat), nil
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

// VersionLine reads a line starting with "R Version" and parses out the version
// information.
func VersionLine(line []byte) (*RVersion, error) {
	spl := bytes.Split(line, []byte(" "))
	if len(spl) < 3 {
		return nil, errors.New("error getting R version")
	}

	semvers := bytes.Split(spl[2], []byte("."))
	if len(semvers) == 3 {
		var err error
		version := RVersion{}

		if version.Major, err = strconv.Atoi(string(semvers[0])); err != nil {
			return nil, fmt.Errorf("couldn't parse major: %w", err)
		}
		if version.Minor, err = strconv.Atoi(string(semvers[1])); err != nil {
			return nil, fmt.Errorf("couldn't parse minor: %w", err)
		}
		if version.Patch, err = strconv.Atoi(string(semvers[2])); err != nil {
			return nil, fmt.Errorf("couldn't parse patch: %w", err)
		}

		return &version, nil
	}

	return nil, errors.New("error getting R version")
}
