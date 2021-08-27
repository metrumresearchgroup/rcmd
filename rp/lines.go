package rp

import (
	"bufio"
	"bytes"
	"io"
	"sync"

	"github.com/metrumresearchgroup/rcmd/v2/writers"
)

// ScanLines scans lines from Rscript output and returns an array with
// the line numbers removed and whitespace trimmed.
func ScanLines(b []byte) ([]string, error) {
	return ScanROutput(b, false)
}

// ScanROutput scans lines from RScript output and returns an array with
// the line numbers removed, whitespace trimmed, and (optionally)
// with all input-like lines (which start with ">") excluded.
func ScanROutput(b []byte, outputOnly bool) ([]string, error) {
	if !bytes.HasSuffix(b, []byte{'\n'}) {
		b = append(b, '\n')
	}
	var fns []writers.FilterFunc
	if outputOnly {
		fns = append(fns, writers.InputFilter)
	}
	fns = append(fns, writers.LineNumberStripper)
	fns = append(fns, bytes.TrimSpace)

	r, w := io.Pipe()

	var lines []string
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		scanner := bufio.NewScanner(r)
		for scanner.Scan() {
			lines = append(lines, scanner.Text())
		}
		wg.Done()
	}()

	filter := writers.NewFilter(w, fns...)

	_, err := filter.Write(b)
	if err != nil {
		return nil, err
	}

	err = filter.Close()
	if err != nil {
		return nil, err
	}

	wg.Wait()

	return lines, nil
}
