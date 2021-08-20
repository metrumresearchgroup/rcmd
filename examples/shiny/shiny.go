package main

import (
	"bufio"
	"context"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/metrumresearchgroup/rcmd/writers"
	"io"
	"os"
	"sync"
)

func main()  {
	rs, _ := rcmd.NewRSettings("R")
	rc := rcmd.NewRunConfig()


	inter, err  := rs.StartR(context.Background(), rc, "examples/shiny", "--quiet", "-e", "shiny::runApp('app.R', port = 8101)")
	if err != nil {
		panic(err)
	}
	stdoutScanner := inter.StdoutScanner()
	stderrScanner := inter.StderrScanner()
	wg := sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		for stdoutScanner.Scan() {
			outputPrinter(rc, stdoutScanner, os.Stdout)
		}
	}()
	go func() {
		wg.Add(1)
		defer wg.Done()
		for stderrScanner.Scan() {
			outputPrinter(rc, stderrScanner, os.Stderr)
		}
	}()

	inter.Wait()
	wg.Wait()
}

func outputPrinter(rc *rcmd.RunCfg, b *bufio.Scanner, w io.Writer) {
	var text string
	if rc.StripLineNumbers {
		text = string(writers.InputFilter(writers.LineNumberStripper(b.Bytes())))
	}
	if rc.Prefix != "" {
		fmt.Fprintln(w, rc.Prefix, text)
	} else {
		fmt.Fprintln(w, text)
	}
}