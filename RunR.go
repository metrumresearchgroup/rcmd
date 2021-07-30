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

// RunCfg contains the configuration for use when executing a run.
type RunCfg struct {
	Prefix           string
	StripLineNumbers bool
	Script           bool
}

// RunOption are specific run option funcs to change configuration.
type RunOption func(*RunCfg)

func NewRunConfig(opts ...RunOption) *RunCfg {
	cfg := &RunCfg{
		Prefix:           "",
		StripLineNumbers: true,
		Script:           false,
	}
	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

// WithScript sets whether the output will be a call to R or Rscript.
func WithScript(s bool) RunOption {
	return func(r *RunCfg) {
		r.Script = s
	}
}

// WithLineNumbers controls whether to keep the leading line numbers
// R includes in all outputs under the format [<num>] <output>.
func WithLineNumbers(ln bool) RunOption {
	return func(r *RunCfg) {
		r.StripLineNumbers = !ln
	}
}

// WithPrefix sets a prefix string before any message is sent to the stdout/stderr writers
// This is useful when printing concurrent results out and want to differentiate
// where messages are coming from.
func WithPrefix(prefix string) RunOption {
	return func(r *RunCfg) {
		r.Prefix = prefix
	}
}

// StartR launches an interactive R console given the same
// configuration as a specific package.
func (rs *RSettings) StartR(ctx context.Context, rc *RunCfg, dir string, cmdArgs ...string) (*pipes.Pipes, func() error, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, nil, err
	}
	capture := command.New(command.WithDir(dir), command.WithEnv(envVars))

	rpath := rs.R(runtime.GOOS, rc.Script)
	p, err := capture.Start(ctx, rpath, cmdArgs...)
	if err != nil {
		return nil, nil, err
	}

	return p, func() error {
		return capture.Stop()
	}, nil
}

func (rc *RunCfg) configPipe(p *pipes.Pipes) {
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
}

// RunR runs a non-interactive R command and streams back the results of
// the stderr and stdout to the *pipes.Pipes writers.
// RunR returns the exit code of the process and and error, if relevant.
func (rs *RSettings) RunR(ctx context.Context, rc *RunCfg, dir string, cmdArgs ...string) (*pipes.Pipes, int, error) {
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

	rc.configPipe(p)

	if err = p.Stdin.Close(); err != nil {
		return p, 0, err
	}

	if err = capture.Stop(); err != nil {
		return p, 0, err
	}

	return p, capture.ExitCode, nil
}

// RunRWithOutput runs a non-interactive R command and returns the combined output.
func (rs *RSettings) RunRWithOutput(ctx context.Context, rc *RunCfg, dir string, args ...string) ([]byte, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, err
	}
	name := rs.R(runtime.GOOS, rc.Script)

	return combinedOutput(ctx, envVars, dir, name, args...)
}

func run(ctx context.Context, env []string, dir string, name string, args ...string) (*pipes.Pipes, *command.Capture, error) {
	cmd := command.New(command.WithEnv(env), command.WithDir(dir))
	ps, err := cmd.Run(ctx, name, args...)

	return ps, cmd, err
}

func combinedOutput(ctx context.Context, env []string, dir string, name string, args ...string) ([]byte, error) {
	capture := command.New(command.WithEnv(env), command.WithDir(dir))
	co, err := capture.CombinedOutput(ctx, name, args...)

	return co, err
}

type RCapture struct {
	*command.Capture
}

func (rs *RSettings) Run(ctx context.Context, _ *RunCfg, dir string, args ...string) (*pipes.Pipes, *command.Capture, error) {
	envVars, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, nil, err
	}
	name := rs.R(runtime.GOOS, false)

	return run(ctx, envVars, dir, name, args...)
}
