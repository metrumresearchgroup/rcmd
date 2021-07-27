package rcmd

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"os/exec"
	"runtime"

	"github.com/metrumresearchgroup/command"
	log "github.com/sirupsen/logrus"

	"github.com/metrumresearchgroup/rcmd/writers"
)

const defaultFailedCode = 1
const defaultSuccessCode = 0

// RunCfg contains the configuration for use when executing a run
type RunCfg struct {
	Stdout           io.Writer
	Stderr           io.Writer
	Stdin            io.Reader
	Prefix           string
	StripLineNumbers bool
	Script           bool
}

// RunOption are specific run option funcs to change configuration
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

// WithScript sets whether the output will be a call to R or Rscript
func WithScript(s bool) RunOption {
	return func(r *RunCfg) {
		r.Script = s
	}
}

// WithLineNumbers controls whether to keep the leading line numbers
// R includes in all outputs under the format [<num>] <output>
func WithLineNumbers(ln bool) RunOption {
	return func(r *RunCfg) {
		r.StripLineNumbers = !ln
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
	ctx context.Context,
	rs RSettings,
	dir string, // this should be put into RSettings
	cmdArgs []string,
	rc RunCfg,
) error {
	envVars := configureEnv(os.Environ(), &rs)
	capture := command.New(command.WithDir(dir), command.WithEnv(envVars))

	rpath := rs.R(runtime.GOOS, rc.Script)
	_, err := capture.Start(ctx, rpath, cmdArgs...)
	if err != nil {
		return err
	}

	cmd := exec.CommandContext(
		ctx,
		rs.R(runtime.GOOS, rc.Script),
		cmdArgs...,
	)

	cmd.Dir = dir
	cmd.Env = envVars
	cmd.Stdout = rc.Stdout
	cmd.Stderr = rc.Stderr
	cmd.Stdin = rc.Stdin
	return cmd.Run()
}

// RunR runs a non-interactive R command and streams back the results of
// the stderr and stdout to the RunCfg writers.
// RunR returns the exit code of the process and and error, if relevant
func RunR(
	ctx context.Context,
	rs RSettings,
	dir string,
	cmdArgs []string,
	rc RunCfg,
) (int, error) {
	envVars := configureEnv(os.Environ(), &rs)
	capture := command.New(command.WithDir(dir), command.WithEnv(envVars))

	rpath := rs.R(runtime.GOOS, rc.Script)
	p, err := capture.Start(ctx, rpath, cmdArgs...)
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		return 0, err
	}

	var fns []writers.FilterFunc

	if rc.StripLineNumbers {
		fns = append(fns, writers.LineNumberStripper)
	}
	fns = append(fns, bytes.TrimSpace)
	if rc.Prefix != "" {
		fns = append(fns, writers.NewPrefixFilter(rc.Prefix))
	}

	filter := writers.NewFilter(rc.Stdout, fns...)

	stdoutScanner := bufio.NewScanner(p.Stdout) // Notice that this is not in a loop
	go func() {
		for stdoutScanner.Scan() {
			_, err := filter.Write(stdoutScanner.Bytes())
			if err != nil {
				break
			}
		}
	}()

	stderrScanner := bufio.NewScanner(p.Stderr) // Notice that this is not in a loop
	go func() {
		for stderrScanner.Scan() {
			_, err := filter.Write(stdoutScanner.Bytes())
			if err != nil {
				break
			}
		}
	}()

	if err = p.Stdin.Close(); err != nil {
		return 0, err
	}

	if err = capture.Stop(); err != nil {
		return 0, err
	}

	return capture.ExitCode, nil
}

// RunRWithOutput runs a non-interactive R command and returns the combined output
func (rs *RSettings) RunRWithOutput(ctx context.Context, dir string, args ...string) (*command.Capture, error) {
	envVars := configureEnv(os.Environ(), rs)
	name := rs.R(runtime.GOOS, false)

	return run(ctx, envVars, dir, name, args...)
}

func run(ctx context.Context, env []string, dir string, name string, args ...string) (*command.Capture, error) {
	cmd := command.New(command.WithEnv(env), command.WithDir(dir))
	err := cmd.Run(ctx, name, args...)
	return cmd, err
}
