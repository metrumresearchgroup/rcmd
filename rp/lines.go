package rp

import (
	"bufio"
	"bytes"

	"github.com/metrumresearchgroup/rcmd/writers"
)

// ScanLines scans lines from Rscript output and returns an array with
// the line numbers removed and whitespace trimmed
func ScanLines(b []byte) ([]string, error) {
	return ScanROutput(b, false)
}

// ScanROutput scans lines from RScript output and returns an array with
// the line numbers removed, whitespace trimmed, and (optionally)
// with all input-like lines (which start with ">") excluded.
func ScanROutput(b []byte, outputOnly bool) ([]string, error) {
	var fns []writers.FilterFunc
	if outputOnly {
		fns = append(fns, writers.InputFilter)
	}
	fns = append(fns, writers.LineNumberStripper)
	fns = append(fns, bytes.TrimSpace)

	output := &bytes.Buffer{}
	filter := writers.NewFilter(output, fns...)

	_, err := filter.Write(b)
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(output)
	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}
