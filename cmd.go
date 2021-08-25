package rcmd

import (
	"context"
	"os"
	"runtime"

	"github.com/metrumresearchgroup/command"
)

type RCmd struct {
	*command.Cmd

	Prefix         string
	UseLineNumbers bool
	Script         bool
	RSettings      *RSettings
}

func NewRScript(ctx context.Context, dir string, args ...string) (*RCmd, error) {
	return newRCmd(ctx, true, dir, args...)
}

func newRCmd(ctx context.Context, script bool, dir string, args ...string) (*RCmd, error) {
	rcmd := &RCmd{
		Script: script,
	}

	rs, err := NewRSettings("")
	if err != nil {
		return nil, err
	}

	env, err := configureEnv(os.Environ(), rs)
	if err != nil {
		return nil, err
	}

	rpath := rs.R(runtime.GOOS, rcmd.Script)

	cmd := command.NewWithContext(ctx, rpath, args...)
	cmd.Dir = dir
	cmd.Env = env

	rcmd.Cmd = cmd

	return rcmd, nil
}

func New(ctx context.Context, dir string, args ...string) (*RCmd, error) {
	return newRCmd(ctx, false, dir, args...)
}
