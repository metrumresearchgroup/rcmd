package shiny_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/metrumresearchgroup/rcmd/shiny"
)

func TestName(_ *testing.T) {
	i, err := shiny.LaunchApp("testdata", "app.R", 0)
	if err != nil {
		panic(err)
	}

	go func() {
		time.Sleep(20 * time.Second)
		if err := i.Stop(); err != nil {
			panic(err)
		}
	}()

	outScanner := i.StdoutScanner()
	errScanner := i.StderrScanner()

	wg := sync.WaitGroup{}

	wg.Add(2)

	go func() {
		defer wg.Done()

		for outScanner.Scan() {
			if _, err := fmt.Println(outScanner.Text()); err != nil {
				panic(err)
			}
		}
		fmt.Println("exit scan out")
	}()

	go func() {
		defer wg.Done()

		for errScanner.Scan() {
			if _, err := fmt.Println(errScanner.Text()); err != nil {
				panic(err)
			}
		}
		fmt.Println("exit scan err")
	}()

	wg.Wait()
}
