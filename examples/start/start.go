package main

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/metrumresearchgroup/rcmd/rp"
	"github.com/mitchellh/go-homedir"
	log "github.com/sirupsen/logrus"
	"strings"
	"sync"
	"time"
)

func main() {
	//startR_example()
	//runR_expression()
	runR_exampleCancel()
}

func startR_example() {
	if err := rcmd.StartR(rcmd.NewRSettings("R"), "", []string{}, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}
}

func startR_example2() {
	if err := rcmd.StartR(rcmd.NewRSettings("R"), "", []string{"-e", "2+2", "slave"}, *rcmd.NewRunConfig()); err != nil {
		panic(err)
	}
}

func runR_expression() {
	res, err := rcmd.RunRWithOutput(rcmd.NewRSettings("R"), "", []string{"-e", "2+2", "--slave"})
	if err != nil {
		panic(err)
	}
	fmt.Println(strings.Join(rp.ScanLines(res), "\n"))
}
func runRWithOutput_example() {
	res, err := rcmd.RunR(context.Background(), rcmd.NewRSettings("R"), "", []string{"-e", "2+2", "--slave", "--interactive"}, *rcmd.NewRunConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	fmt.Println("----with prefix----")
	res, err = rcmd.RunR(context.Background(), rcmd.NewRSettings("R"), "", []string{"-e", "2+2", "--slave", "--interactive"}, *rcmd.NewRunConfig(rcmd.WithPrefix("custom-prefix:")))
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
func runR_exampleTimeout() {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)
	res, err := rcmd.RunR(ctx, rcmd.NewRSettings("R"), "", []string{"-e", "Sys.sleep(1.5); 2+2", "--slave", "--interactive"}, *rcmd.NewRunConfig())
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
		res, err := rcmd.RunR(ctx, rcmd.NewRSettings("R"), "", []string{"-e", "2+2", "--slave", "--interactive"}, *rcmd.NewRunConfig())

		if err != nil {
			panic(err)
		}
		fmt.Println("go routine 1: ", res)
	}()
	go func() {
		defer wg.Done()
		res, err := rcmd.RunR(ctx, rcmd.NewRSettings("R"), "", []string{"-e", "Sys.sleep(0.5); stop('failed')", "--slave", "--interactive"}, *rcmd.NewRunConfig())
		if err != nil {
			log.Error("goroutine 2 error:", err)
			log.Warn("cancelling ongoing work...")
			cancel()
		}
		fmt.Println("go routine 2 res: ", res)

	}()
	go func() {
		defer wg.Done()
		res, err := rcmd.RunR(ctx, rcmd.NewRSettings("R"), "", []string{"-e", "Sys.sleep(1); 2+2", "--slave", "--interactive"}, *rcmd.NewRunConfig())
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
	res, err := rcmd.RunR(context.Background(), rcmd.NewRSettings("R"), dir, []string{"-e", "options(crayon.enabled = TRUE); devtools::test()", "--slave", "--interactive"}, *rcmd.NewRunConfig())
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
