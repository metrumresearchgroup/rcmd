package shiny

import (
	"context"
	"fmt"

	"github.com/metrumresearchgroup/command"

	"github.com/metrumresearchgroup/rcmd"
)

func LaunchApp(dir, name string, port int) (*command.Interact, error) {
	if port == 0 {
		port = 9999
	}
	rc := rcmd.NewRunConfig(rcmd.WithPrefix("foo"))
	rs, err := rcmd.NewRSettings("R")
	if err != nil {
		panic(err)
	}

	return rs.StartR(context.Background(), rc, dir, "-e", fmt.Sprintf("shiny::runApp('%s', port = %d)", name, port), "--slave", "--quiet")
}
