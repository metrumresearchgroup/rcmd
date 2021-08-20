package main

import (
	"context"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/metrumresearchgroup/rcmd/writers"
	"os"
	"os/exec"
)

func main()  {
	rs, _ := rcmd.NewRSettings("R")
	rc :=rcmd.NewRunConfig()
	rc.CmdModifierFunc = func(cmd *exec.Cmd) error {
		cmd.Env = append(cmd.Env, "FOO=bar")
		return nil
	}
	res, err  := rs.RunRWithOutput(context.Background(), rc, "",  "--quiet", "-e", "Sys.getenv('FOO')")
	if err != nil {
		panic(err)
	}
	outputter := writers.NewFilter(os.Stdout, writers.LineNumberStripper, writers.InputFilter)
	outputter.Write(res)
}
