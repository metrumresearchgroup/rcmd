// +Build R
package rcmd

import (
	"context"
	"fmt"
	"io"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/metrumresearchgroup/wrapt"
	log "github.com/sirupsen/logrus"

	"github.com/metrumresearchgroup/rcmd/rp"
)

func Test_startR_example(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	p, stop, err := rs.StartR(context.Background(), NewRunConfig(WithPrefix("foo")), "")
	t.A.NoError(err)
	defer fmt.Println(stop())
	_, err = io.ReadAll(p.Stdout)
	t.A.NoError(err)
}

func Test_startR_example2(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	_, stop, err := rs.StartR(context.Background(), NewRunConfig(), "", "-e", "2+2", "slave")
	t.A.NoError(err)
	defer fmt.Println(stop())
}

func Test_runR_expression(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	ps, err := rs.RunRWithOutput(context.Background(), NewRunConfig(), "", "-e", "2+2", "--slave")
	t.A.NoError(err)

	lines, err := rp.ScanLines(ps)
	t.A.NoError(err)

	fmt.Println(strings.Join(lines, "\n"))
}

func Test_runRWithOutput_example(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	bs, err := rs.RunRWithOutput(context.Background(), NewRunConfig(), "", "-e", "2+2", "--slave", "--interactive")
	t.A.NoError(err)

	fmt.Println(string(bs))

	fmt.Println("----with prefix----")

	rs, err = NewRSettings("R")
	t.A.NoError(err)

	bs, err = rs.RunRWithOutput(context.Background(), NewRunConfig(), "", "-e", "2+2", "--slave", "--interactive")
	t.A.NoError(err)

	fmt.Println(string(bs))
}

func Test_runR_exampleTimeout(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	ctx, ccl := context.WithTimeout(context.Background(), 1*time.Second)
	defer ccl()

	p, res, err := rs.RunR(ctx, NewRunConfig(), "", "-e", "Sys.sleep(1.5); 2+2", "--slave", "--interactive")
	t.A.NoError(err)
	_, err = io.ReadAll(p.Stdout)
	t.A.NoError(err)

	fmt.Println(res)
}

func Test_runR_exampleCancel(tt *testing.T) {
	t := wrapt.WrapT(tt)

	ctx, cancel := context.WithCancel(context.Background())
	var wg sync.WaitGroup
	// fire off 3 go routines where the 2nd fails, and all existing should stop processing
	wg.Add(3)
	go func() {
		defer wg.Done()
		rs, err := NewRSettings("R")
		t.A.NoError(err)

		_, res, err := rs.RunR(ctx, NewRunConfig(), "", "-e", "2+2", "--slave", "--interactive")
		t.A.NoError(err)

		fmt.Println("go routine 1: ", res)
	}()
	go func() {
		defer wg.Done()

		rs, err := NewRSettings("R")
		t.A.NoError(err)

		_, res, err := rs.RunR(ctx, NewRunConfig(), "", "-e", "Sys.sleep(0.5); stop('failed')", "--slave", "--interactive")
		t.A.NoError(err)
		if err != nil {
			log.Error("goroutine 2 error:", err)
			log.Warn("cancelling ongoing work...")
			cancel()
		}
		fmt.Println("go routine 2 res: ", res)
	}()
	go func() {
		defer wg.Done()
		rs, err := NewRSettings("R")
		t.A.NoError(err)
		_, res, err := rs.RunR(ctx, NewRunConfig(), "", "-e", "Sys.sleep(1); 2+2", "--slave", "--interactive")
		t.A.NoError(err)

		if err != nil {
			log.Error("goroutine 3 error:", err)
		}
		fmt.Println("go routine 3 res: ", res)
	}()
	wg.Wait()
	fmt.Println("completed everything....")
}
