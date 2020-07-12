package rcmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/metrumresearchgroup/rcmd/rp"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
	"os/exec"
	"runtime"
)

const defaultFailedCode = 1
const defaultSuccessCode = 0

// RunOpts
type RunCfg struct {
	Stdout           io.Writer
	Stderr           io.Writer
	Stdin            io.Reader
	Prefix           string
	StripLineNumbers bool
	Script           bool
}

type RunOption func(*RunCfg)

func NewRunConfig(opts ...RunOption) *RunCfg {
	cfg := &RunCfg{
		Stdout:           os.Stdout,
		Stderr:           os.Stderr,
		Stdin:            os.Stdin,
		Prefix:           "",
		StripLineNumbers: true,
		Script:           false,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	return cfg
}

func WithScript(s bool) RunOption {
	return func(r *RunCfg) {
		r.Script = s
	}
}
// WithLineNumbers controls whether to keep the leading line numbers
// R includes in all outputs under the format [<num>] <output>
func WithLineNumbers(ln bool) RunOption {
	return func(r *RunCfg) {
		r.StripLineNumbers = ln
	}
}

// WithStdOut Sets the writer to send results sent to stdout
// for example to suppress stdout could provide `WithStdOut(ioutil.Discard)`
func WithStdOut(w io.Writer) RunOption {
	return func(r *RunCfg) {
		r.Stdout = w
	}
}
// WithStdErr Sets the writer to send results sent to stderr
// for example to suppress stderr could provide `WithStdOut(ioutil.Discard)`
func WithStdErr(w io.Writer) RunOption {
	return func(r *RunCfg) {
		r.Stdout = w
	}
}

// WithStdin Sets the writer to send information the stdin. Not all R commands
// accept stdin.
func WithStdin(r io.Reader) RunOption {
	return func(rc *RunCfg) {
		rc.Stdin = r
	}
}
// WithPrefix sets a prefix string before any message is sent to the stdout/stderr writers
// This is useful when printing concurrent results out and want to differentiate
// where messages are coming from
func WithPrefix(prefix string) RunOption {
	return func(r *RunCfg) {
		r.Prefix = prefix
	}
}

// StartR launches an interactive R console given the same
// configuration as a specific package.
func StartR(
	rs RSettings,
	dir string, // this should be put into RSettings
	cmdArgs []string,
	rc RunCfg,
) error {

	envVars := configureEnv(os.Environ(), rs)

	log.WithFields(
		log.Fields{
			"cmdArgs":   cmdArgs,
			"RSettings": rs,
			"env":       envVars,
		}).Trace("command args")

	// --vanilla is a command for R and should be specified before the CMD, eg
	// R --vanilla CMD check
	// if cs.Vanilla {
	// 	cmdArgs = append([]string{"--vanilla"}, cmdArgs...)
	// }
	cmd := exec.Command(
		rs.R(runtime.GOOS, rc.Script),
		cmdArgs...,
	)

	if dir == "" {
		dir, _ = os.Getwd()
	}
	cmd.Dir = dir
	cmd.Env = envVars
	cmd.Stdout = rc.Stdout
	cmd.Stderr = rc.Stderr
	cmd.Stdin = rc.Stdin
	return cmd.Run()
}

// RunRWithOutput runs a non-interactive R command and streams back the results of
// the stderr and stdout to the RunCfg writers
func RunRWithOutput(
	ctx context.Context,
	rs RSettings,
	dir string,
	cmdArgs []string,
	rc RunCfg,
) (int, error) {
	envVars := configureEnv(os.Environ(), rs)
	rpath := rs.R(runtime.GOOS, false)
	cmd := exec.CommandContext(
		ctx,
		rpath,
		cmdArgs...,
	)
	if dir == "" {
		dir, _ = os.Getwd()
	}
	cmd.Dir = dir
	cmd.Env = envVars
	stdoutPipe, err := cmd.StdoutPipe()
	stderrPipe, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	err = cmd.Start()
	if err != nil {
		fmt.Println("error starting cmd")
		panic(err)
	}
	stdoutScanner := bufio.NewScanner(stdoutPipe) // Notice that this is not in a loop
	stderrScanner := bufio.NewScanner(stderrPipe) // Notice that this is not in a loop
	go func() {
		for stdoutScanner.Scan() {
			text := stdoutScanner.Text()
			if rc.StripLineNumbers {
				text = rp.StripLineNumber(text)
			}
			if rc.Prefix != "" {
				fmt.Fprintln(rc.Stdout, rc.Prefix, text)
			} else {
				fmt.Fprintln(rc.Stdout, text)
			}
		}
	}()
	go func() {
		for stderrScanner.Scan() {
			text := stderrScanner.Text()
			if rc.StripLineNumbers {
				text = rp.StripLineNumber(text)
			}
			if rc.Prefix != "" {
				fmt.Fprintln(rc.Stdout, rc.Prefix, text)
			} else {
				fmt.Fprintln(rc.Stdout, text)
			}
		}
	}()
	if err := cmd.Wait(); err != nil {
		return cmd.ProcessState.ExitCode(), err
	}
	return cmd.ProcessState.ExitCode(), nil
}

// RunR runs a non-interactive R command and returns the combined output
func RunR(
	rs RSettings,
	dir string,
	cmdArgs []string,
) ([]byte, error) {
	envVars := configureEnv(os.Environ(), rs)
	rpath := rs.R(runtime.GOOS, false)
	cmd := exec.Command(
		rpath,
		cmdArgs...,
	)
	if dir == "" {
		dir, _ = os.Getwd()
	}
	cmd.Dir = dir
	cmd.Env = envVars

	return cmd.CombinedOutput()
}
