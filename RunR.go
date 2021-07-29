package rcmd

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"runtime"

	"github.com/metrumresearchgroup/command"
	"github.com/metrumresearchgroup/command/pipes"

	"github.com/metrumresearchgroup/rcmd/writers"
)

// RunCfg contains the configuration for use when executing a run
type RunCfg struct {
	// Stdout           io.ReadCloser
	// Stderr           io.ReadCloser
	// Stdin            io.WriteCloser
	Prefix           string
	StripLineNumbers bool
	Script           bool
}

// RunOption are specific run option funcs to change configuration
type RunOption func(*RunCfg)

func NewRunConfig(opts ...RunOption) *RunCfg {
	cfg := &RunCfg{
		// Stdout:           os.Stdout,
		// Stderr:           os.Stderr,
		// Stdin:            os.Stdin,
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

//
// // WithStdOut Sets the writer to send results sent to stdout
// // for example to suppress stdout could provide `WithStdOut(ioutil.Discard)`
// func WithStdOut(r io.ReadCloser) RunOption {
// 	return func(rc *RunCfg) {
// 		rc.Stdout = r
// 	}
// }
//
// // WithStdErr Sets the writer to send results sent to stderr
// // for example to suppress stderr could provide `WithStdOut(ioutil.Discard)`
// func WithStdErr(r io.ReadCloser) RunOption {
// 	return func(rc *RunCfg) {
// 		rc.Stdout = r
// 	}
// }
//
// // WithStdin Sets the writer to send information the stdin. Not all R commands
// // accept stdin.
// func WithStdin(w io.WriteCloser) RunOption {
// 	return func(rc *RunCfg) {
// 		rc.Stdin = w
// 	}
// }

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
) (*pipes.Pipes, error) {
	envVars, err := configureEnv(os.Environ(), &rs)
	if err != nil {
		return nil, err
	}
	capture := command.New(command.WithDir(dir), command.WithEnv(envVars))

	rpath := rs.R(runtime.GOOS, rc.Script)
	p, err := capture.Start(ctx, rpath, cmdArgs...)
	if err != nil {
		return nil, err
	}

	return p, nil
}

// RunR runs a non-interactive R command and streams back the results of
// the stderr and stdout to the RunCfg writers.
// RunR returns the exit code of the process and and error, if relevant
func RunR(
	ctx context.Context,
	rs *RSettings,
	dir string,
	cmdArgs []string,
	rc *RunCfg,
) (*pipes.Pipes, int, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, 0, err
	}
	capture := command.New(command.WithDir(dir), command.WithEnv(envVars))

	rpath := rs.R(runtime.GOOS, rc.Script)
	p, err := capture.Start(ctx, rpath, cmdArgs...)
	if err != nil {
		return p, 0, err
	}

	// Assign pipe's original out err to locals
	sout, serr := p.Stdout, p.Stderr

	// Read both streams into one
	mr := io.MultiReader(sout, serr)

	// filter results

	var fns []writers.FilterFunc

	if rc.StripLineNumbers {
		fns = append(fns, writers.LineNumberStripper)
	}
	fns = append(fns, bytes.TrimSpace)
	if rc.Prefix != "" {
		fns = append(fns, writers.NewPrefixFilter(rc.Prefix))
	}

	// Create a pipe to keep stdout fed with merged out/err stream
	pipeReader, PipeWriter := io.Pipe()

	writeFilter := writers.NewFilter(PipeWriter, fns...)

	// re-assign the pipe to stdout on p
	p.Stdout = pipeReader

	scanner := bufio.NewScanner(mr) // Notice that this is not in a loop
	go func() {
		for scanner.Scan() {
			_, err := writeFilter.Write(scanner.Bytes())
			if err != nil {
				break
			}
		}
	}()

	if err = p.Stdin.Close(); err != nil {
		return p, 0, err
	}

	if err = capture.Stop(); err != nil {
		return p, 0, err
	}

	return p, capture.ExitCode, nil
}

// RunRWithOutput runs a non-interactive R command and returns the combined output
func (rs *RSettings) RunRWithOutput(ctx context.Context, dir string, args ...string) ([]byte, *command.Capture, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, nil, err
	}
	name := rs.R(runtime.GOOS, false)

	return combinedOutput(ctx, envVars, dir, name, args...)
}

func run(ctx context.Context, env []string, dir string, name string, args ...string) (*pipes.Pipes, *command.Capture, error) {
	cmd := command.New(command.WithEnv(env), command.WithDir(dir))
	ps, err := cmd.Run(ctx, name, args...)
	return ps, cmd, err
}

func combinedOutput(ctx context.Context, env []string, dir string, name string, args ...string) ([]byte, *command.Capture, error) {
	capture := command.New(command.WithEnv(env), command.WithDir(dir))
	co, err := capture.CombinedOutput(ctx, name, args...)

	return co, capture, err
}

type RCapture struct {
	*command.Capture
}

func (rs *RSettings) Run(ctx context.Context, rc *RunCfg, dir string, args ...string) (*pipes.Pipes, *command.Capture, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, nil, err
	}
	name := rs.R(runtime.GOOS, false)
	return run(ctx, envVars, dir, name, args...)
}
