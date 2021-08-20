package main

import (
	"context"
	"fmt"
	"github.com/metrumresearchgroup/rcmd"
	"github.com/metrumresearchgroup/rcmd/writers"
	"os"
)

func main()  {
	rs, _ := rcmd.NewRSettings("R")
	rc :=rcmd.NewRunConfig()
	res, err  := rs.RunRWithOutput(context.Background(), rc, "", "--quiet", "-e", "2+2")
	if err != nil {
		panic(err)
	}
	fmt.Println("-------- simple addition -------")
	fmt.Println(string(res))

	res, err  = rs.RunRWithOutput(context.Background(), rc, "", "--quiet", "-e", "2+2")
	if err != nil {
		panic(err)
	}

	// can use an filter
	outputter := writers.NewFilter(os.Stdout, writers.LineNumberStripper, writers.InputFilter)
	fmt.Println("-------- simple addition with results cleaned -------")
	outputter.Write(res)

}
