package main

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"

	"github.com/metrumresearchgroup/rcmd"
	"github.com/metrumresearchgroup/rcmd/rp"
)

func main() {
	// startR_example()
	// runR_expression()
	runRWithOutput_example()
	// runR_exampleCancel()
}

func startR_example() {
	rs, err := rcmd.NewRSettings("R")
	if err != nil {
		panic(err)
	}
	if _, err := rcmd.StartR(context.Background(), *rs, "", []string{}, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}
}

func startR_example2() {
	rs, err := rcmd.NewRSettings("R")
	if err != nil {
		panic(err)
	}
	if _, err := rcmd.StartR(context.Background(), *rs, "", []string{"-e", "2+2", "slave"}, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}
}

func runR_expression() {
	rs, err := rcmd.NewRSettings("R")
	if err != nil {
		panic(err)
	}

	ps, _, err := rs.RunRWithOutput(context.Background(), "", "-e", "2+2", "--slave")
	if err != nil {
		panic(err)
	}

	lines, err := rp.ScanLines(ps)
	fmt.Println(strings.Join(lines, "\n"))
}
func runRWithOutput_example() {
	rs, err := rcmd.NewRSettings("R")
	bs, _, err := rs.RunRWithOutput(context.Background(), "", "-e", "2+2", "--slave", "--interactive")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
	fmt.Println("----with prefix----")
	rs, err = rcmd.NewRSettings("R")
	bs, _, err = rs.RunRWithOutput(context.Background(), "", "-e", "2+2", "--slave", "--interactive")
	if err != nil {
		panic(err)
	}
	fmt.Println(string(bs))
}

func runR_exampleTimeout() {
	rs, err := rcmd.NewRSettings("R")
	ctx, ccl := context.WithTimeout(context.Background(), 1*time.Second)
	defer ccl()
	_, res, err := rcmd.RunR(ctx, rs, "", []string{"-e", "Sys.sleep(1.5); 2+2", "--slave", "--interactive"}, rcmd.NewRunConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}

func runR_exampleCancel() {
	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// fire off 3 go routines where the 2nd fails, and all existing should stop processing
	wg.Add(3)
	go func() {
		defer wg.Done()
		rs, err := rcmd.NewRSettings("R")
		_, res, err := rcmd.RunR(ctx, rs, "", []string{"-e", "2+2", "--slave", "--interactive"}, rcmd.NewRunConfig())

		if err != nil {
			panic(err)
		}
		fmt.Println("go routine 1: ", res)
	}()
	go func() {
		defer wg.Done()
		rs, err := rcmd.NewRSettings("R")
		_, res, err := rcmd.RunR(ctx, rs, "", []string{"-e", "Sys.sleep(0.5); stop('failed')", "--slave", "--interactive"}, rcmd.NewRunConfig())
		if err != nil {
			log.Error("goroutine 2 error:", err)
			log.Warn("cancelling ongoing work...")
			cancel()
		}
		fmt.Println("go routine 2 res: ", res)

	}()
	go func() {
		defer wg.Done()
		rs, err := rcmd.NewRSettings("R")

		_, res, err := rcmd.RunR(ctx, rs, "", []string{"-e", "Sys.sleep(1); 2+2", "--slave", "--interactive"}, rcmd.NewRunConfig())
		if err != nil {
			log.Error("goroutine 3 error:", err)
		}
		fmt.Println("go routine 3 res: ", res)
	}()
	wg.Wait()
	fmt.Println("completed everything....")
}

func runR_examplepkg() {
	dir, _ := homedir.Expand("~/metrum/metrumresearchgroup/rbabylon")
	rs, err := rcmd.NewRSettings("R")
	_, res, err := rcmd.RunR(context.Background(), rs, dir, []string{"-e", "options(crayon.enabled = TRUE); devtools::test()", "--slave", "--interactive"}, rcmd.NewRunConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
