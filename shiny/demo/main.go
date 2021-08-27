package main

import (
	"context"
	"fmt"
	"time"

	"github.com/metrumresearchgroup/command"

	"github.com/metrumresearchgroup/rcmd/v2/shiny"
)

const timeout = 20 * time.Second

func main() {
	app, err := shiny.ConfigureApp(context.Background(), "testdata", "app.R", 0)
	if err != nil {
		panic(err)
	}

	// i := command.WireIO(nil, os.Stdout, os.Stderr)
	command.InteractiveIO().Apply(app.Cmd)

	err = app.Start()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Sleeping %d…\n", timeout)
	time.Sleep(timeout)
	if err := app.Kill(); err != nil {
		panic(err)
	}
}
