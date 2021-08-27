package shiny

import (
	"context"
	"fmt"

	"github.com/metrumresearchgroup/rcmd/v2"
)

func ConfigureApp(ctx context.Context, dir, name string, port int) (*rcmd.RCmd, error) {
	if port == 0 {
		port = 9999
	}

	rc, err := rcmd.New(ctx, dir, "-e", fmt.Sprintf("shiny::runApp('%s', port = %d)", name, port), "--slave", "--quiet")
	if err != nil {
		return nil, err
	}

	return rc, nil
}
