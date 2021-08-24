// +build !windows

package rcmd

import (
	"context"
	"testing"

	"github.com/metrumresearchgroup/wrapt"
)

func Test_example(tt *testing.T) {
	t := wrapt.WrapT(tt)

	cmd, err := New(context.Background(), "", "--quiet", "-e", "2+2")
	t.R.NoError(err)

	co, err := cmd.CombinedOutput()
	t.R.NoError(err)

	t.R.Equal("> 2+2\n[1] 4\n> \n> \n", string(co))
}

/*
func Test_RunR_example(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	p, n, err := rs.RunR(context.Background(), New(WithPrefix("foo")), "", "-e", "{2+2}")
	t.A.NoError(err)
	t.A.Equal(0, n)

	_, err = io.ReadAll(p.Stdout)
	t.A.NoError(err)
}

func Test_startR_example2(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	i, err := rs.StartR(context.Background(), New(), "", "-e", "2+2", "slave")
	t.A.NoError(err)
	defer fmt.Println(i.Stop())
}

func Test_runR_expression(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	ps, err := rs.RunRWithOutput(context.Background(), New(), "", "-e", "2+2", "--slave")
	t.A.NoError(err)

	lines, err := rp.ScanLines(ps)
	t.A.NoError(err)

	fmt.Println(strings.Join(lines, "\n"))
}

func Test_runRWithOutput_example(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	bs, err := rs.RunRWithOutput(context.Background(), New(), "", "-e", "2+2", "--slave", "--interactive")
	t.A.NoError(err)

	fmt.Println(string(bs))

	fmt.Println("----with prefix----")

	rs, err = NewRSettings("R")
	t.A.NoError(err)

	bs, err = rs.RunRWithOutput(context.Background(), New(), "", "-e", "2+2", "--slave", "--interactive")
	t.A.NoError(err)

	fmt.Println(string(bs))
}

func Test_runR_exampleTimeout(tt *testing.T) {
	t := wrapt.WrapT(tt)

	rs, err := NewRSettings("R")
	t.A.NoError(err)

	ctx, ccl := context.WithTimeout(context.Background(), 1*time.Second)
	defer ccl()

	_, res, err := rs.RunR(ctx, New(), "", "-e", "Sys.sleep(1.5); 2+2", "--slave", "--interactive")
	t.A.Error(err)

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

		i, err := rs.StartR(ctx, New(), "", "-e", "2+2", "--slave", "--interactive")
		t.A.NoError(err)
		defer fmt.Println(i.Stop())
		res, err := io.ReadAll(i.Pipes().Stdout)
		t.A.NoError(err)
		fmt.Println("go routine 1: ", res)
	}()
	go func() {
		defer wg.Done()

		rs, err := NewRSettings("R")
		t.A.NoError(err)

		i, err := rs.StartR(ctx, New(), "", "-e", "Sys.sleep(0.5); stop('failed')", "--slave", "--interactive")
		t.A.NoError(err)
		defer fmt.Println(i.Stop())
		if err != nil {
			log.Error("goroutine 2 error:", err)
			log.Warn("cancelling ongoing work...")
			cancel()
		}
		res, err := io.ReadAll(i.Pipes().Stdout)
		t.A.NoError(err)
		fmt.Println("go routine 2 res: ", res)
	}()
	go func() {
		defer wg.Done()
		rs, err := NewRSettings("R")
		t.A.NoError(err)
		i, err := rs.StartR(ctx, New(), "", "-e", "Sys.sleep(1); 2+2", "--slave", "--interactive")
		t.A.NoError(err)
		defer fmt.Println(i.Stop())
		if err != nil {
			log.Error("goroutine 3 error:", err)
		}
		res, err := io.ReadAll(i.Pipes().Stdout)
		t.A.NoError(err)
		fmt.Println("go routine 3 res: ", res)
	}()
	wg.Wait()
	fmt.Println("completed everything....")
}
*/
